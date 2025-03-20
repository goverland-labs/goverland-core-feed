package item

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var emptyID uuid.UUID

type Repo struct {
	conn *gorm.DB
}

func NewRepo(conn *gorm.DB) *Repo {
	return &Repo{conn: conn}
}

func (r *Repo) Save(item *FeedItem) error {
	var (
		_ = item.DaoID
		_ = item.ProposalID
		_ = item.Type
		_ = item.Action
	)

	if item.ID == emptyID {
		item.ID = uuid.New()
	}

	conn := r.conn
	// there is no unique key for delegate type
	if item.Type != TypeDelegate {
		conn = conn.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "dao_id"},
				{Name: "proposal_id"},
				{Name: "type"},
				{Name: "action"},
			},
			TargetWhere: clause.Where{Exprs: []clause.Expression{
				clause.Neq{Column: "type", Value: TypeDelegate}},
			},
			UpdateAll: true,
		})
	}

	err := conn.
		Create(item).
		Error

	return err
}

func (r *Repo) GetDaoItem(id uuid.UUID) (*FeedItem, error) {
	var (
		item FeedItem
		_    = item.DaoID
	)

	err := r.conn.Where("dao_id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Repo) GetProposalItem(id string) (*FeedItem, error) {
	var (
		item FeedItem
		_    = item.ProposalID
	)

	err := r.conn.Where("proposal_id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Repo) GetLastItems(subscriberID string, lastUpdatedAt time.Time, limit int) ([]FeedItem, error) {
	var feedItems []FeedItem

	err := r.conn.
		Joins("JOIN subscriptions ON subscriptions.dao_id = feed_items.dao_id").
		Where("subscriptions.subscriber_id = ?", subscriberID).
		Where("feed_items.updated_at > ?", lastUpdatedAt).
		Order("feed_items.updated_at asc").
		Limit(limit).
		Find(&feedItems).Error
	if err != nil {
		return nil, err
	}

	return feedItems, nil
}

func (r *Repo) GetByFilters(filters []Filter) (FeedList, error) {
	db := r.conn.Model(&FeedItem{})
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
