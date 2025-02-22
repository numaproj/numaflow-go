package servingstore

import (
	"context"
)

// ServingStorer is the interface for serving store to store and retrieve from a custom store.
type ServingStorer interface {
	// Put is to put data into the Serving Store.
	Put(ctx context.Context, put PutRequester)

	// Get is to retrieve data from the Serving Store.
	Get(ctx context.Context) StoredResults
}

// PutRequester interface exposes methods to retrieve data from the Put rpc.
type PutRequester interface {
	Origin() string
	Payload() [][]byte
}

type PutRequest struct {
	origin   string
	payloads [][]byte
}

// Origin returns the origin name.
func (p *PutRequest) Origin() string {
	return p.origin
}

// Payload returns the payloads to be stored.
func (p *PutRequest) Payload() [][]byte {
	return p.payloads
}
