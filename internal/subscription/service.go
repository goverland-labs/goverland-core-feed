package subscription

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=mocks_test.go -package=subscription . DataProvider,Cacher

type DataProvider interface {
	Create(Subscription) error
	Delete(Subscription) error
	GetByID(uuid.UUID, uuid.UUID) (Subscription, error)
	GetSubscribers(daoID uuid.UUID) ([]Subscription, error)
}

type Cacher interface {
	AddItems(string, ...uuid.UUID)
	RemoveItem(string, uuid.UUID)
	UpdateItems(string, ...uuid.UUID)
	GetItems(string) ([]uuid.UUID, bool)
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

func (s *Service) Subscribe(_ context.Context, item Subscription) (*Subscription, error) {
	sub, err := s.repo.GetByID(item.SubscriberID, item.DaoID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get subscription: %w", err)
	}

	if err == nil {
		return &sub, nil
	}

	err = s.repo.Create(item)
	if err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}

	go s.cache.AddItems(item.DaoID.String(), item.SubscriberID)

	return &item, err
}

func (s *Service) Unsubscribe(_ context.Context, item Subscription) error {
	sub, err := s.repo.GetByID(item.SubscriberID, item.DaoID)
	if err != nil {
		return fmt.Errorf("get subscription: %w", err)
	}

	err = s.repo.Delete(sub)
	if err != nil {
		return fmt.Errorf("delete scubscription[%s - %s]: %w", item.SubscriberID, item.DaoID, err)
	}

	go s.cache.RemoveItem(item.DaoID.String(), item.SubscriberID)

	return nil
}

func (s *Service) GetSubscribers(_ context.Context, daoID uuid.UUID) ([]uuid.UUID, error) {
	if list, ok := s.cache.GetItems(daoID.String()); ok {
		return list, nil
	}

	data, err := s.repo.GetSubscribers(daoID)
	if err != nil {
		return nil, fmt.Errorf("get subscribers: %w", err)
	}

	response := make([]uuid.UUID, len(data))
	for i, sub := range data {
		response[i] = sub.SubscriberID
	}

	go s.cache.UpdateItems(daoID.String(), response...)

	return response, nil
}
