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
	"github.com/goverland-labs/goverland-core-feed/internal/metrics"
)

const (
	proposalMaxPendingElements = 100
	proposalRateLimit          = 500 * client.KiB
	proposalExecutionTtl       = time.Minute
)

// proposalEvents is map event:string => is_unique_event:bool
type eventConfig struct {
	isUnique  bool
	action    TimelineAction
	extractor func(payload pevents.ProposalPayload) time.Time
}

var defaultConfig = eventConfig{isUnique: false, action: ProposalUpdated}

var proposalEvents = map[string]eventConfig{
	pevents.SubjectProposalCreated: {isUnique: true, action: ProposalCreated, extractor: func(payload pevents.ProposalPayload) time.Time {
		return time.Unix(int64(payload.Created), 0).UTC()
	}},
	pevents.SubjectProposalVotingStartsSoon: {isUnique: true, action: ProposalVotingStartsSoon},
	pevents.SubjectProposalVotingEndsSoon:   {isUnique: true, action: ProposalVotingEndsSoon},
	pevents.SubjectProposalVotingStarted: {isUnique: true, action: ProposalVotingStarted, extractor: func(payload pevents.ProposalPayload) time.Time {
		return time.Unix(int64(payload.Start), 0).UTC()
	}},
	pevents.SubjectProposalVotingQuorumReached: {isUnique: true, action: ProposalVotingQuorumReached},
	pevents.SubjectProposalVotingEnded: {isUnique: true, action: ProposalVotingEnded, extractor: func(payload pevents.ProposalPayload) time.Time {
		return time.Unix(int64(payload.End), 0).UTC()
	}},
	pevents.SubjectProposalUpdated:      defaultConfig,
	pevents.SubjectProposalUpdatedState: defaultConfig,
}

type ProposalConsumer struct {
	conn      *nats.Conn
	service   *Service
	consumers []*client.Consumer[pevents.ProposalPayload]
}

func NewProposalConsumer(nc *nats.Conn, s *Service) (*ProposalConsumer, error) {
	c := &ProposalConsumer{
		conn:      nc,
		service:   s,
		consumers: make([]*client.Consumer[pevents.ProposalPayload], 0),
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

		item, err := c.service.GetProposalItem(context.TODO(), payload.ID)
		if err != nil {
			return err
		}

		var timeline Timeline
		if item != nil {
			timeline = item.Timeline
		}

		cfg, exist := proposalEvents[action]
		if !exist {
			cfg = defaultConfig
		}

		eventTime := time.Now().UTC()
		if cfg.extractor != nil {
			eventTime = cfg.extractor(payload)
		}

		var sendUpdates = true
		if cfg.isUnique {
			sendUpdates = timeline.AddUniqueAction(eventTime, cfg.action)
		} else {
			timeline.AddNonUniqueAction(eventTime, cfg.action)
		}

		timeline = c.prefillTimelineInNeeded(payload, timeline)

		if item == nil {
			item, err = c.convertToFeedItem(payload, timeline)
			if err != nil {
				return err
			}
		} else {
			sn, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("cant marshal payload: %w", err)
			}

			item.Snapshot = sn
			item.Timeline = timeline
		}

		err = c.service.HandleItem(context.TODO(), item, sendUpdates)
		if err != nil {
			log.Error().Err(err).Msg("process proposal")
			return err
		}

		log.Debug().Msgf("proposal was processed: %s", payload.ID)

		return nil
	}
}

func (c *ProposalConsumer) convertToFeedItem(pl pevents.ProposalPayload, timeline Timeline) (*FeedItem, error) {
	b, err := json.Marshal(pl)
	if err != nil {
		return nil, fmt.Errorf("cant marshal payload: %w", err)
	}

	return &FeedItem{
		DaoID:      pl.DaoID,
		ProposalID: pl.ID,
		Type:       TypeProposal,
		Action:     timeline.LastAction(),
		Snapshot:   b,
		Timeline:   timeline,
	}, nil
}

func (c *ProposalConsumer) Start(ctx context.Context) error {
	group := config.GenerateGroupName("item_proposal")

	opts := []client.ConsumerOpt{
		client.WithRateLimit(proposalRateLimit),
		client.WithMaxAckPending(proposalMaxPendingElements),
		client.WithAckWait(proposalExecutionTtl),
	}

	for event := range proposalEvents {
		cc, err := client.NewConsumer(ctx, c.conn, group, event, c.handler(event), opts...)
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

func (c *ProposalConsumer) prefillTimelineInNeeded(payload pevents.ProposalPayload, timeline Timeline) Timeline {
	prepend := make([]TimelineItem, 0, 3)

	if !timeline.ContainsAction(ProposalCreated) {
		prepend = append(prepend, TimelineItem{
			CreatedAt: time.Unix(int64(payload.Created), 0).UTC(),
			Action:    ProposalCreated,
		})
	}

	votingStartsAt := time.Unix(int64(payload.Start), 0).UTC()
	if votingStartsAt.Before(time.Now()) && !timeline.ContainsAction(ProposalVotingStarted) {
		prepend = append(prepend, TimelineItem{
			CreatedAt: votingStartsAt,
			Action:    ProposalVotingStarted,
		})
	}

	votingEndsAt := time.Unix(int64(payload.End), 0).UTC()
	if votingEndsAt.Before(time.Now()) && !timeline.ContainsAction(ProposalVotingEnded) {
		prepend = append(prepend, TimelineItem{
			CreatedAt: votingEndsAt,
			Action:    ProposalVotingEnded,
		})
	}

	timeline = append(prepend, timeline...)

	return timeline
}
