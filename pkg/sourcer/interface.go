package sourcer

import (
	"context"
	"time"
)

// Sourcer is the interface for implementation of the source.
type Sourcer interface {
	// Read reads the data from the source and sends the data to the message channel.
	// Read should never attempt to close the message channel as the caller owns the channel.
	Read(ctx context.Context, readRequest ReadRequest, messageCh chan<- Message)
	// Ack acknowledges the data from the source.
	Ack(ctx context.Context, request AckRequest)
	// Pending returns the number of pending messages.
	// When the return value is negative, it indicates the pending information is not available.
	Pending(ctx context.Context) int64
}

// ReadRequest is the interface of read request.
type ReadRequest interface {
	// Count returns the number of records to read.
	Count() uint64
	// TimeOut returns the timeout of the read request.
	TimeOut() time.Duration
}

// AckRequest is the interface of ack request.
type AckRequest interface {
	// Offsets returns the offsets of the records to ack.
	Offsets() []Offset
}
