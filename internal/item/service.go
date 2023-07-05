package item

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goverland-labs/platform-events/events/core"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/core-feed/internal/subscriber"
)

//go:generate mockgen -destination=mocks_test.go -package=item . DataProvider,Publisher

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

type DataProvider interface {
	Create(item FeedItem) error
	GetByFilters(filters []Filter) (FeedList, error)
}

type SubscriberProvider interface {
	GetByID(_ context.Context, id string) (*subscriber.Subscriber, error)
}

type SubscriptionProvider interface {
	GetSubscribers(_ context.Context, daoID string) ([]string, error)
}

type Service struct {
	repo          DataProvider
	events        Publisher
	subscribers   SubscriberProvider
	subscriptions SubscriptionProvider
}

func NewService(r DataProvider, p Publisher, sub SubscriberProvider, sp SubscriptionProvider) (*Service, error) {
	return &Service{
		repo:          r,
		events:        p,
		subscribers:   sub,
		subscriptions: sp,
	}, nil
}

func (s *Service) HandleItem(ctx context.Context, item FeedItem) error {
	err := s.repo.Create(item)
	if err != nil {
		return fmt.Errorf("can't create item: %w", err)
	}

	// todo: refactor and move to separated logic
	go func() {
		subs, err := s.subscriptions.GetSubscribers(ctx, item.DaoID)
		if err != nil {
			log.Error().Err(err).Msg("get subscribers")
			return
		}

		feed := convertToExternalFeed(item)
		data, err := json.Marshal(feed)
		if err != nil {
			log.Error().Err(err).Msgf("marshal feed: %d", item.ID)
			return
		}

		for _, sub := range subs {
			info, err := s.subscribers.GetByID(ctx, sub)
			if err != nil {
				log.Error().Err(err).Msgf("get details for: %s", sub)
				continue
			}

			payload := core.CallbackPayload{
				WebhookURL: info.WebhookURL,
				Body:       data,
			}

			err = s.events.PublishJSON(ctx, core.SubjectCallback, payload)
			if err != nil {
				log.Error().Err(err).Msgf("publish callback for: %s", info.WebhookURL)
			}
		}
	}()

	return nil
}

func convertFeedType(ftype Type) core.Type {
	switch ftype {
	case TypeDao:
		return core.TypeDao
	case TypeProposal:
		return core.TypeProposal
	default:
		return core.TypeDao
	}
}

func convertToExternalFeed(item FeedItem) core.FeedItem {
	return core.FeedItem{
		DaoID:        item.DaoID,
		ProposalID:   item.ProposalID,
		DiscussionID: item.DiscussionID,
		Type:         convertFeedType(item.Type),
		Action:       core.ConvertActionToExternal(item.Action),
		Snapshot:     item.Snapshot,
	}
}

func (s *Service) GetByFilters(filters []Filter) (FeedList, error) {
	list, err := s.repo.GetByFilters(filters)
	if err != nil {
		return FeedList{}, fmt.Errorf("get by filters: %w", err)
	}

	return list, nil
}
