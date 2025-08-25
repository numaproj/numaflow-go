package main

import (
	"context"
	"log"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

// serveSink is a sinker implementation that logs the input to stdout
type serveSink struct {
}

func (l *serveSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		id := d.ID()
		result = result.Append(sinksdk.ResponseServe(id, d.Value()))
	}
	return result
}

func main() {
	err := sinksdk.NewServer(&serveSink{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start sink function server: ", err)
	}
}
