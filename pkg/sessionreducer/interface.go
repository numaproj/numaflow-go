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
	Aggregator(ctx context.Context) []byte
	MergeAggregator(ctx context.Context, aggregator []byte)
}

type CreateSessionReducer func() SessionReducer
