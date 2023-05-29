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

const (
	actionCreated = "created"
	actionUpdated = "updated"
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
