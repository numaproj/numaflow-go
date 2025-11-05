package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

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
		if writeToPrimarySink() {
			fmt.Println("Primary sink write succeeded, writing to onSuccess sink - ", string(d.Value()))
			// write to onSuccess sink
			result = result.Append(sinksdk.ResponseOnSuccess(id, sinksdk.NewMessage([]byte("primary sink write succeeded"))))
		} else {
			fmt.Println("Primary sink write failed, writing to fallback sink - ", string(d.Value()))
			// write to fallback sink
			result = result.Append(sinksdk.ResponseFallback(id))
		}
	}
	return result
}

func writeToPrimarySink() bool {
	// simulate primary sink write failure/success
	return rand.Intn(2) == 1
}

func main() {
	err := sinksdk.NewServer(&onSuccessLogSink{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start sink function server: ", err)
	}
}
