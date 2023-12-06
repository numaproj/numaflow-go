package reducer

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

// Reducer is the interface of reduce function implementation.
type Reducer interface {
	Reduce(ctx context.Context, keys []string, inputCh <-chan Datum, md Metadata) Messages
}

// ReducerFunc is a utility type used to convert a Reduce function to a Reducer.
type ReducerFunc func(ctx context.Context, keys []string, inputCh <-chan Datum, md Metadata) Messages

// Reduce implements the function of reduce function.
func (rf ReducerFunc) Reduce(ctx context.Context, keys []string, inputCh <-chan Datum, md Metadata) Messages {
	return rf(ctx, keys, inputCh, md)
}
