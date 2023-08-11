package main

import (
	"context"

	sourcesdk "github.com/numaproj/numaflow-go/pkg/source"
	"github.com/numaproj/numaflow-go/pkg/source/server"
)

func handle(_ context.Context) uint64 {
	// The simple source always returns 0 to indicate no pending records.
	return 0
}

func main() {
	server.New(sourcesdk.PendingFunc(handle)).Start(context.Background())
}
