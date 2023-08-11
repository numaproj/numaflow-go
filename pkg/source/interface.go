package source

import (
	"context"

	"github.com/numaproj/numaflow-go/pkg/source/model"
)

// ReadRequest is the interface of read request.
type ReadRequest interface {
	// Count returns the number of records to read.
	Count() uint64
}

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

// ReadHandler is the interface of read function implementation.
type ReadHandler interface {
	// HandleDo is the function to read the data from the source.
	HandleDo(ctx context.Context, readRequest ReadRequest, messageCh chan<- model.Message)
}

// ReadFunc is a utility type used to convert a HandleDo function to a ReadHandler.
type ReadFunc func(ctx context.Context, readRequest ReadRequest, messageCh chan<- model.Message)

// HandleDo implements the function of read function.
func (msf ReadFunc) HandleDo(ctx context.Context, readRequest ReadRequest, messageCh chan<- model.Message) {
	msf(ctx, readRequest, messageCh)
}
