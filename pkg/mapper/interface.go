package mapper

import (
	"context"
	"time"
)

// Datum contains methods to get the payload information.
type Datum interface {
	Value() []byte
	EventTime() time.Time
	Watermark() time.Time
	Headers() map[string]string
}

// Mapper is the interface of map function implementation.
type Mapper interface {
	// Map is the function to process each coming message.
	Map(ctx context.Context, keys []string, datum Datum) Messages
}

// MapperFunc is a utility type used to convert a map function to a Mapper.
type MapperFunc func(ctx context.Context, keys []string, datum Datum) Messages

// Map implements the function of map function.
func (mf MapperFunc) Map(ctx context.Context, keys []string, datum Datum) Messages {
	return mf(ctx, keys, datum)
}
