package sideinput

import (
	"context"
)

// RetrieveSideInputHandler is the interface of the RetrieveSideInput function implementation.
type RetrieveSideInputHandler interface {
	// HandleDo is the function to process each side-input request.
	HandleDo(ctx context.Context) MessageSI
}

// RetrieveSideInput is utility type used to convert a HandleDo function to a RetrieveSideInputHandler.
type RetrieveSideInput func(ctx context.Context) MessageSI

// HandleDo implements the function of RetrieveSideInput function.
func (mf RetrieveSideInput) HandleDo(ctx context.Context) MessageSI {
	return mf(ctx)
}
