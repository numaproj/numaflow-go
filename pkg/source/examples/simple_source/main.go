package main

import (
	"context"

	sourcesdk "github.com/numaproj/numaflow-go/pkg/source"
	"github.com/numaproj/numaflow-go/pkg/source/model"
	"github.com/numaproj/numaflow-go/pkg/source/server"
)

func pending(_ context.Context) uint64 {
	// The simple source always returns 0 to indicate no pending records.
	return 0
}

func read(ctx context.Context, readRequest sourcesdk.ReadRequest, messageCh chan<- model.Message) {
}

func ack(ctx context.Context, request sourcesdk.AckRequest) {
}

func main() {
	server.New(
		sourcesdk.PendingFunc(pending),
		sourcesdk.ReadFunc(read),
		sourcesdk.AckFunc(ack)).Start(context.Background())
}
