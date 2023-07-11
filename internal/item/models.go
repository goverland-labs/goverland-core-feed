package item

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Type string

const (
	TypeDao      = "dao"
	TypeProposal = "proposal"
)

type FeedItem struct {
	gorm.Model

	DaoID        string
	ProposalID   string
	DiscussionID string
	Type         Type
	Action       string

	Snapshot json.RawMessage
}

type FeedList struct {
	Items      []FeedItem
	TotalCount int64
}
