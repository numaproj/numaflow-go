package main

import (
	"context"
	"fmt"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sink"
	"github.com/numaproj/numaflow-go/pkg/sink/server"
)

func handle(ctx context.Context, sinkCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for d := range sinkCh {
		_ = d.EventTime()
		_ = d.Watermark()
		fmt.Println("User Defined Sink:", string(d.Value()))
		id := d.ID()
		result = result.Append(sinksdk.ResponseOK(id))
	}
	return result
}

func main() {
	server.New().RegisterSinker(sinksdk.SinkFunc(handle)).Start(context.Background())
}
