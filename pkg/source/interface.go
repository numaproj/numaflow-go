package source

import (
	"context"
	"time"

	"github.com/numaproj/numaflow-go/pkg/source/model"
)

// PendingHandler is the interface of pending function implementation.
type PendingHandler interface {
	// HandleDo is the function to return the number of pending records at the user-defined source.
	HandleDo(ctx context.Context) uint64
}

// PendingFunc is a utility type used to convert a HandleDo function to a PendingHandler.
type PendingFunc func(ctx context.Context) uint64

// HandleDo implements the function of pending function.
func (pf PendingFunc) HandleDo(ctx context.Context) uint64 {
	return pf(ctx)
}

// ReadRequest is the interface of read request.
type ReadRequest interface {
	// Count returns the number of records to read.
	Count() uint64
	// TimeOut returns the timeout of the read request.
	TimeOut() time.Duration
}

// ReadHandler is the interface of read function implementation.
type ReadHandler interface {
	// HandleDo is the function to read the data from the source.
	HandleDo(ctx context.Context, readRequest ReadRequest, messageCh chan<- model.Message)
}

// ReadFunc is a utility type used to convert a HandleDo function to a ReadHandler.
type ReadFunc func(ctx context.Context, readRequest ReadRequest, messageCh chan<- model.Message)

// HandleDo implements the function of read function.
func (rf ReadFunc) HandleDo(ctx context.Context, readRequest ReadRequest, messageCh chan<- model.Message) {
	rf(ctx, readRequest, messageCh)
}

// AckRequest is the interface of ack request.
type AckRequest interface {
	// Offsets returns the offsets of the records to ack.
	Offsets() []model.Offset
}

// AckHandler is the interface of ack function implementation.
type AckHandler interface {
	// HandleDo is the function to ack the data from the source.
	HandleDo(ctx context.Context, request AckRequest)
}

// AckFunc is a utility type used to convert a HandleDo function to an AckHandler.
type AckFunc func(ctx context.Context, request AckRequest)

// HandleDo implements the function of ack function.
func (af AckFunc) HandleDo(ctx context.Context, request AckRequest) {
	af(ctx, request)
}
