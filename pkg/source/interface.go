package source

import (
	"context"
	sourcepb "github.com/numaproj/numaflow-go/pkg/apis/proto/source/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

// Datum contains methods to get the payload information.
type Datum interface {
	Value() []byte
	EventTime() time.Time
	Watermark() time.Time
}

// Client contains methods to call a gRPC client.
type Client interface {
	CloseConn(ctx context.Context) error
	IsReady(ctx context.Context, in *emptypb.Empty) (bool, error)
	TransformFn(ctx context.Context, datum *sourcepb.Datum) ([]*sourcepb.Datum, error)
}

// TransformHandler is the interface of transform function implementation.
type TransformHandler interface {
	// HandleTransform is the function to process each coming message
	HandleTransform(ctx context.Context, key string, datum Datum) Messages
}

// TransformFunc is utility type used to convert a HandleTransform function to a TransformHandler.
type TransformFunc func(ctx context.Context, key string, datum Datum) Messages

// HandleTransform implements the function of transform function.
func (tf TransformFunc) HandleTransform(ctx context.Context, key string, datum Datum) Messages {
	return tf(ctx, key, datum)
}
