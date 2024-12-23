package item

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pevents "github.com/goverland-labs/goverland-platform-events/events/core"
	client "github.com/goverland-labs/goverland-platform-events/pkg/natsclient"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-core-feed/internal/config"
)

const (
	delegatesMaxPendingElements = 100
	delegatesRateLimit          = 500 * client.KiB
	delegatesExecutionTtl       = time.Minute
)

type DelegatesConsumer struct {
	conn      *nats.Conn
	service   *Service
	consumers []*client.Consumer[pevents.DelegatePayload]
}

func NewDelegatesConsumer(nc *nats.Conn, s *Service) (*DelegatesConsumer, error) {
	c := &DelegatesConsumer{
		conn:      nc,
		service:   s,
		consumers: make([]*client.Consumer[pevents.DelegatePayload], 0),
	}

	return c, nil
}

func (c *DelegatesConsumer) handler(action string) pevents.DelegatesHandler {
	return func(payload pevents.DelegatePayload) error {
		item, err := c.convertToFeedItem(convertToTimeLineAction(action), payload)
		if err != nil {
			return err
		}

		if err = c.service.HandleItem(context.TODO(), item, true); err != nil {
			log.Error().
				Str("dao_id", payload.DaoID.String()).
				Err(err).
				Msg("process delegate event")

			return err
		}

		log.Debug().Msgf("delegate was processed: %s", payload.Initiator)

		return nil
	}
}

func convertToTimeLineAction(action string) TimelineAction {
	switch action {
	case pevents.SubjectDelegateCreateProposal:
		return DelegateCreateProposal
	case pevents.SubjectDelegateVotingVoted:
		return DelegateVotingVoted
	case pevents.SubjectDelegateVotingSkipVote:
		return DelegateVotingSkipVote
	case pevents.SubjectDelegateCreated:
		return DelegateCreated
	default:
		return None
	}
}

func (c *DelegatesConsumer) convertToFeedItem(action TimelineAction, pl pevents.DelegatePayload) (*FeedItem, error) {
	b, err := json.Marshal(pl)
	if err != nil {
		return nil, fmt.Errorf("cant marshal payload: %w", err)
	}

	// we don't need timeline cause delegates feed its events
	return &FeedItem{
		DaoID:      pl.DaoID,
		ProposalID: pl.ProposalID,
		Type:       TypeDelegate,
		Action:     action,
		Snapshot:   b,
	}, nil
}

func (c *DelegatesConsumer) Start(ctx context.Context) error {
	group := config.GenerateGroupName("item_dao")

	opts := []client.ConsumerOpt{
		client.WithRateLimit(delegatesRateLimit),
		client.WithMaxAckPending(delegatesMaxPendingElements),
		client.WithAckWait(delegatesExecutionTtl),
	}

	for _, subj := range []string{
		pevents.SubjectDelegateCreateProposal,
		pevents.SubjectDelegateVotingVoted,
		pevents.SubjectDelegateVotingSkipVote,
		pevents.SubjectDelegateCreated,
	} {
		consumer, err := client.NewConsumer(ctx, c.conn, group, subj, c.handler(subj), opts...)
		if err != nil {
			return fmt.Errorf("consume for %s/%s: %w", group, subj, err)
		}

		c.consumers = append(c.consumers, consumer)
	}

	log.Info().Msg("feed item Delegates consumers is started")

	// todo: handle correct stopping the consumer by context
	<-ctx.Done()
	return c.stop()
}

func (c *DelegatesConsumer) stop() error {
	for _, cs := range c.consumers {
		if err := cs.Close(); err != nil {
			log.Error().Err(err).Msg("unable to close delegates consumer")
		}
	}

	return nil
}
