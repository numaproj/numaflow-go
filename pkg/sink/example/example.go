package main

import (
	"context"
	"fmt"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sink"
	"github.com/numaproj/numaflow-go/pkg/sink/server"
)

func handle(ctx context.Context, datumList []sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for _, d := range datumList {
		fmt.Println(string(d.Value()))
		id := d.ID()
		result = result.Append(sinksdk.ResponseOK(id))
	}
	return result
}

func main() {
	server.New().RegisterSinker(sinksdk.SinkFunc(handle)).Start(context.Background())
}
