package item

import (
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// Create creates one feed item object
func (r *Repo) Create(item FeedItem) error {
	return r.db.Create(&item).Error
}
