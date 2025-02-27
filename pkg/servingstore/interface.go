package servingstore

import (
	"context"
)

// ServingStorer is the interface for serving store to store and retrieve from a custom store.
type ServingStorer interface {
	// Put is to put data into the Serving Store.
	Put(ctx context.Context, put PutDatum)

	// Get is to retrieve data from the Serving Store.
	Get(ctx context.Context, get GetDatum) StoredResult
}

// PutDatum interface exposes methods to retrieve data from the Put rpc.
type PutDatum interface {
	ID() string
	Payloads() []Payload
}

// PutRequest contains the details to store the payload to the Store.
type PutRequest struct {
	id       string
	payloads []Payload
}

// ID returns the id of the original request.
func (p *PutRequest) ID() string {
	return p.id
}

// Payloads returns the payloads to be stored.
func (p *PutRequest) Payloads() []Payload {
	return p.payloads
}

// GetDatum is the interface to expose methods to retrieve from the Get rpc.
type GetDatum interface {
	ID() string
}

// GetRequest has details on the Get rpc.
type GetRequest struct {
	id string
}

// ID is the unique ID original request which is used get the data stored in the Store.
func (g *GetRequest) ID() string {
	return g.id
}
