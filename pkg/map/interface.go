package _map

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

// MapHandler is the interface of map function implementation.
type MapHandler interface {
	// HandleDo is the function to process each coming message.
	HandleDo(ctx context.Context, keys []string, datum Datum) Messages
}

// MapFunc is a utility type used to convert a HandleDo function to a MapHandler.
type MapFunc func(ctx context.Context, keys []string, datum Datum) Messages

// HandleDo implements the function of map function.
func (mf MapFunc) HandleDo(ctx context.Context, keys []string, datum Datum) Messages {
	return mf(ctx, keys, datum)
}
