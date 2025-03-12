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

// Accumulator is the interface which can be used to implement a session reduce operation.
type Accumulator interface {
	Accumulate(ctx context.Context, input <-chan Datum, output <-chan Datum)
}
