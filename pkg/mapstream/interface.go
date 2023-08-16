package mapstream

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

// MapStreamer is the interface of map stream function implementation.
type MapStreamer interface {
	// MapStream is the function to process each coming message and streams
	// the result back using a channel.
	MapStream(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message)
}

// MapStreamerFunc is a utility type used to convert a function to a MapStreamer.
type MapStreamerFunc func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message)

// MapStream implements the function of map stream function.
func (msf MapStreamerFunc) MapStream(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
	msf(ctx, keys, datum, messageCh)
}
