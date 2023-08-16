package main

import (
	"context"
	"fmt"
	"log"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sink"
)

func logSink(_ context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		_ = d.EventTime()
		_ = d.Watermark()
		fmt.Println("User Defined Sink:", string(d.Value()))
		id := d.ID()
		result = result.Append(sinksdk.ResponseOK(id))
	}
	return result
}

func main() {
	err := sinksdk.NewSinkServer(sinksdk.SinkerFunc(logSink)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start sink function server: ", err)
	}
}
