package main

import (
	"context"
	"log"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

// noResponseSink is a sinker implementation that does not return any response
type noResponseSink struct {
}

func (l *noResponseSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for range datumStreamCh {
		// Process the datum but don't append any response
		// This simulates a sink that doesn't respond
	}
	return result
}

func main() {
	err := sinksdk.NewServer(&noResponseSink{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start no response sink server: ", err)
	}
}
