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
	}
}

func convertTypeToAPI(t Type) proto.FeedInfo_Type {
	switch t {
	case TypeDao:
		return proto.FeedInfo_TYPE_DAO
	case TypeProposal:
		return proto.FeedInfo_TYPE_PROPOSAL
	default:
		return proto.FeedInfo_TYPE_UNSPECIFIED
	}
}
