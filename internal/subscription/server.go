package subscription

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	proto "github.com/goverland-labs/feed/protobuf/internalapi"
)

type SubscriptionProvider interface {
	Subscribe(_ context.Context, item Subscription) (*Subscription, error)
	Unsubscribe(_ context.Context, item Subscription) error
}

type Server struct {
	proto.UnimplementedSubscriptionServer

	sp SubscriptionProvider
}

func NewServer(sp SubscriptionProvider) *Server {
	return &Server{
		sp: sp,
	}
}

func (s *Server) Subscribe(ctx context.Context, req *proto.SubscribeRequest) (*emptypb.Empty, error) {
	if req.GetSubscriberId() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid subscriber ID")
	}

	if req.GetDaoId() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID")
	}

	_, err := s.sp.Subscribe(ctx, Subscription{
		SubscriberID: req.GetSubscriberId(),
		DaoID:        req.GetDaoId(),
	})
	if err != nil {
		log.Error().Err(err).Msgf("subscribe: %s - %s", req.GetSubscriberId(), req.GetDaoId())
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.Debug().Msgf("subscribe: %s - %s", req.GetSubscriberId(), req.GetDaoId())

	return &emptypb.Empty{}, nil
}

func (s *Server) Unsubscribe(ctx context.Context, req *proto.UnsubscribeRequest) (*emptypb.Empty, error) {
	if req.GetSubscriberId() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid subscriber ID")
	}

	if req.GetDaoId() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID")
	}

	err := s.sp.Unsubscribe(ctx, Subscription{
		SubscriberID: req.GetSubscriberId(),
		DaoID:        req.GetDaoId(),
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.InvalidArgument, "invalid subscription")
	}

	if err != nil {
		log.Error().Err(err).Msgf("unsubscribe: %s - %s", req.GetSubscriberId(), req.GetDaoId())
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.Debug().Msgf("unsubscribe: %s - %s", req.GetSubscriberId(), req.GetDaoId())

	return &emptypb.Empty{}, nil
}
