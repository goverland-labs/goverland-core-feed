package feedevent

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-core-feed/internal/item"
	"github.com/goverland-labs/goverland-core-feed/internal/pubsub"
)

const (
	forcedFetchTime = 1 * time.Minute

	feedItemsLimit = 1000
)

type FeedItemsProvider interface {
	GetLastItems(subscriberID string, lastUpdatedAt time.Time, limit int) ([]item.FeedItem, error)
}

type Service struct {
	notifier *pubsub.PubSub[string]

	feedItemsProvider FeedItemsProvider
}

func NewService(notifier *pubsub.PubSub[string], feedItemsProvider FeedItemsProvider) *Service {
	return &Service{notifier: notifier, feedItemsProvider: feedItemsProvider}
}

func (s *Service) Watch(
	ctx context.Context,
	subscriberID uuid.UUID,
	lastUpdatedAt time.Time,
	handler func(entity item.FeedItem) error,
) error {
	notificationsCh := s.notifier.Subscribe()
	defer func() {
		s.notifier.Unsubscribe(notificationsCh)
	}()

	for {
		feedItems, err := s.feedItemsProvider.GetLastItems(subscriberID.String(), lastUpdatedAt, feedItemsLimit)
		if err != nil {
			return fmt.Errorf("fail to fetch last feed items: %v", err)
		}

		log.Info().
			Int("count", len(feedItems)).
			Str("subscriber", subscriberID.String()).
			Msg("fetched feed items")

		for _, feedItem := range feedItems {
			err := handler(feedItem)
			if err != nil {
				return fmt.Errorf("fail to handle feed item in subscription: %v", err)
			}

			lastUpdatedAt = feedItem.UpdatedAt
		}

		log.Info().
			Str("value", lastUpdatedAt.String()).
			Msg("change lastUpdatedAt")

		if len(feedItems) == feedItemsLimit {
			continue
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
