package subscriber

import (
	"context"
	"errors"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	proto "github.com/goverland-labs/core-api/protobuf/internalapi"
)

type SubscriberProvider interface {
	GetByID(_ context.Context, id uuid.UUID) (*Subscriber, error)
	Create(_ context.Context, webhookURL string) (*Subscriber, error)
	Update(_ context.Context, item Subscriber) error
}

type Server struct {
	proto.UnimplementedSubscriberServer

	sp SubscriberProvider
}

func NewServer(sp SubscriberProvider) *Server {
	return &Server{
		sp: sp,
	}
}

func (s *Server) Create(ctx context.Context, req *proto.CreateSubscriberRequest) (*proto.CreateSubscriberResponse, error) {
	if req.GetWebhookUrl() != "" {
		if _, err := url.ParseRequestURI(req.GetWebhookUrl()); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid webhook url")
		}
	}

	sub, err := s.sp.Create(ctx, req.GetWebhookUrl())
	if err != nil {
		log.Error().Err(err).Msg("create subscriber")

		return nil, status.Error(codes.Internal, "internal error")
	}

	log.Debug().Msgf("create subscriber: %s", sub.ID)

	return &proto.CreateSubscriberResponse{SubscriberId: sub.ID.String()}, nil
}

func (s *Server) Update(ctx context.Context, req *proto.UpdateSubscriberRequest) (*emptypb.Empty, error) {
	subID := GetSubscriberID(ctx)
	if req.GetWebhookUrl() != "" {
		if _, err := url.ParseRequestURI(req.GetWebhookUrl()); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid webhook url")
		}
	}

	err := s.sp.Update(ctx, Subscriber{
		ID:         subID,
		WebhookURL: req.GetWebhookUrl(),
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.InvalidArgument, "invalid subscriber ID")
	}

	if err != nil {
		log.Error().Err(err).Msgf("update subscriber: %s", subID)
		return nil, status.Error(codes.Internal, "internal error")
	}

	log.Debug().Msgf("update subscriber: %s", subID)

	return &emptypb.Empty{}, nil
}
