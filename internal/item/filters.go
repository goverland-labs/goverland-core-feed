package item

import (
	"gorm.io/gorm"
)

type Filter interface {
	Apply(*gorm.DB) *gorm.DB
}

type PageFilter struct {
	Offset int
	Limit  int
}

func (f PageFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Offset(f.Offset).Limit(f.Limit)
}

type DaoIDFilter struct {
	ID string
}

func (f DaoIDFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("dao_id = ?", f.ID)
}

type TypeFilter struct {
	Types []string
}

func (f TypeFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("type IN ?", f.Types)
}

type ActionFilter struct {
	Actions []string
}

func (f ActionFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("action IN ?", f.Actions)
}

type OrderByCreatedFilter struct {
}

func (f OrderByCreatedFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Order("created_at desc")
}
