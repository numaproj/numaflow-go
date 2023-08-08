package main

import (
	"context"
	"fmt"

	"github.com/KeranYang/numaflow-go/pkg/sink/server"

	sinksdk "github.com/KeranYang/numaflow-go/pkg/sink"
)

func handle(_ context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
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
	server.New().RegisterSinker(sinksdk.SinkFunc(handle)).Start(context.Background())
}
