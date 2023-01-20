package function

import (
	"context"
	"time"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Datum contains methods to get the payload information.
type Datum interface {
	Value() []byte
	EventTime() time.Time
	Watermark() time.Time
}

// Metadata contains methods to get the metadata for the reduce operation.
type Metadata interface {
	IntervalWindow() IntervalWindow
}

// IntervalWindow contains methods to get the information for a given interval window.
type IntervalWindow interface {
	StartTime() time.Time
	EndTime() time.Time
}

// Client contains methods to call a gRPC client.
type Client interface {
	CloseConn(ctx context.Context) error
	IsReady(ctx context.Context, in *emptypb.Empty) (bool, error)
	MapFn(ctx context.Context, datum *functionpb.Datum) ([]*functionpb.Datum, error)
	MapTFn(ctx context.Context, datum *functionpb.Datum) ([]*functionpb.Datum, error)
	ReduceFn(ctx context.Context, datumStreamCh <-chan *functionpb.Datum) ([]*functionpb.Datum, error)
}

// MapHandler is the interface of map function implementation.
type MapHandler interface {
	// HandleDo is the function to process each coming message.
	HandleDo(ctx context.Context, key string, datum Datum) Messages
}

// MapTHandler is the interface of mapT function implementation.
type MapTHandler interface {
	// HandleDo is the function to process each coming message.
	HandleDo(ctx context.Context, key string, datum Datum) MessageTs
}

// ReduceHandler is the interface of reduce function implementation.
type ReduceHandler interface {
	HandleDo(ctx context.Context, key string, reduceCh <-chan Datum, md Metadata) Messages
}

// MapFunc is utility type used to convert a HandleDo function to a MapHandler.
type MapFunc func(ctx context.Context, key string, datum Datum) Messages

// HandleDo implements the function of map function.
func (mf MapFunc) HandleDo(ctx context.Context, key string, datum Datum) Messages {
	return mf(ctx, key, datum)
}

// MapTFunc is utility type used to convert a HandleDo function to a MapTHandler.
type MapTFunc func(ctx context.Context, key string, datum Datum) MessageTs

// HandleDo implements the function of mapT function.
func (mf MapTFunc) HandleDo(ctx context.Context, key string, datum Datum) MessageTs {
	return mf(ctx, key, datum)
}

// ReduceFunc is utility type used to convert a HandleDo function to a ReduceHandler.
type ReduceFunc func(ctx context.Context, key string, reduceCh <-chan Datum, md Metadata) Messages

// HandleDo implements the function of reduce function.
func (rf ReduceFunc) HandleDo(ctx context.Context, key string, reduceCh <-chan Datum, md Metadata) Messages {
	return rf(ctx, key, reduceCh, md)
}
