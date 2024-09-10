package item

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Type string

const (
	TypeDao      Type = "dao"
	TypeProposal Type = "proposal"
	TypeDelegate Type = "delegate"

	None                        TimelineAction = ""
	DaoCreated                  TimelineAction = "dao.created"
	DaoUpdated                  TimelineAction = "dao.updated"
	ProposalCreated             TimelineAction = "proposal.created"
	ProposalUpdated             TimelineAction = "proposal.updated"
	ProposalVotingStartsSoon    TimelineAction = "proposal.voting.starts_soon"
	ProposalVotingEndsSoon      TimelineAction = "proposal.voting.ends_soon"
	ProposalVotingStarted       TimelineAction = "proposal.voting.started"
	ProposalVotingQuorumReached TimelineAction = "proposal.voting.quorum_reached"
	ProposalVotingEnded         TimelineAction = "proposal.voting.ended"
	DelegateCreateProposal      TimelineAction = "delegate.proposal.created"
)

type FeedItem struct {
	ID          uuid.UUID `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	TriggeredAt time.Time      `gorm:"index"`

	DaoID        uuid.UUID
	ProposalID   string
	DiscussionID string
	Type         Type
	Action       TimelineAction

	Snapshot json.RawMessage
	Timeline Timeline `gorm:"serializer:json"`
}

type Timeline []TimelineItem

func (t *Timeline) AddUniqueAction(createdAt time.Time, action TimelineAction) (isNew bool) {
	if *t == nil {
		*t = make(Timeline, 0, 1)
	}

	for i := range *t {
		if (*t)[i].Action.Equals(action) {
			return false
		}
	}

	*t = append(*t, TimelineItem{
		CreatedAt: createdAt,
		Action:    action,
	})

	return true
}

func (t *Timeline) AddNonUniqueAction(createdAt time.Time, action TimelineAction) {
	if *t == nil {
		*t = make(Timeline, 0, 1)
	}

	*t = append(*t, TimelineItem{
		CreatedAt: createdAt,
		Action:    action,
	})
}

func (t *Timeline) ContainsAction(action TimelineAction) bool {
	if t == nil || len(*t) == 0 {
		return false
	}

	for _, item := range *t {
		if item.Action.Equals(action) {
			return true
		}
	}

	return false
}

func (t *Timeline) Sort() {
	if t == nil || len(*t) == 0 {
		return
	}

	sort.SliceStable(*t, func(i, j int) bool {
		return (*t)[i].CreatedAt.Before((*t)[j].CreatedAt)
	})
}

func (t *Timeline) LastAction() TimelineAction {
	if t == nil {
		return None
	}

	if len(*t) == 0 {
		return None
	}

	return (*t)[len(*t)-1].Action
}

type TimelineItem struct {
	CreatedAt time.Time      `json:"created_at"`
	Action    TimelineAction `json:"action"`
}

type TimelineAction string

func (a TimelineAction) Equals(action TimelineAction) bool {
	return strings.EqualFold(string(a), string(action))
}

type FeedList struct {
	Items      []FeedItem
	TotalCount int64
}
