package subscriber

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=mocks_test.go -package=subscriber . DataProvider

type DataProvider interface {
	Create(*Subscriber) error
	Update(*Subscriber) error
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

func (s *Service) Create(ctx context.Context, webhookURL string) (*Subscriber, error) {
	subID, err := s.generateSubscriberID(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate subscriber id: %w", err)
	}

	item := &Subscriber{
		ID:         subID,
		WebhookURL: webhookURL,
	}
	err = s.repo.Create(item)
	if err != nil {
		return nil, fmt.Errorf("create subscriber: %w", err)
	}

	go s.cache.UpsertItem(item.ID, item)

	return item, err
}

func (s *Service) generateSubscriberID(ctx context.Context) (string, error) {
	subID := uuid.New().String()
	_, err := s.GetByID(ctx, subID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return subID, nil
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("get subscriber: %w", err)
	}

	return s.generateSubscriberID(ctx)
}

func (s *Service) Update(ctx context.Context, item Subscriber) error {
	sub, err := s.GetByID(ctx, item.ID)
	if err != nil {
		return fmt.Errorf("get subscriber: %w", err)
	}

	sub.WebhookURL = item.WebhookURL
	err = s.repo.Update(sub)
	if err != nil {
		return fmt.Errorf("update subscriber: %w", err)
	}

	go s.cache.UpsertItem(item.ID, sub)

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
