package subscriber

import (
	"context"
	"fmt"
)

//go:generate mockgen -destination=mocks_test.go -package=subscriber . DataProvider

type DataProvider interface {
	Create(Subscriber) error
	Update(Subscriber) error
	GetByID(string) (*Subscriber, error)
}

type Service struct {
	repo DataProvider
}

func NewService(r DataProvider) (*Service, error) {
	return &Service{
		repo: r,
	}, nil
}

func (s *Service) Create(_ context.Context, item Subscriber) error {
	err := s.repo.Create(item)
	if err != nil {
		return fmt.Errorf("create subscriber: %w", err)
	}

	return nil
}

func (s *Service) Update(_ context.Context, item Subscriber) error {
	err := s.repo.Update(item)
	if err != nil {
		return fmt.Errorf("update subscriber: %w", err)
	}

	return nil
}

func (s *Service) GetByID(_ context.Context, id string) (*Subscriber, error) {
	sub, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}

	return sub, nil
}
