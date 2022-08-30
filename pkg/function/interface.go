package function

import (
	"context"
	"time"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client interface {
	CloseConn(ctx context.Context) error
	IsReady(ctx context.Context, in *emptypb.Empty) (bool, error)
	MapFn(ctx context.Context, datum *v1.Datum) ([]*v1.Datum, error)
	ReduceFn(ctx context.Context, datumStream []*v1.Datum) ([]*v1.Datum, error)
}

type Datum interface {
	Key() string
	Value() []byte
	EventTime() time.Time
	Watermark() time.Time
}

type Metadata interface {
	IntervalWindow() IntervalWindow
	NonKey() bool
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

type ReduceFunc func(ctx context.Context, reduceCh <-chan Datum, md Metadata) (Messages, error)

func (rf ReduceFunc) HandleDo(ctx context.Context, reduceCh <-chan Datum, md Metadata) (Messages, error) {
	return rf(ctx, reduceCh, md)
}
