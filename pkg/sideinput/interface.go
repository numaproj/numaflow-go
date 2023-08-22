package sideinput

import (
	"context"
)

type RetrieverSideInput interface {
	// RetrieveSideInput is the function to process each side-input request.
	RetrieveSideInput(ctx context.Context) MessageSI
}

// RetrieverFunc is a utility type used to convert a RetrieveSideInput function to a RetrieverSideInput.
type RetrieverFunc func(ctx context.Context) MessageSI

// RetrieveSideInput implements the function of RetrieveSideInput function.
func (mf RetrieverFunc) RetrieveSideInput(ctx context.Context) MessageSI {
	return mf(ctx)
}
