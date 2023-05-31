package subscriber

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

//go:generate mockgen -destination=mocks_test.go -package=subscriber . DataProvider

type DataProvider interface {
	Create(Subscriber) error
	Update(Subscriber) error
	GetByID(string) (*Subscriber, error)
}

type Cacher interface {
	UpsertItem(key string, value *Subscriber)
	GetItem(key string) (*Subscriber, bool)
}

type Service struct {
	repo  DataProvider
	cache Cacher
}

func NewService(r DataProvider, c Cacher) (*Service, error) {
	return &Service{
		repo:  r,
		cache: c,
	}, nil
}

func (s *Service) Create(ctx context.Context, item Subscriber) (*Subscriber, error) {
	sub, err := s.GetByID(ctx, item.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get subscriber: %w", err)
	}

	if err == nil {
		return sub, nil
	}

	err = s.repo.Create(item)
	if err != nil {
		return nil, fmt.Errorf("create subscriber: %w", err)
	}

	go s.cache.UpsertItem(item.ID, &item)

	return &item, err
}

func (s *Service) Update(ctx context.Context, item Subscriber) error {
	_, err := s.GetByID(ctx, item.ID)
	if err != nil {
		return fmt.Errorf("get subscriber: %w", err)
	}

	err = s.repo.Update(item)
	if err != nil {
		return fmt.Errorf("update subscriber: %w", err)
	}

	go s.cache.UpsertItem(item.ID, &item)

	return nil
}

func (s *Service) GetByID(_ context.Context, id string) (*Subscriber, error) {
	if el, ok := s.cache.GetItem(id); ok {
		return el, nil
	}

	sub, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}

	go s.cache.UpsertItem(sub.ID, sub)

	return sub, nil
}
