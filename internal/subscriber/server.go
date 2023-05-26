package subscriber

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"gorm.io/gorm"

	proto "github.com/goverland-labs/feed/protobuf/internalapi"
)

type SubscriberProvider interface {
	GetByID(_ context.Context, id string) (*Subscriber, error)
	Create(_ context.Context, item Subscriber) error
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
		return nil, errors.New("invalid ID")
	}

	if req.GetWebhookURL() != "" {
		if _, err := url.ParseRequestURI(req.GetWebhookURL()); err != nil {
			return nil, fmt.Errorf("invalid webhook url: %w", err)
		}
	}

	sub, err := s.sp.GetByID(ctx, req.GetId())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get subscriber: %w", err)
	}

	if err == nil {
		return &proto.CreateSubscriberResponse{UserID: sub.ID}, nil
	}

	err = s.sp.Create(ctx, Subscriber{
		ID:         req.GetId(),
		WebhookURL: req.GetWebhookURL(),
	})
	if err != nil {
		return nil, fmt.Errorf("create subscriber: %w", err)
	}

	return &proto.CreateSubscriberResponse{UserID: req.GetId()}, nil
}
