package reduceStreamer

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

// ReduceStreamer is the interface of reduceStream function implementation.
type ReduceStreamer interface {
	ReduceStream(ctx context.Context, keys []string, inputCh <-chan Datum, outputCh chan<- Message, md Metadata)
}

// ReducerFunc is a utility type used to convert a ReduceStream function to a ReduceStreamer.
type ReducerFunc func(ctx context.Context, keys []string, inputCh <-chan Datum, outputCh chan<- Message, md Metadata)

// ReduceStream implements the function of reduce function.
func (rf ReducerFunc) ReduceStream(ctx context.Context, keys []string, inputCh <-chan Datum, outputCh chan<- Message, md Metadata) {
	rf(ctx, keys, inputCh, outputCh, md)
}