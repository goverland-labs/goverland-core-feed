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

	"github.com/goverland-labs/feed/internal/metrics"
)

const (
	groupName = "item"
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

func (c *DaoConsumer) convertToFeedItem(pl pevents.DaoPayload, action string) (FeedItem, error) {
	b, err := json.Marshal(pl)
	if err != nil {
		return FeedItem{}, fmt.Errorf("cant marshal payload: %w", err)
	}

	return FeedItem{
		DaoID:    pl.ID,
		Type:     TypeDao,
		Action:   action,
		Snapshot: b,
	}, nil
}

func (c *DaoConsumer) Start(ctx context.Context) error {
	cc, err := client.NewConsumer(ctx, c.conn, groupName, pevents.SubjectDaoCreated, c.handler(pevents.SubjectDaoCreated))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", groupName, pevents.SubjectDaoCreated, err)
	}
	cu, err := client.NewConsumer(ctx, c.conn, groupName, pevents.SubjectDaoUpdated, c.handler(pevents.SubjectDaoUpdated))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", groupName, pevents.SubjectDaoUpdated, err)
	}

	c.consumers = append(c.consumers, cc, cu)

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
