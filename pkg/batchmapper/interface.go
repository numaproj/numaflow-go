package batchmapper

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
	// Id returns the unique ID set for the given message
	ID() string
	// Keys returns the keys associated with a given datum
	Keys() []string
}

// BatchMapper is the interface for a Batch Map mode where the user is given a list
// of messages, and they return the consolidated response for all of them together.
type BatchMapper interface {
	// BatchMap is the function which processes a list of input messages
	BatchMap(ctx context.Context, datumStreamCh <-chan Datum) BatchResponses
}

// BatchMapperFunc is a utility type used to convert a batch map function to a BatchMapper.
type BatchMapperFunc func(ctx context.Context, datumStreamCh <-chan Datum) BatchResponses

// BatchMap implements the functionality of BatchMap function.
func (mf BatchMapperFunc) BatchMap(ctx context.Context, datumStreamCh <-chan Datum) BatchResponses {
	return mf(ctx, datumStreamCh)
}
