package sinker

import (
	"context"
	"time"
)

// Datum is the interface of incoming message payload for sink function.
type Datum interface {
	// Keys returns the keys of the message.
	Keys() []string
	// Value returns the payload of the message.
	Value() []byte
	// EventTime returns the event time of the message.
	EventTime() time.Time
	// Watermark returns the watermark of the message.
	Watermark() time.Time
	// ID returns the ID of the message.
	ID() string
	// Headers returns the headers of the message.
	Headers() map[string]string
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
