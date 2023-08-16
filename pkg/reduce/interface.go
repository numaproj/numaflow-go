package reduce

import (
	"context"
	"time"
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

// ReduceHandler is the interface of reduce function implementation.
type ReduceHandler interface {
	HandleDo(ctx context.Context, keys []string, reduceCh <-chan Datum, md Metadata) Messages
}

// ReduceFunc is a utility type used to convert a HandleDo function to a ReduceHandler.
type ReduceFunc func(ctx context.Context, keys []string, reduceCh <-chan Datum, md Metadata) Messages

// HandleDo implements the function of reduce function.
func (rf ReduceFunc) HandleDo(ctx context.Context, keys []string, reduceCh <-chan Datum, md Metadata) Messages {
	return rf(ctx, keys, reduceCh, md)
}
