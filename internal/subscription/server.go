package subscription

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	proto "github.com/goverland-labs/core-api/protobuf/internalapi"

	"github.com/goverland-labs/core-feed/internal/subscriber"
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
	subID := subscriber.GetSubscriberID(ctx)

	if req.GetDaoId() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid dao id")
	}

	daoId, err := uuid.Parse(req.GetDaoId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid dao id: %s", err))
	}

	_, err = s.sp.Subscribe(ctx, Subscription{
		SubscriberID: subID,
		DaoID:        daoId,
	})
	if err != nil {
		log.Error().Err(err).Msgf("subscribe: %s - %s", subID, req.GetDaoId())
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.Debug().Msgf("subscribe: %s - %s", subID, req.GetDaoId())

	return &emptypb.Empty{}, nil
}

func (s *Server) Unsubscribe(ctx context.Context, req *proto.UnsubscribeRequest) (*emptypb.Empty, error) {
	subID := subscriber.GetSubscriberID(ctx)

	if req.GetDaoId() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid dao id")
	}

	daoId, err := uuid.Parse(req.GetDaoId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid dao id: %s", err))
	}

	err = s.sp.Unsubscribe(ctx, Subscription{
		SubscriberID: subID,
		DaoID:        daoId,
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.InvalidArgument, "invalid subscription")
	}

	if err != nil {
		log.Error().Err(err).Msgf("unsubscribe: %s - %s", subID, req.GetDaoId())
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.Debug().Msgf("unsubscribe: %s - %s", subID, req.GetDaoId())

	return &emptypb.Empty{}, nil
}
