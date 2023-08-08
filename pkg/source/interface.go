package source

import (
	"context"
)

// PendingHandler is the interface of pending function implementation.
type PendingHandler interface {
	// HandleDo is the function to return the number of pending records at the user defined source.
	HandleDo(ctx context.Context) uint64
}

// PendingFunc is a utility type used to convert a HandleDo function to a PendingHandler.
type PendingFunc func(ctx context.Context) uint64

// HandleDo implements the function of pending function.
func (pf PendingFunc) HandleDo(ctx context.Context) uint64 {
	return pf(ctx)
}
