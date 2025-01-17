package feedevent

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-core-feed/internal/item"
)

const (
	forcedFetchTime = 5 * time.Minute
)

type Service struct {
	name     string
	notifier *PubSub
}

func NewService(name string, notifier *PubSub) *Service {
	return &Service{name: name, notifier: notifier}
}

func (s *Service) Watch(
	ctx context.Context,
	name string,
	lastUpdatedAt time.Time,
	handler func(entity item.FeedItem) error, // TODO generic?
) error {
	notificationsCh := s.notifier.Subscribe()
	defer func() {
		s.notifier.Unsubscribe(notificationsCh)
	}()

	for {
		entities := []item.FeedItem{} // TODO fetch from db or read from notificationsCh, discuss ??
		//if err != nil {
		//	return fmt.Errorf("fail to fetch entites for subscription from database: %v", err)
		//}

		for _, entity := range entities {
			err := handler(entity)
			if err != nil {
				return fmt.Errorf("fail to handle entity in subscription: %v", err)
			}

			// TODO: update lastUpdatedAt
		}

		select {
		case <-ctx.Done():
			log.Info().Msg("ctx is done, finished subscription")
			return nil

		case <-notificationsCh:
		case <-time.After(forcedFetchTime):
		}
	}
}
