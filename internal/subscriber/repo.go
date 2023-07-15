package subscriber

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repo struct {
	conn *gorm.DB
}

func NewRepo(conn *gorm.DB) *Repo {
	return &Repo{conn: conn}
}

func (r *Repo) Create(item *Subscriber) error {
	return r.conn.Create(item).Error
}

func (r *Repo) Update(item *Subscriber) error {
	return r.conn.Save(item).Error
}

func (r *Repo) GetByID(id uuid.UUID) (*Subscriber, error) {
	var sub Subscriber

	err := r.conn.
		Where(Subscriber{ID: id}).
		First(&sub).
		Error

	if err != nil {
		return nil, fmt.Errorf("get subscriber by id #%s: %w", id, err)
	}

	return &sub, nil
}
