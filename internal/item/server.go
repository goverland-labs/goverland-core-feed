package item

import (
	"context"

	proto "github.com/goverland-labs/core-api/protobuf/internalapi"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultLimit  = 50
	defaultOffset = 0
)

var feedItemTypeMap = map[Type]proto.FeedInfo_Type{
	TypeDao:      proto.FeedInfo_DAO,
	TypeProposal: proto.FeedInfo_Proposal,
}

var timelineActionsMap = map[TimelineAction]proto.FeedTimelineItem_TimelineAction{
	DaoCreated:                  proto.FeedTimelineItem_DaoCreated,
	DaoUpdated:                  proto.FeedTimelineItem_DaoUpdated,
	ProposalCreated:             proto.FeedTimelineItem_ProposalCreated,
	ProposalUpdated:             proto.FeedTimelineItem_ProposalCreated,
	ProposalVotingStartsSoon:    proto.FeedTimelineItem_ProposalCreated,
	ProposalVotingStarted:       proto.FeedTimelineItem_ProposalCreated,
	ProposalVotingQuorumReached: proto.FeedTimelineItem_ProposalCreated,
	ProposalVotingEnded:         proto.FeedTimelineItem_ProposalCreated,
}

type Server struct {
	proto.UnimplementedFeedServer

	service *Service
}

func NewServer(sp *Service) *Server {
	return &Server{
		service: sp,
	}
}

func (s *Server) GetByFilter(_ context.Context, req *proto.FeedByFilterRequest) (*proto.FeedByFilterResponse, error) {
	limit, offset := defaultLimit, defaultOffset
	if req.GetLimit() > 0 {
		limit = int(req.GetLimit())
	}
	if req.GetOffset() > 0 {
		offset = int(req.GetOffset())
	}
	filters := []Filter{
		PageFilter{Limit: limit, Offset: offset},
		OrderByCreatedFilter{},
	}

	if req.GetDaoId() != "" {
		filters = append(filters, DaoIDFilter{ID: req.GetDaoId()})
	}

	if len(req.GetActions()) != 0 {
		filters = append(filters, ActionFilter{Actions: req.GetActions()})
	}

	if len(req.GetActions()) != 0 {
		filters = append(filters, ActionFilter{Actions: req.GetActions()})
	}

	if len(req.GetTypes()) != 0 {
		filters = append(filters, TypeFilter{Types: req.GetTypes()})
	}

	list, err := s.service.GetByFilters(filters)
	if err != nil {
		log.Error().Err(err).Msg("get by filters")

		return nil, status.Error(codes.Internal, "internal error")
	}

	items := make([]*proto.FeedInfo, len(list.Items))
	for i, info := range list.Items {
		items[i] = convertFeedItemToAPI(&info)
	}

	return &proto.FeedByFilterResponse{
		Items:      items,
		TotalCount: uint64(list.TotalCount),
	}, nil
}

func convertFeedItemToAPI(item *FeedItem) *proto.FeedInfo {
	return &proto.FeedInfo{
		Id:           item.ID.String(),
		CreatedAt:    timestamppb.New(item.CreatedAt),
		UpdatedAt:    timestamppb.New(item.UpdatedAt),
		DaoId:        item.DaoID,
		ProposalId:   item.ProposalID,
		DiscussionId: item.DiscussionID,
		Action:       string(convertActionToExternal(item.Action)),
		Type:         convertTypeToAPI(item.Type),
		Snapshot:     &anypb.Any{Value: item.Snapshot},
		Timeline:     convertTimelineToProto(item.Timeline),
	}
}

func convertTimelineToProto(timeline Timeline) []*proto.FeedTimelineItem {
	converted := make([]*proto.FeedTimelineItem, 0, len(timeline))
	for i := range timeline {
		converted = append(converted, &proto.FeedTimelineItem{
			CreatedAt: timestamppb.New(timeline[i].CreatedAt),
			Action:    convertTimelineActionToProto(timeline[i].Action),
		})
	}

	return converted
}

func convertTypeToAPI(t Type) proto.FeedInfo_Type {
	converted, exists := feedItemTypeMap[t]
	if !exists {
		log.Warn().Str("action", string(t)).Msg("unknown feed item type")

		return proto.FeedInfo_Unspecified
	}

	return converted
}

func convertTimelineActionToProto(action TimelineAction) proto.FeedTimelineItem_TimelineAction {
	converted, exists := timelineActionsMap[action]
	if !exists {
		log.Warn().Str("action", string(action)).Msg("unknown timeline action")

		return proto.FeedTimelineItem_Unspecified
	}

	return converted
}
