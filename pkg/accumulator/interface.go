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
	// UpdateValue updates the payload of the datum.
	UpdateValue([]byte)
	// SetTags sets the tags of the datum, tags are used for conditional forwarding.
	SetTags([]string)
	// Headers returns the headers of the datum.
	Headers() map[string]string
	// ID returns the ID of the datum, id will be used for deduplication.
	ID() string
}

// Accumulator is the interface which can be used to implement the accumulator operation.
type Accumulator interface {
	Accumulate(ctx context.Context, input <-chan Datum, output chan<- Datum)
}

// AccumulatorCreator is the interface which is used to create an Accumulator.
type AccumulatorCreator interface {
	Create() Accumulator
}

// accumulatorFn is a utility type used to convert an Accumulate function to an Accumulator.
type accumulatorFn func(ctx context.Context, input <-chan Datum, output chan<- Datum)

// Accumulate implements the function of accumulate function.
func (af accumulatorFn) Accumulate(ctx context.Context, input <-chan Datum, output chan<- Datum) {
	af(ctx, input, output)
}

// simpleAccumulatorCreator is an implementation of AccumulatorCreator, which creates an Accumulator for the given function.
type simpleAccumulatorCreator struct {
	f func(context.Context, <-chan Datum, chan<- Datum)
}

// Create creates an Accumulator for the given function.
func (s *simpleAccumulatorCreator) Create() Accumulator {
	return accumulatorFn(s.f)
}

// SimpleCreatorWithAccumulateFn creates a simple AccumulatorCreator for the given accumulate function.
func SimpleCreatorWithAccumulateFn(f func(context.Context, <-chan Datum, chan<- Datum)) AccumulatorCreator {
	return &simpleAccumulatorCreator{f: f}
}
