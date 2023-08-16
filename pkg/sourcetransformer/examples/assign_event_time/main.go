package main

import (
	"context"
	"time"

	"github.com/numaproj/numaflow-go/pkg/sourcetransformer"
)

func mapTHandle(_ context.Context, keys []string, d sourcetransformer.Datum) sourcetransformer.Messages {
	// Update message event time to time.Now()
	eventTime := time.Now()
	return sourcetransformer.MessagesBuilder().Append(sourcetransformer.NewMessage(d.Value(), eventTime).WithKeys(keys))
}

func main() {
	sourcetransformer.NewServer(sourcetransformer.MapTFunc(mapTHandle)).Start(context.Background())
}
