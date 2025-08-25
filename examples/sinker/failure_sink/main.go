package main

import (
	"context"
	"fmt"
	"log"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

// failureSink is a sinker implementation that creates a failed response
type failureSink struct {
}

func (l *failureSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		id := d.ID()
		result = result.Append(sinksdk.ResponseFailure(id, fmt.Sprintf("test retry strategy for id: %s", id)))
	}
	return result
}

func main() {
	err := sinksdk.NewServer(&failureSink{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start failure sink server: ", err)
	}
}
