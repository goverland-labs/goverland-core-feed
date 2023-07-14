package grpcsrv

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goverland-labs/core-feed/internal/subscriber"
)

type SubscriberProvider interface {
	GetByID(_ context.Context, id uuid.UUID) (*subscriber.Subscriber, error)
}

const (
	subscriberIDKey = "subscriber_id"
)

var (
	errWrongSubscriberID = status.Errorf(codes.Unauthenticated, "wrong subscriber identifier")
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
		return nil, errWrongSubscriberID
	}

	parsed, err := uuid.Parse(requestSubID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errWrongSubscriberID, err.Error())
	}

	if _, err := a.subs.GetByID(ctx, parsed); err != nil {
		return nil, errWrongSubscriberID
	}

	newCtx := context.WithValue(ctx, subscriber.IDKey, parsed)
	return newCtx, nil
}
