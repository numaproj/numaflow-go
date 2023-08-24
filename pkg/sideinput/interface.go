package sideinput

import (
	"context"
)

// SideInputRetriever is the interface for side input retrieval implementation.
type SideInputRetriever interface {
	// RetrieveSideInput is the function to process each side-input request.
	RetrieveSideInput(ctx context.Context) Message
}

// RetrieveFunc is a utility type used to convert a RetrieveSideInput function to a SideInputRetriever.
type RetrieveFunc func(ctx context.Context) Message

// RetrieveSideInput implements the function of RetrieveSideInput function.
func (mf RetrieveFunc) RetrieveSideInput(ctx context.Context) Message {
	return mf(ctx)
}
