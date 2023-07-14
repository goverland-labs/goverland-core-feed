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

type DaoConsumer struct {
	conn      *nats.Conn
	service   *Service
	consumers []*client.Consumer
}

func NewDaoConsumer(nc *nats.Conn, s *Service) (*DaoConsumer, error) {
	c := &DaoConsumer{
		conn:      nc,
		service:   s,
		consumers: make([]*client.Consumer, 0),
	}

	return c, nil
}

func (c *DaoConsumer) handler(action string) pevents.DaoHandler {
	return func(payload pevents.DaoPayload) error {
		var err error
		defer func(start time.Time) {
			metricHandleHistogram.
				WithLabelValues("dao", metrics.ErrLabelValue(err)).
				Observe(time.Since(start).Seconds())
		}(time.Now())

		item, err := c.service.GetDaoItem(context.TODO(), payload.ID)
		if err != nil {
			return err
		}

		var timeline Timeline
		if item != nil {
			timeline = item.Timeline
		}

		var sendUpdates = true
		now := time.Now().UTC()
		switch action {
		case pevents.SubjectDaoCreated:
			sendUpdates = timeline.AddUniqueAction(now, DaoCreated)
		case pevents.SubjectDaoUpdated:
			timeline.AddNonUniqueAction(now, DaoUpdated)
		}

		timeline = c.prefillTimelineInNeeded(payload, timeline)

		if item == nil {
			item, err = c.convertToFeedItem(payload, timeline)
			if err != nil {
				return err
			}
		} else {
			item.Timeline = timeline
		}

		err = c.service.HandleItem(context.TODO(), item, sendUpdates)
		if err != nil {
			log.Error().Str("dao_id", payload.ID.String()).Err(err).Msg("process dao")
			return err
		}

		log.Debug().Msgf("dao was processed: %s", payload.ID)

		return nil
	}
}

func (c *DaoConsumer) convertToFeedItem(pl pevents.DaoPayload, timeline Timeline) (*FeedItem, error) {
	b, err := json.Marshal(pl)
	if err != nil {
		return nil, fmt.Errorf("cant marshal payload: %w", err)
	}

	return &FeedItem{
		DaoID:    pl.ID,
		Type:     TypeDao,
		Action:   timeline.LastAction(),
		Snapshot: b,
		Timeline: timeline,
	}, nil
}

func (c *DaoConsumer) Start(ctx context.Context) error {
	group := config.GenerateGroupName("item_dao")

	for _, subj := range []string{pevents.SubjectDaoCreated, pevents.SubjectDaoUpdated} {
		consumer, err := client.NewConsumer(ctx, c.conn, group, subj, c.handler(subj))
		if err != nil {
			return fmt.Errorf("consume for %s/%s: %w", group, subj, err)
		}

		c.consumers = append(c.consumers, consumer)
	}

	log.Info().Msg("feed item DAO consumers is started")

	// todo: handle correct stopping the consumer by context
	<-ctx.Done()
	return c.stop()
}

func (c *DaoConsumer) stop() error {
	for _, cs := range c.consumers {
		if err := cs.Close(); err != nil {
			log.Error().Err(err).Msg("cant close feed item consumer")
		}
	}

	return nil
}

func (c *DaoConsumer) prefillTimelineInNeeded(payload pevents.DaoPayload, timeline Timeline) Timeline {
	prepend := make([]TimelineItem, 0, 3)

	if !timeline.ContainsAction(DaoCreated) {
		prepend = append(prepend, TimelineItem{
			CreatedAt: time.Now(),
			Action:    DaoCreated,
		})
	}

	timeline = append(prepend, timeline...)
	timeline.Sort()

	return timeline
}
