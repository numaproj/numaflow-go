package sink

import (
	"context"
	"time"
)

// Datum is the interface of incoming message payload for sink function.
type Datum interface {
	Keys() []string
	Value() []byte
	EventTime() time.Time
	Watermark() time.Time
	ID() string
}

// Sinker is the interface of sink function implementation.
type Sinker interface {
	// Sink is the function to process a list of incoming messages
	Sink(ctx context.Context, datumStreamCh <-chan Datum) Responses
}

// SinkerFunc is utility type used to convert a Sink function to a Sinker.
type SinkerFunc func(ctx context.Context, datumStreamCh <-chan Datum) Responses

// Sink implements the function of sink function.
func (sf SinkerFunc) Sink(ctx context.Context, datumStreamCh <-chan Datum) Responses {
	return sf(ctx, datumStreamCh)
}
