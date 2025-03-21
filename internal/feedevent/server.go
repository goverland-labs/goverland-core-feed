package feedevent

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goverland-labs/goverland-core-feed/internal/item"
	"github.com/goverland-labs/goverland-core-feed/protocol/feedpb"
)

var subscriptionTypesMapping = map[feedpb.FeedItemType]item.Type{
	feedpb.FeedItemType_FEED_ITEM_TYPE_DAO:      item.TypeDao,
	feedpb.FeedItemType_FEED_ITEM_TYPE_PROPOSAL: item.TypeProposal,
	feedpb.FeedItemType_FEED_ITEM_TYPE_DELEGATE: item.TypeDelegate,
}

type Server struct {
	feedpb.UnimplementedFeedEventsServer

	service *Service
}

func NewServer(sp *Service) *Server {
	return &Server{
		service: sp,
	}
}

func (s *Server) EventsSubscribe(req *feedpb.EventsSubscribeRequest, stream grpc.ServerStreamingServer[feedpb.FeedItem]) error {
	ctx := stream.Context()

	log.Info().Msg("events subscribe start")

	subscriberID, err := uuid.Parse(req.SubscriberId)
	if err != nil {
		log.Error().
			Str("subscriber", req.SubscriberId).
			Err(err).Msg("error parse subscriber id")

		return status.Error(codes.InvalidArgument, "invalid subscriber id")
	}

	var fTypes []item.Type
	for _, t := range req.GetSubscriptionTypes() {
		fType, ok := subscriptionTypesMapping[t]
		if !ok {
			log.Error().
				Str("subscriber", req.SubscriberId).
				Str("type", t.String()).
				Msg("unknown feed item type")

			return status.Error(codes.InvalidArgument, "unknown feed item type")
		}

		fTypes = append(fTypes, fType)
	}

	err = s.service.Watch(ctx, subscriberID, fTypes, req.LastUpdatedAt.AsTime(), func(entity item.FeedItem) error {
		feedItem, err := convertToFeedItem(entity)
		if err != nil {
			log.Error().
				Str("subscriber", req.SubscriberId).
				Err(err).Msg("error convert feed item")

			return nil
		}

		return stream.Send(feedItem)
	})
	if err != nil {
		log.Error().
			Str("subscriber", req.SubscriberId).
			Err(err).Msg("error watch feed events")

		return status.Error(codes.Internal, "internal error")
	}

	return nil
}
