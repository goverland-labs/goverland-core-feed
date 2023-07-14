package subscription

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(item Subscription) error {
	return r.db.Create(&item).Error
}

func (r *Repo) Delete(item Subscription) error {
	return r.db.Delete(&item).Error
}

func (r *Repo) GetByID(subscriberID, daoID uuid.UUID) (Subscription, error) {
	var res Subscription
	err := r.db.
		Where(&Subscription{
			SubscriberID: subscriberID,
			DaoID:        daoID,
		}).
		First(&res).
		Error

	return res, err
}

// todo: think about getting this elements by chunks

func (r *Repo) GetSubscribers(daoID uuid.UUID) ([]Subscription, error) {
	var res []Subscription
	err := r.db.
		Where(&Subscription{
			DaoID: daoID,
		}).
		Find(&res).
		Error

	return res, err
}
