package feedevent

import (
	"encoding/json"
	"errors"
	"time"

	pevents "github.com/goverland-labs/goverland-platform-events/events/core"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/goverland-labs/goverland-core-feed/protocol/feedpb"

	"github.com/goverland-labs/goverland-core-feed/internal/item"
)

var typesMapping = map[item.Type]feedpb.FeedItemType{
	item.TypeDao:      feedpb.FeedItemType_FEED_ITEM_TYPE_DAO,
	item.TypeProposal: feedpb.FeedItemType_FEED_ITEM_TYPE_PROPOSAL,
	item.TypeDelegate: feedpb.FeedItemType_FEED_ITEM_TYPE_DELEGATE,
}

var snapshotConverters = map[item.Type]func(item.FeedItem) (any, error){
	item.TypeDao:      daoSnapshotConverter,
	item.TypeProposal: proposalSnapshotConverter,
	item.TypeDelegate: delegateSnapshotConverter,
}

func convertToFeedItem(fItem item.FeedItem) (*feedpb.FeedItem, error) {
	fItemType, ok := typesMapping[fItem.Type]
	if !ok {
		return nil, errors.New("unknown feed item type")
	}

	converter, ok := snapshotConverters[fItem.Type]
	if !ok {
		return nil, errors.New("snapshot converter not found")
	}

	converted, err := converter(fItem)
	if err != nil {
		return nil, err
	}

	feedItem := &feedpb.FeedItem{
		CreatedAt: timestamppb.New(fItem.CreatedAt),
		UpdatedAt: timestamppb.New(fItem.UpdatedAt),
		Type:      fItemType,
	}

	// TODO: move to opaque grpc to avoid type switch
	switch v := converted.(type) {
	case *feedpb.FeedItem_Dao:
		feedItem.Snapshot = v
	case *feedpb.FeedItem_Proposal:
		feedItem.Snapshot = v
	case *feedpb.FeedItem_Delegate:
		feedItem.Snapshot = v
	}

	return feedItem, nil
}

func daoSnapshotConverter(fItem item.FeedItem) (any, error) {
	var dPayload pevents.DaoPayload
	err := json.Unmarshal(fItem.Snapshot, &dPayload)
	if err != nil {
		return nil, err
	}

	var popularityIndex float64
	if dPayload.PopularityIndex != nil {
		popularityIndex = *dPayload.PopularityIndex
	}

	var timeline []*feedpb.Timeline
	for _, tItem := range fItem.Timeline {
		timeline = append(timeline, &feedpb.Timeline{
			Action:    string(tItem.Action),
			CreatedAt: timestamppb.New(tItem.CreatedAt),
		})
	}

	return &feedpb.FeedItem_Dao{
		Dao: &feedpb.DAO{
			CreatedAt:       timestamppb.New(dPayload.CreatedAt),
			InternalId:      dPayload.ID.String(),
			OriginalId:      dPayload.Alias,
			Name:            dPayload.Name,
			Avatar:          dPayload.Avatar,
			PopularityIndex: popularityIndex,
			Verified:        dPayload.Verified,
			Timeline:        timeline,
		},
	}, nil
}

func proposalSnapshotConverter(fItem item.FeedItem) (any, error) {
	var dPayload pevents.ProposalPayload
	err := json.Unmarshal(fItem.Snapshot, &dPayload)
	if err != nil {
		return nil, err
	}

	var timeline []*feedpb.Timeline
	for _, tItem := range fItem.Timeline {
		timeline = append(timeline, &feedpb.Timeline{
			Action:    string(tItem.Action),
			CreatedAt: timestamppb.New(tItem.CreatedAt),
		})
	}

	return &feedpb.FeedItem_Proposal{
		Proposal: &feedpb.Proposal{
			CreatedAt:     timestamppb.New(time.Unix(int64(dPayload.Created), 0)),
			Id:            dPayload.ID,
			DaoInternalId: dPayload.DaoID.String(),
			Author:        dPayload.Author,
			Title:         dPayload.Title,
			State:         dPayload.State,
			Spam:          dPayload.Spam,
			Timeline:      timeline,
			Type:          dPayload.Type,
			Privacy:       dPayload.Privacy,
			Choices:       dPayload.Choices,
			VoteStart:     timestamppb.New(time.Unix(int64(dPayload.Start), 0)),
			VoteEnd:       timestamppb.New(time.Unix(int64(dPayload.End), 0)),
		},
	}, nil
}

func delegateSnapshotConverter(fItem item.FeedItem) (any, error) {
	var dPayload pevents.DelegatePayload
	err := json.Unmarshal(fItem.Snapshot, &dPayload)
	if err != nil {
		return nil, err
	}

	var dueDate *timestamppb.Timestamp
	if dPayload.DueDate != nil {
		dueDate = timestamppb.New(*dPayload.DueDate)
	}

	return &feedpb.FeedItem_Delegate{
		Delegate: &feedpb.Delegate{
			AddressFrom:   dPayload.Initiator,
			AddressTo:     dPayload.Delegator,
			DaoInternalId: dPayload.DaoID.String(),
			ProposalId:    dPayload.ProposalID,
			Action:        string(fItem.Action),
			DueDate:       dueDate,
		},
	}, nil
}
