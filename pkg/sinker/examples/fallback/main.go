package main

import (
	"context"
	"fmt"
	"log"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

// fbLogSink is a sinker implementation that logs the input to stdout
type fbLogSink struct {
}

func (l *fbLogSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		_ = d.EventTime()
		_ = d.Watermark()
		id := d.ID()
		fmt.Println("Primary sink under maintenance, writing to fallback sink - ", string(d.Value()))
		// write to fallback sink
		result = result.Append(sinksdk.ResponseFallback(id))
	}
	return result
}

func main() {
	err := sinksdk.NewServer(&fbLogSink{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start sink function server: ", err)
	}
}
