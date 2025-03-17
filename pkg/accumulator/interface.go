package accumulator

import (
	"context"
	"time"
)

// Datum contains methods to get the payload information.
type Datum interface {
	// Value returns the payload of the datum.
	Value() []byte
	// EventTime returns the event time of the datum.
	EventTime() time.Time
	// Watermark returns the watermark of the datum.
	Watermark() time.Time
	// Keys returns the keys of the datum.
	Keys() []string
	// Headers returns the headers of the datum.
	Headers() map[string]string
	// ID returns the ID of the datum, id will be used for deduplication.
	ID() string
}

// Accumulator is the interface which can be used to implement the accumulator operation.
type Accumulator interface {
	// Accumulate can read unordered from the input stream and emit the ordered data to the output stream.
	// Once the watermark (WM) of the output stream progresses, the data in WAL until that WM will be garbage collected.
	// NOTE: A message can be silently dropped if need be, and it will be cleared from the WAL when the WM progresses.
	Accumulate(ctx context.Context, input <-chan Datum, output chan<- Message)
}

// AccumulatorCreator is the interface which is used to create an Accumulator.
type AccumulatorCreator interface {
	// Create is called for every key and will be closed after the keyed stream is idle for the timeout duration.
	Create() Accumulator
}

// accumulatorFn is a utility type used to convert an Accumulate function to an Accumulator.
type accumulatorFn func(ctx context.Context, input <-chan Datum, output chan<- Message)

// Accumulate implements the function of accumulate function.
func (af accumulatorFn) Accumulate(ctx context.Context, input <-chan Datum, output chan<- Message) {
	af(ctx, input, output)
}

// simpleAccumulatorCreator is an implementation of AccumulatorCreator, which creates an Accumulator for the given function.
type simpleAccumulatorCreator struct {
	f func(context.Context, <-chan Datum, chan<- Message)
}

// Create creates an Accumulator for the given function.
func (s *simpleAccumulatorCreator) Create() Accumulator {
	return accumulatorFn(s.f)
}

// SimpleCreatorWithAccumulateFn creates a simple AccumulatorCreator for the given accumulate function.
func SimpleCreatorWithAccumulateFn(f func(context.Context, <-chan Datum, chan<- Message)) AccumulatorCreator {
	return &simpleAccumulatorCreator{f: f}
}
