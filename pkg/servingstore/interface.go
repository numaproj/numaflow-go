package servingstore

import (
	"context"
)

// ServingStorer is the interface for serving store to store and retrieve from a custom store.
type ServingStorer interface {
	// Put is to put data into the Serving Store.
	Put(ctx context.Context, put PutRequester)

	// Get is to retrieve data from the Serving Store.
	Get(ctx context.Context, get GetRequester) StoredResults
}

// PutRequester interface exposes methods to retrieve data from the Put rpc.
type PutRequester interface {
	Origin() string
	Payload() [][]byte
}

// PutRequest contains the details to store the payload to the Store.
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

// GetRequester is the interface to expose methods to retrieve from the Get rpc.
type GetRequester interface {
	Id() string
}

// GetRequest has details on the Get rpc.
type GetRequest struct {
	id string
}

// Id is the unique ID original request which is used get the data stored in the Store.
func (g *GetRequest) Id() string {
	return g.id
}
