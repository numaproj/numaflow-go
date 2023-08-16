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

// SourceTransformer is the interface of SourceTransformer function implementation.
type SourceTransformer interface {
	// SourceTransformer is the function to process each coming message.
	SourceTransformer(ctx context.Context, keys []string, datum Datum) Messages
}

// SourceTransformFunc is a utility type used to convert a function to a SourceTransformer.
type SourceTransformFunc func(ctx context.Context, keys []string, datum Datum) Messages

// SourceTransformer implements the function of source transformer function.
func (mf SourceTransformFunc) SourceTransformer(ctx context.Context, keys []string, datum Datum) Messages {
	return mf(ctx, keys, datum)
}
