package main

import (
	"context"
	"fmt"
	"log"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

// logSink is a sinker implementation that logs the input to stdout
type logSink struct {
}

func (l *logSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		_ = d.EventTime()
		_ = d.Watermark()
		fmt.Println("User Defined Sink:", string(d.Value()))
		id := d.ID()
		result = result.Append(sinksdk.ResponseOK(id))
		// if we are not able to write to sink and if we have a fallback sink configured
		// we can use sinksdk.ResponseFallback(id)) to write the message to fallback sink
	}
	return result
}

func main() {
	err := sinksdk.NewServer(&logSink{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start sink function server: ", err)
	}
}
