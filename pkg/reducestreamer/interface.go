package reducestreamer

import (
	"context"
	"time"
)

// Datum contains methods to get the payload information.
type Datum interface {
	// Value returns the payload of the message.
	Value() []byte
	// EventTime returns the event time of the message.
	EventTime() time.Time
	// Watermark returns the watermark of the message.
	Watermark() time.Time
	// Headers returns the headers of the message.
	Headers() map[string]string
}

// Metadata contains methods to get the metadata for the reduceStream operation.
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

// ReduceStreamerCreator is the interface which is used to create a ReduceStreamer.
type ReduceStreamerCreator interface {
	// Create creates a ReduceStreamer, will be invoked once for every keyed window.
	Create() ReduceStreamer
}

// simpleReducerCreator is an implementation of ReduceStreamerCreator, which creates a ReduceStreamer for the given function.
type simpleReduceStreamerCreator struct {
	f func(ctx context.Context, keys []string, inputCh <-chan Datum, outputCh chan<- Message, md Metadata)
}

// Create creates a Reducer for the given function.
func (s *simpleReduceStreamerCreator) Create() ReduceStreamer {
	return reduceStreamFn(s.f)
}

// SimpleCreatorWithReduceStreamFn creates a simple ReduceStreamerCreator for the given reduceStream function.
func SimpleCreatorWithReduceStreamFn(f func(ctx context.Context, keys []string, inputCh <-chan Datum, outputCh chan<- Message, md Metadata)) ReduceStreamerCreator {
	return &simpleReduceStreamerCreator{f: f}
}

// reduceStreamFn is a utility type used to convert a ReduceStream function to a ReduceStreamer.
type reduceStreamFn func(ctx context.Context, keys []string, inputCh <-chan Datum, outputCh chan<- Message, md Metadata)

// ReduceStream implements the function of ReduceStreamer interface.
func (rf reduceStreamFn) ReduceStream(ctx context.Context, keys []string, inputCh <-chan Datum, outputCh chan<- Message, md Metadata) {
	rf(ctx, keys, inputCh, outputCh, md)
}
