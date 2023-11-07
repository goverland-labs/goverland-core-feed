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
	IDs []string
}

func (f DaoIDFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("dao_id in ?", f.IDs)
}

type TypeFilter struct {
	Types []string
}

func (f TypeFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("type IN ?", f.Types)
}

type ActiveFilter struct {
	IsActive bool
}

func (f ActiveFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.IsActive {
		return db.Where("to_timestamp((snapshot->'end')::double precision) >= now()")
	}

	return db.Where("to_timestamp((snapshot->'end')::double precision) < now()")
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

type OrderByTriggeredFilter struct {
}

func (f OrderByTriggeredFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Order("triggered_at desc")
}
