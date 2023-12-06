package sourcer

import (
	"context"
	"time"
)

// Sourcer is the interface for implementation of the source.
type Sourcer interface {
	// Read reads the data from the source and sends the data to the message channel.
	// If the read request is timed out, the function returns without reading new data.
	// Right after reading a message, the function marks the offset as to be acked.
	// Read should never attempt to close the message channel as the caller owns the channel.
	Read(ctx context.Context, readRequest ReadRequest, messageCh chan<- Message)
	// Ack acknowledges the data from the source.
	Ack(ctx context.Context, request AckRequest)
	// Pending returns the number of pending messages.
	// When the return value is negative, it indicates the pending information is not available.
	// With pending information being not available, the Numaflow platform doesn't auto-scale the source.
	Pending(ctx context.Context) int64
	//Partitions returns the partitions associated with the source.
	Partitions(ctx context.Context) []int32
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
