package item

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pevents "github.com/goverland-labs/platform-events/events/core"
	client "github.com/goverland-labs/platform-events/pkg/natsclient"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/core-feed/internal/config"
	"github.com/goverland-labs/core-feed/internal/metrics"
)

type ProposalConsumer struct {
	conn      *nats.Conn
	service   *Service
	consumers []*client.Consumer
}

func NewProposalConsumer(nc *nats.Conn, s *Service) (*ProposalConsumer, error) {
	c := &ProposalConsumer{
		conn:      nc,
		service:   s,
		consumers: make([]*client.Consumer, 0),
	}

	return c, nil
}

func (c *ProposalConsumer) handler(action string) pevents.ProposalHandler {
	return func(payload pevents.ProposalPayload) error {
		var err error
		defer func(start time.Time) {
			metricHandleHistogram.
				WithLabelValues("dao", metrics.ErrLabelValue(err)).
				Observe(time.Since(start).Seconds())
		}(time.Now())

		item, err := c.convertToFeedItem(payload, action)
		if err != nil {
			log.Error().Err(err).Msg("converting feed item")
			return err
		}

		err = c.service.HandleItem(context.TODO(), item)
		if err != nil {
			log.Error().Err(err).Msg("process dao")
			return err
		}

		log.Debug().Msgf("dao was processed: %s", payload.ID)

		return nil
	}
}

func (c *ProposalConsumer) convertToFeedItem(pl pevents.ProposalPayload, action string) (FeedItem, error) {
	b, err := json.Marshal(pl)
	if err != nil {
		return FeedItem{}, fmt.Errorf("cant marshal payload: %w", err)
	}

	return FeedItem{
		DaoID:      pl.DaoID,
		ProposalID: pl.ID,
		Type:       TypeProposal,
		Action:     action,
		Snapshot:   b,
	}, nil
}

func (c *ProposalConsumer) Start(ctx context.Context) error {
	group := config.GenerateGroupName("item_proposal")

	for _, event := range []string{
		pevents.SubjectProposalCreated,
		pevents.SubjectProposalUpdated,
		pevents.SubjectProposalVotingStartsSoon,
		pevents.SubjectProposalVotingStarted,
		pevents.SubjectProposalVotingEnded,
		pevents.SubjectProposalVotingQuorumReached,
		pevents.SubjectProposalUpdatedState,
	} {
		cc, err := client.NewConsumer(ctx, c.conn, group, event, c.handler(event))
		if err != nil {
			return fmt.Errorf("consume for %s/%s: %w", group, event, err)
		}
		c.consumers = append(c.consumers, cc)
	}

	log.Info().Msg("feed item proposal consumers is started")

	// todo: handle correct stopping the consumer by context
	<-ctx.Done()
	return c.stop()
}

func (c *ProposalConsumer) stop() error {
	for _, cs := range c.consumers {
		if err := cs.Close(); err != nil {
			log.Error().Err(err).Msg("cant close feed item proposal consumer")
		}
	}

	return nil
}
