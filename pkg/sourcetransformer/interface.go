package sourcetransformer

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

// MapTHandler is the interface of mapT function implementation.
type MapTHandler interface {
	// HandleDo is the function to process each coming message.
	HandleDo(ctx context.Context, keys []string, datum Datum) Messages
}

// MapTFunc is a utility type used to convert a HandleDo function to a MapTHandler.
type MapTFunc func(ctx context.Context, keys []string, datum Datum) Messages

// HandleDo implements the function of mapT function.
func (mf MapTFunc) HandleDo(ctx context.Context, keys []string, datum Datum) Messages {
	return mf(ctx, keys, datum)
}
