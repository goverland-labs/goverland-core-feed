package subscription

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

//go:generate mockgen -destination=mocks_test.go -package=subscription . DataProvider

type DataProvider interface {
	Create(Subscription) error
	Delete(Subscription) error
	GetByID(string, string) (Subscription, error)
}

type Service struct {
	repo DataProvider
}

func NewService(r DataProvider) (*Service, error) {
	return &Service{
		repo: r,
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

	return nil
}
