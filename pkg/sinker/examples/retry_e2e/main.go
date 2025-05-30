package main

import (
	"context"
	"fmt"
	"log"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

// retrySink is a sinker implementation that creates a failed response
type retrySink struct {
}

func (l *retrySink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		id := d.ID()
		result = result.Append(sinksdk.ResponseFailure(id, fmt.Sprintf("test retry strategy for id: %s", id)))
	}
	return result
}

func main() {
	err := sinksdk.NewServer(&retrySink{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start sink function server: ", err)
	}
}
