package item

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/goverland-labs/platform-events/events/core"
	"github.com/goverland-labs/platform-events/events/inbox"
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

	data, err := json.Marshal(convertToExternalFeed(item))
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

		err = s.events.PublishJSON(ctx, core.SubjectCallback, core.CallbackPayload{
			WebhookURL: info.WebhookURL,
			Body:       data,
		})
		if err != nil {
			log.Error().Str("subscriber", sub.String()).Str("webhook_url", info.WebhookURL).Err(err).Msgf("publish callback")
		}
	}

	return nil
}

func convertFeedType(ftype Type) inbox.Type {
	switch ftype {
	case TypeDao:
		return inbox.TypeDao
	case TypeProposal:
		return inbox.TypeProposal
	default:
		return inbox.TypeDao
	}
}

func convertToExternalFeed(item *FeedItem) inbox.FeedPayload {
	// TODO: TBD: Might be we should export feed.id?

	return inbox.FeedPayload{
		DaoID:        item.DaoID,
		ProposalID:   item.ProposalID,
		DiscussionID: item.DiscussionID,
		Type:         convertFeedType(item.Type),
		Action:       convertActionToExternal(item.Action),
		Snapshot:     item.Snapshot,
		Timeline:     convertToExternalTimeline(item.Timeline),
	}
}

func convertToExternalTimeline(timeline Timeline) []inbox.TimelineItem {
	converted := make([]inbox.TimelineItem, 0, len(timeline))

	for _, t := range timeline {
		action := convertActionToExternal(t.Action)
		if action == "" {
			// TODO: log warning
			continue
		}

		converted = append(converted, inbox.TimelineItem{
			CreatedAt: t.CreatedAt,
			Action:    action,
		})
	}

	return converted
}

var inboxTimelineActionMap = map[TimelineAction]inbox.TimelineAction{
	DaoCreated:                  inbox.DaoCreated,
	DaoUpdated:                  inbox.DaoUpdated,
	ProposalCreated:             inbox.ProposalCreated,
	ProposalUpdated:             inbox.ProposalUpdated,
	ProposalVotingStartsSoon:    inbox.ProposalVotingStartsSoon,
	ProposalVotingStarted:       inbox.ProposalVotingStarted,
	ProposalVotingQuorumReached: inbox.ProposalVotingQuorumReached,
	ProposalVotingEnded:         inbox.ProposalVotingEnded,
}

func convertActionToExternal(action TimelineAction) inbox.TimelineAction {
	return inboxTimelineActionMap[action]
}

func (s *Service) GetByFilters(filters []Filter) (FeedList, error) {
	list, err := s.repo.GetByFilters(filters)
	if err != nil {
		return FeedList{}, fmt.Errorf("get by filters: %w", err)
	}

	return list, nil
}
