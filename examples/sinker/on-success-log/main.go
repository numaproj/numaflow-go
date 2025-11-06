package main

import (
	"context"
	"fmt"
	"log"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

// onSuccessLogSink is a sinker implementation that logs the input to stdout
type onSuccessLogSink struct {
}

func (l *onSuccessLogSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		_ = d.EventTime()
		_ = d.Watermark()
		id := d.ID()
		fmt.Println("Primary sink write succeeded, writing to onSuccess sink - ", string(d.Value()))
		// write to onSuccess sink
		result = result.Append(sinksdk.ResponseOnSuccess(id, sinksdk.NewMessage([]byte("on-success-message"))))
	}
	return result
}

func main() {
	err := sinksdk.NewServer(&onSuccessLogSink{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start sink function server: ", err)
	}
}
