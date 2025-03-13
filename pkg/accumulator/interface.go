package accumulator

import (
	"context"
	"time"
)

// Datum contains methods to get the payload information.
type Datum interface {
	Value() []byte
	EventTime() time.Time
	Watermark() time.Time
	Keys() []string
	UpdateValue([]byte)
	SetTags([]string)
	Headers() map[string]string
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
