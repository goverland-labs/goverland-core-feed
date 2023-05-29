package subscriber

import (
	"fmt"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(item Subscriber) error {
	return r.db.Create(&item).Error
}

func (r *Repo) Update(item Subscriber) error {
	return r.db.Save(&item).Error
}

func (r *Repo) GetByID(id string) (*Subscriber, error) {
	t := Subscriber{ID: id}
	request := r.db.Take(&t)
	if err := request.Error; err != nil {
		return nil, fmt.Errorf("get subscriber by id #%s: %w", id, err)
	}

	return &t, nil
}
