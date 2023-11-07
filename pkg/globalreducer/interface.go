package globalreducer

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

// GlobalReducer is the interface which can be used to implement a session reduce operation.
type GlobalReducer interface {
	GlobalReduce(ctx context.Context, keys []string, input <-chan Datum, output <-chan Messages)
}
