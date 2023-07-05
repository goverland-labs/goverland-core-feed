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

func (r *Repo) GetByFilters(filters []Filter) (FeedList, error) {
	db := r.db.Model(&FeedItem{})
	for _, f := range filters {
		if _, ok := f.(PageFilter); ok {
			continue
		}
		db = f.Apply(db)
	}

	var cnt int64
	err := db.Count(&cnt).Error
	if err != nil {
		return FeedList{}, err
	}

	for _, f := range filters {
		if _, ok := f.(PageFilter); ok {
			db = f.Apply(db)
		}
	}

	var list []FeedItem
	err = db.Find(&list).Error
	if err != nil {
		return FeedList{}, err
	}

	return FeedList{
		Items:      list,
		TotalCount: cnt,
	}, nil
}
