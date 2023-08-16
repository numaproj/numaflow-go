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

// MapStreamHandler is the interface of map stream function implementation.
type MapStreamHandler interface {
	// HandleDo is the function to process each coming message and streams
	// the result back using a channel.
	HandleDo(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message)
}

// MapStreamFunc is a utility type used to convert a HandleDo function to a MapStreamHandler.
type MapStreamFunc func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message)

// HandleDo implements the function of map stream function.
func (msf MapStreamFunc) HandleDo(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
	msf(ctx, keys, datum, messageCh)
}
