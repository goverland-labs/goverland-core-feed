package item

import (
	"fmt"

	"gorm.io/gorm"
)

type Direction string

const (
	DirectionAsc  Direction = "asc"
	DirectionDesc Direction = "desc"
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

type OrderByTriggeredFilter struct {
}

func (f OrderByTriggeredFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Order("triggered_at desc")
}

type SkipSpammed struct {
}

func (f SkipSpammed) Apply(db *gorm.DB) *gorm.DB {
	var (
		dummy FeedItem
		_     = dummy.Snapshot // spam flag
	)

	return db.Where(`snapshot->>'spam' != 'true'`)
}

type SkipCanceled struct {
}

func (f SkipCanceled) Apply(db *gorm.DB) *gorm.DB {
	var (
		dummy FeedItem
		_     = dummy.Snapshot // state
	)

	return db.Where(`snapshot->>'state' != 'canceled'`)
}

type SortedByActuality struct {
}

func (f SortedByActuality) Apply(db *gorm.DB) *gorm.DB {
	var (
		dummy FeedItem
		_     = dummy.CreatedAt
		_     = dummy.Snapshot // state
	)

	return db.Order(`
				array_position(array [
					'active',
					'pending',
					'succeeded',
					'failed',
					'defeated',
					'canceled'
				], snapshot->>'state'), 
				created_at desc`)
}

type SortedByCreated struct {
	Direction Direction
}

func (f SortedByCreated) Apply(db *gorm.DB) *gorm.DB {
	var (
		dummy FeedItem
		_     = dummy.CreatedAt
	)

	return db.Order(fmt.Sprintf("created_at %s", f.Direction))
}
