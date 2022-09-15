package main

import (
	"context"
	"fmt"

	"github.com/numaproj/numaflow-go/pkg/datum"
	"github.com/numaproj/numaflow-go/pkg/sink"
	sinksdk "github.com/numaproj/numaflow-go/pkg/sink"
	"github.com/numaproj/numaflow-go/pkg/sink/server"
)

func handle(ctx context.Context, datumList []datum.Datum) sink.Responses {
	result := sink.ResponsesBuilder()
	for _, d := range datumList {
		fmt.Println(string(d.Value()))
		id := "TODO"
		result = result.Append(sink.ResponseOK(id))
	}
	return result
}

func main() {
	server.New().RegisterSinker(sinksdk.SinkFunc(handle)).Start(context.Background())
}
