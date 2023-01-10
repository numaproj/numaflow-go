package main

import (
	"context"
	sourcesdk "github.com/numaproj/numaflow-go/pkg/source"

	"github.com/numaproj/numaflow-go/pkg/source/server"
)

func transformHandle(_ context.Context, key string, d sourcesdk.Datum) sourcesdk.Messages {
	// directly forward the input to the output without changing the event time
	_ = d.EventTime() // Event time is available
	_ = d.Watermark() // Watermark is available
	return sourcesdk.MessagesBuilder().Append(sourcesdk.MessageTo(d.EventTime(), key, d.Value()))
}

func main() {
	server.New().RegisterTransformer(sourcesdk.TransformFunc(transformHandle)).Start(context.Background())
}
