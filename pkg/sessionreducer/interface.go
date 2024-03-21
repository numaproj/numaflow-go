package sessionreducer

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

// SessionReducer is the interface which can be used to implement a session reduce operation.
type SessionReducer interface {
	// SessionReduce applies a session reduce function to a request stream and streams the results.
	SessionReduce(ctx context.Context, keys []string, inputCh <-chan Datum, outputCh chan<- Message)
	// Accumulator returns the accumulator for the session reducer, will be invoked when this session is merged
	// with another session.
	Accumulator(ctx context.Context) []byte
	// MergeAccumulator merges the accumulator for the session reducer, will be invoked when another session is merged
	// with this session.
	MergeAccumulator(ctx context.Context, accumulator []byte)
}

// SessionReducerCreator is the interface which can be used to create a session reducer.
type SessionReducerCreator interface {
	// Create creates a session reducer, will be invoked once for every keyed window.
	Create() SessionReducer
}
