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

// ReducerCreator is the interface which is used to create a Reducer.
type ReducerCreator interface {
	// Create creates a Reducer, will be invoked once for every keyed window.
	Create() Reducer
}

// simpleReducerCreator is an implementation of ReducerCreator, which creates a Reducer for the given function.
type simpleReducerCreator struct {
	f func(context.Context, []string, <-chan Datum, Metadata) Messages
}

// Create creates a Reducer for the given function.
func (s *simpleReducerCreator) Create() Reducer {
	return reducerFn(s.f)
}

// SimpleCreatorWithReduceFn creates a simple ReducerCreator for the given reduce function.
func SimpleCreatorWithReduceFn(f func(context.Context, []string, <-chan Datum, Metadata) Messages) ReducerCreator {
	return &simpleReducerCreator{f: f}
}

// Reducer is the interface of reduce function implementation.
type Reducer interface {
	Reduce(ctx context.Context, keys []string, reduceCh <-chan Datum, md Metadata) Messages
}

// reducerFn is a utility type used to convert a Reduce function to a Reducer.
type reducerFn func(ctx context.Context, keys []string, reduceCh <-chan Datum, md Metadata) Messages

// Reduce implements the function of reduce function.
func (rf reducerFn) Reduce(ctx context.Context, keys []string, reduceCh <-chan Datum, md Metadata) Messages {
	return rf(ctx, keys, reduceCh, md)
}
