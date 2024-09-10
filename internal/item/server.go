package item

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/goverland-labs/goverland-core-feed/protocol/feedpb"
)

const (
	defaultLimit  = 50
	defaultOffset = 0
)

var feedItemTypeMap = map[Type]feedpb.FeedInfo_Type{
	TypeDao:      feedpb.FeedInfo_DAO,
	TypeProposal: feedpb.FeedInfo_Proposal,
	TypeDelegate: feedpb.FeedInfo_Delegate,
}

var timelineActionsMap = map[TimelineAction]feedpb.FeedTimelineItem_TimelineAction{
	DaoCreated:                  feedpb.FeedTimelineItem_DaoCreated,
	DaoUpdated:                  feedpb.FeedTimelineItem_DaoUpdated,
	ProposalCreated:             feedpb.FeedTimelineItem_ProposalCreated,
	ProposalUpdated:             feedpb.FeedTimelineItem_ProposalUpdated,
	ProposalVotingStartsSoon:    feedpb.FeedTimelineItem_ProposalVotingStartsSoon,
	ProposalVotingEndsSoon:      feedpb.FeedTimelineItem_ProposalVotingEndsSoon,
	ProposalVotingStarted:       feedpb.FeedTimelineItem_ProposalVotingStarted,
	ProposalVotingQuorumReached: feedpb.FeedTimelineItem_ProposalVotingQuorumReached,
	ProposalVotingEnded:         feedpb.FeedTimelineItem_ProposalVotingEnded,
	DelegateCreateProposal:      feedpb.FeedTimelineItem_DelegateCreateProposal,
}

type Server struct {
	feedpb.UnimplementedFeedServer

	service *Service
}

func NewServer(sp *Service) *Server {
	return &Server{
		service: sp,
	}
}

func (s *Server) GetByFilter(_ context.Context, req *feedpb.FeedByFilterRequest) (*feedpb.FeedByFilterResponse, error) {
	limit, offset := defaultLimit, defaultOffset
	if req.GetLimit() > 0 {
		limit = int(req.GetLimit())
	}
	if req.GetOffset() > 0 {
		offset = int(req.GetOffset())
	}
	filters := []Filter{
		SkipSpammed{},
		SkipCanceled{},
		SkipDelegates{}, // we don't want to see delegates events in the dao feed
		PageFilter{Limit: limit, Offset: offset},
		SortedByCreated{
			Direction: DirectionDesc,
		},
	}

	// nolint:staticcheck // todo: deprecated. remove after updating core-api version in all related services
	if req.GetDaoId() != "" {
		filters = append(filters, DaoIDFilter{IDs: []string{req.GetDaoId()}})
	}
	if len(req.GetDaoIds()) > 0 {
		filters = append(filters, DaoIDFilter{IDs: req.GetDaoIds()})
	}

	if len(req.GetActions()) != 0 {
		filters = append(filters, ActionFilter{Actions: req.GetActions()})
	}

	if len(req.GetTypes()) != 0 {
		filters = append(filters, TypeFilter{Types: req.GetTypes()})
	}

	if req.IsActive != nil {
		filters = append(filters, ActiveFilter{IsActive: req.GetIsActive()})
	}

	list, err := s.service.GetByFilters(filters)
	if err != nil {
		log.Error().Err(err).Msg("get by filters")

		return nil, status.Error(codes.Internal, "internal error")
	}

	items := make([]*feedpb.FeedInfo, len(list.Items))
	for i, info := range list.Items {
		items[i] = convertFeedItemToAPI(&info)
	}

	return &feedpb.FeedByFilterResponse{
		Items:      items,
		TotalCount: uint64(list.TotalCount),
	}, nil
}

func convertFeedItemToAPI(item *FeedItem) *feedpb.FeedInfo {
	return &feedpb.FeedInfo{
		Id:           item.ID.String(),
		CreatedAt:    timestamppb.New(item.CreatedAt),
		UpdatedAt:    timestamppb.New(item.UpdatedAt),
		DaoId:        item.DaoID.String(),
		ProposalId:   item.ProposalID,
		DiscussionId: item.DiscussionID,
		Action:       string(convertActionToExternal(item.Action)),
		Type:         convertTypeToAPI(item.Type),
		Snapshot:     &anypb.Any{Value: item.Snapshot},
		Timeline:     convertTimelineToProto(item.Timeline),
	}
}

func convertTimelineToProto(timeline Timeline) []*feedpb.FeedTimelineItem {
	converted := make([]*feedpb.FeedTimelineItem, 0, len(timeline))
	for i := range timeline {
		converted = append(converted, &feedpb.FeedTimelineItem{
			CreatedAt: timestamppb.New(timeline[i].CreatedAt),
			Action:    convertTimelineActionToProto(timeline[i].Action),
		})
	}

	return converted
}

func convertTypeToAPI(t Type) feedpb.FeedInfo_Type {
	converted, exists := feedItemTypeMap[t]
	if !exists {
		log.Warn().Str("action", string(t)).Msg("unknown feed item type")

		return feedpb.FeedInfo_Unspecified
	}

	return converted
}

func convertTimelineActionToProto(action TimelineAction) feedpb.FeedTimelineItem_TimelineAction {
	converted, exists := timelineActionsMap[action]
	if !exists {
		log.Warn().Str("action", string(action)).Msg("unknown timeline action")

		return feedpb.FeedTimelineItem_Unspecified
	}

	return converted
}
