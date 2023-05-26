package subscriber

import (
	"context"
	"errors"
	"net/url"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	proto "github.com/goverland-labs/feed/protobuf/internalapi"
)

type SubscriberProvider interface {
	GetByID(_ context.Context, id string) (*Subscriber, error)
	Create(_ context.Context, item Subscriber) error
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
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid subscriber ID")
	}

	if req.GetWebhookURL() != "" {
		if _, err := url.ParseRequestURI(req.GetWebhookURL()); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid webhook url")
		}
	}

	sub, err := s.sp.GetByID(ctx, req.GetId())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Msgf("get subscriber: %s", req.GetId())
		return nil, status.Error(codes.InvalidArgument, "invalid subscriber ID")
	}

	if err == nil {
		return &proto.CreateSubscriberResponse{UserID: sub.ID}, nil
	}

	err = s.sp.Create(ctx, Subscriber{
		ID:         req.GetId(),
		WebhookURL: req.GetWebhookURL(),
	})
	if err != nil {
		log.Error().Err(err).Msgf("create subscriber: %s", req.GetId())
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &proto.CreateSubscriberResponse{UserID: req.GetId()}, nil
}

func (s *Server) Update(ctx context.Context, req *proto.UpdateSubscriberRequest) (*emptypb.Empty, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid subscriber ID")
	}

	if req.GetWebhookURL() != "" {
		if _, err := url.ParseRequestURI(req.GetWebhookURL()); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid webhook url")
		}
	}

	_, err := s.sp.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid subscriber ID")
	}

	err = s.sp.Update(ctx, Subscriber{
		ID:         req.GetId(),
		WebhookURL: req.GetWebhookURL(),
	})
	if err != nil {
		log.Error().Err(err).Msgf("update subscriber: %s", req.GetId())
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &emptypb.Empty{}, nil
}
