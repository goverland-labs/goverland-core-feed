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

type Action string

const (
	actionCreated = "created"
	actionUpdated = "updated"
)

type FeedItem struct {
	gorm.Model

	Type   Type
	TypeID string
	Action Action

	Snapshot json.RawMessage
}
