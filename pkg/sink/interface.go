package sink

import (
	"context"
	"time"
)

// Datum is the interface of incoming message payload for user defined sink function.
type Datum interface {
	Keys() []string
	Value() []byte
	EventTime() time.Time
	Watermark() time.Time
	ID() string
}

// SinkHandler is the interface of sink function implementation.
type SinkHandler interface {
	// HandleDo is the function to process a list of incoming messages
	HandleDo(ctx context.Context, datumStreamCh <-chan Datum) Responses
}

// SinkFunc is utility type used to convert a HandleDo function to a SinkHandler.
type SinkFunc func(ctx context.Context, datumStreamCh <-chan Datum) Responses

// HandleDo implements the function of sink function.
func (sf SinkFunc) HandleDo(ctx context.Context, datumStreamCh <-chan Datum) Responses {
	return sf(ctx, datumStreamCh)
}
