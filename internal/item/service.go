package item

import (
	"context"
	"fmt"
)

//go:generate mockgen -destination=mocks_test.go -package=item . DataProvider,Publisher

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

type DataProvider interface {
	Create(item FeedItem) error
}

type Service struct {
	repo   DataProvider
	events Publisher
}

func NewService(r DataProvider, p Publisher) (*Service, error) {
	return &Service{
		repo:   r,
		events: p,
	}, nil
}

func (s *Service) HandleItem(_ context.Context, item FeedItem) error {
	err := s.repo.Create(item)
	if err != nil {
		return fmt.Errorf("can't create item: %w", err)
	}

	// todo: publish core event

	return nil
}
