package function

import (
	"context"
	"time"

	functionpb "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client interface {
	CloseConn(ctx context.Context) error
	IsReady(ctx context.Context, in *emptypb.Empty) (bool, error)
	MapFn(ctx context.Context, datum *functionpb.Datum) ([]*functionpb.Datum, error)
	ReduceFn(ctx context.Context, datumStreamCh <-chan *functionpb.Datum) ([]*functionpb.Datum, error)
}

type Datum interface {
	Key() string
	Value() []byte
	EventTime() time.Time
	Watermark() time.Time
}

type Metadata interface {
	IntervalWindow() IntervalWindow
}

type IntervalWindow interface {
	StartTime() time.Time
	EndTime() time.Time
}

// MapHandler is the interface of map function implementation.
type MapHandler interface {
	// HandleDo is the function to process each coming message
	HandleDo(ctx context.Context, datum Datum) (Messages, error)
}

// ReduceHandler is the interface of reduce function implementation.
type ReduceHandler interface {
	HandleDo(ctx context.Context, reduceCh <-chan Datum, md Metadata) (Messages, error) // TODO
}

// MapFunc is utility type used to convert a HandleDo function to a MapHandler.
type MapFunc func(ctx context.Context, datum Datum) (Messages, error)

// HandleDo implements the function of map function.
func (mf MapFunc) HandleDo(ctx context.Context, datum Datum) (Messages, error) {
	return mf(ctx, datum)
}

// ReduceFunc is utility type used to convert a HandleDo function to a ReduceHandler.
type ReduceFunc func(ctx context.Context, reduceCh <-chan Datum, md Metadata) (Messages, error)

// HandleDo implements the function of reduce function.
func (rf ReduceFunc) HandleDo(ctx context.Context, reduceCh <-chan Datum, md Metadata) (Messages, error) {
	return rf(ctx, reduceCh, md)
}
