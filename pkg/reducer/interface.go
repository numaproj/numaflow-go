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
	Reduce(ctx context.Context, keys []string, inputCh <-chan Datum, outputCh chan<- Message, md Metadata)
}

// ReducerFunc is a utility type used to convert a Reduce function to a Reducer.
type ReducerFunc func(ctx context.Context, keys []string, inputCh <-chan Datum, outputCh chan<- Message, md Metadata)

// Reduce implements the function of reduce function.
func (rf ReducerFunc) Reduce(ctx context.Context, keys []string, inputCh <-chan Datum, outputCh chan<- Message, md Metadata) {
	rf(ctx, keys, inputCh, outputCh, md)
}
