package feedevent

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/goverland-labs/goverland-core-feed/internal/item"
	"github.com/goverland-labs/goverland-core-feed/protocol/feedpb"
)

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

	err := s.service.Watch(ctx, "feed events", req.LastUpdatedAt.AsTime(), func(entity item.FeedItem) error {
		return stream.Send(&feedpb.FeedItem{
			CreatedAt: timestamppb.New(entity.CreatedAt),
			UpdatedAt: timestamppb.New(entity.UpdatedAt),
			Type:      0,   // TODO resolve type
			Snapshot:  nil, // TODO resolve snapshot
		})
	})
	if err != nil {
		log.Error().
			Str("subscriber", req.SubscriberId).
			Err(err).Msg("error Watch feed events")

		return status.Error(codes.Internal, "internal error")
	}

	return nil
}
