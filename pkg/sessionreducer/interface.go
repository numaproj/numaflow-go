package sessionreducer

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

// SessionReducer is the interface which can be used to implement a session reduce operation.
type SessionReducer interface {
	SessionReduce(ctx context.Context, keys []string, input <-chan Datum) Messages
	Accumulator(ctx context.Context) []byte
	MergeAccumulator(ctx context.Context, accumulator []byte)
}

// CreateSessionReducer is a function which returns a new instance of SessionReducer.
type CreateSessionReducer func() SessionReducer
