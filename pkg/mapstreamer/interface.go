package mapstreamer

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

// MapStreamer is the interface of map stream function implementation.
type MapStreamer interface {
	// MapStream is the function to process each coming message and streams
	// the result back using a channel.
	MapStream(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message)
	MapStreamBatch(ctx context.Context, datumCh <-chan Datum, messageCh chan<- Message)
}

// MapStreamerFunc is a utility type used to convert a function to a MapStreamer.
type MapStreamerFunc func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message)

// MapStream implements the function of map stream function.
func (msf MapStreamerFunc) MapStream(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
	msf(ctx, keys, datum, messageCh)
}
