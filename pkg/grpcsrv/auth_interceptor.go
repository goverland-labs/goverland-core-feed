package grpcsrv

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goverland-labs/feed/internal/subscriber"
)

type SubscriberProvider interface {
	GetByID(_ context.Context, id string) (*subscriber.Subscriber, error)
}

const (
	subscriberIDKey = "subscriber_id"
)

var (
	errWrongSubsciberID = status.Errorf(codes.Unauthenticated, "wrong subscriber identifier")
)

type Auth struct {
	subs SubscriberProvider
}

func NewAuthInterceptor(subs SubscriberProvider) *Auth {
	return &Auth{
		subs: subs,
	}
}

func (a *Auth) AuthAndIdentifyTickerFunc(ctx context.Context) (context.Context, error) {
	md := metautils.ExtractIncoming(ctx)
	requestSubID := md.Get(subscriberIDKey)

	if requestSubID == "" {
		return nil, errWrongSubsciberID
	}

	if _, err := a.subs.GetByID(ctx, requestSubID); err != nil {
		return nil, errWrongSubsciberID
	}

	newCtx := context.WithValue(ctx, subscriber.IDKey, requestSubID)
	return newCtx, nil
}
