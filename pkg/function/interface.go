package function

import (
	"context"
	"time"
)

// Datum contains methods to get the payload information.
type Datum interface {
	Value() []byte
	EventTime() time.Time
	Watermark() time.Time
	Metadata() DatumMetadata
}

// DatumMetadata contains methods to get the metadata information for datum.
type DatumMetadata interface {
	// ID returns the ID of the datum.
	ID() string
	// NumDelivered returns the number of times the datum has been delivered.
	NumDelivered() uint64
}

// Metadata contains methods to get the metadata for the reduce operation.
type Metadata interface {
	IntervalWindow() IntervalWindow
}

// IntervalWindow contains methods to get the information for a given interval window.
type IntervalWindow interface {
	StartTime() time.Time
	EndTime() time.Time
}

// MapHandler is the interface of map function implementation.
type MapHandler interface {
	// HandleDo is the function to process each coming message.
	HandleDo(ctx context.Context, keys []string, datum Datum) Messages
}

// MapStreamHandler is the interface of map stream function implementation.
type MapStreamHandler interface {
	// HandleDo is the function to process each coming message and streams
	// the result back using a channel.
	HandleDo(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message)
}

// MapTHandler is the interface of mapT function implementation.
type MapTHandler interface {
	// HandleDo is the function to process each coming message.
	HandleDo(ctx context.Context, keys []string, datum Datum) MessageTs
}

// ReduceHandler is the interface of reduce function implementation.
type ReduceHandler interface {
	HandleDo(ctx context.Context, keys []string, reduceCh <-chan Datum, md Metadata) Messages
}

// MapFunc is a utility type used to convert a HandleDo function to a MapHandler.
type MapFunc func(ctx context.Context, keys []string, datum Datum) Messages

// HandleDo implements the function of map function.
func (mf MapFunc) HandleDo(ctx context.Context, keys []string, datum Datum) Messages {
	return mf(ctx, keys, datum)
}

// MapStreamFunc is a utility type used to convert a HandleDo function to a MapStreamHandler.
type MapStreamFunc func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message)

// HandleDo implements the function of map stream function.
func (msf MapStreamFunc) HandleDo(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
	msf(ctx, keys, datum, messageCh)
}

// MapTFunc is a utility type used to convert a HandleDo function to a MapTHandler.
type MapTFunc func(ctx context.Context, keys []string, datum Datum) MessageTs

// HandleDo implements the function of mapT function.
func (mf MapTFunc) HandleDo(ctx context.Context, keys []string, datum Datum) MessageTs {
	return mf(ctx, keys, datum)
}

// ReduceFunc is a utility type used to convert a HandleDo function to a ReduceHandler.
type ReduceFunc func(ctx context.Context, keys []string, reduceCh <-chan Datum, md Metadata) Messages

// HandleDo implements the function of reduce function.
func (rf ReduceFunc) HandleDo(ctx context.Context, keys []string, reduceCh <-chan Datum, md Metadata) Messages {
	return rf(ctx, keys, reduceCh, md)
}
