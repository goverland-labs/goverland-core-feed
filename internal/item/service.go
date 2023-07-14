package item

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/goverland-labs/platform-events/events/core"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/goverland-labs/core-feed/internal/subscriber"
)

//go:generate mockgen -destination=mocks_test.go -package=item . DataProvider,Publisher

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

type DataProvider interface {
	Save(item *FeedItem) error
	GetDaoItem(id uuid.UUID) (*FeedItem, error)
	GetProposalItem(id string) (*FeedItem, error)
	GetByFilters(filters []Filter) (FeedList, error)
}

type SubscriberProvider interface {
	GetByID(_ context.Context, id uuid.UUID) (*subscriber.Subscriber, error)
}

type SubscriptionProvider interface {
	GetSubscribers(_ context.Context, daoID uuid.UUID) ([]uuid.UUID, error)
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

func (s *Service) GetDaoItem(_ context.Context, id uuid.UUID) (*FeedItem, error) {
	item, err := s.repo.GetDaoItem(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) GetProposalItem(_ context.Context, id string) (*FeedItem, error) {
	item, err := s.repo.GetProposalItem(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) HandleItem(ctx context.Context, item *FeedItem, sendUpdates bool) error {
	if err := s.repo.Save(item); err != nil {
		return fmt.Errorf("can't save feed item: %w", err)
	}

	if !sendUpdates {
		return nil
	}

	subs, err := s.subscriptions.GetSubscribers(ctx, item.DaoID)
	if err != nil {
		log.Error().Err(err).Msg("get subscribers")
		return nil
	}

	feed := convertToExternalFeed(item)
	data, err := json.Marshal(feed)
	if err != nil {
		log.Error().Err(err).Msgf("marshal feed: %d", item.ID)

		// Suppress an error for the consumer for avoiding duplicated events
		return nil
	}

	for _, sub := range subs {
		info, err := s.subscribers.GetByID(ctx, sub)
		if err != nil {
			log.Error().Str("subscriber", sub.String()).Err(err).Msgf("get details for subscriber")
			continue
		}

		payload := core.CallbackPayload{
			WebhookURL: info.WebhookURL,
			Body:       data,
		}

		err = s.events.PublishJSON(ctx, core.SubjectCallback, payload)
		if err != nil {
			log.Error().Str("subscriber", sub.String()).Str("webhook_url", info.WebhookURL).Err(err).Msgf("publish callback")
		}
	}

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

func convertToExternalFeed(item *FeedItem) core.FeedItem {
	// TODO: TBD: Might be we should export feed.id?

	return core.FeedItem{
		DaoID:        item.DaoID,
		ProposalID:   item.ProposalID,
		DiscussionID: item.DiscussionID,
		Type:         convertFeedType(item.Type),
		Action:       convertActionToExternal(item.Action),
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

func convertActionToExternal(action TimelineAction) core.Action {
	switch action {
	case DaoCreated, ProposalCreated:
		return core.ActionCreated
	case DaoUpdated, ProposalUpdated:
		return core.ActionUpdated
	case ProposalVotingStarted:
		return core.ActionVotingStarted
	case ProposalVotingStartsSoon:
		return core.ActionVotingStartsSoon
	case ProposalVotingQuorumReached:
		return core.ActionVotingQuorumReached
	case ProposalVotingEnded:
		return core.ActionVotingEnded
	default:
		return ""
	}
}
