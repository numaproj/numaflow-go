package main

import (
	"context"
	"log"

	"github.com/numaproj/numaflow-go/pkg/servingstore"
)

type InMemoryStore struct {
	store map[string][]servingstore.Payload
}

func (i *InMemoryStore) Put(ctx context.Context, putDatum servingstore.PutDatum) {
	id := putDatum.ID()
	log.Printf("Received Put request for %s", id)
	if _, ok := i.store[id]; !ok {
		i.store[id] = make([]servingstore.Payload, 0)
	}
	for _, payload := range putDatum.Payloads() {
		i.store[id] = append(i.store[id], servingstore.NewPayload(payload.Origin(), payload.Value()))
	}
}

func (i *InMemoryStore) Get(ctx context.Context, getDatum servingstore.GetDatum) servingstore.StoredResult {
	id := getDatum.ID()
	log.Printf("Received Get request for %s", id)
	if data, ok := i.store[id]; ok {
		return servingstore.NewStoredResult(id, data)
	} else {
		return servingstore.NewStoredResult(id, nil)
	}
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		store: make(map[string][]servingstore.Payload),
	}
}

func main() {
	err := servingstore.NewServer(NewInMemoryStore()).Start(context.Background())
	if err != nil {
		log.Panic("Failed to serving store function server: ", err)
	}
}
