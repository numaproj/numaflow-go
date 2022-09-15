package main

import (
	"context"
	"fmt"

	"github.com/numaproj/numaflow-go/pkg/sink"
	sinksdk "github.com/numaproj/numaflow-go/sink"
)

func handle(ctx context.Context, msgs []sink.Message) (sink.Responses, error) {
	result := sink.ResponsesBuilder()
	for _, m := range msgs {
		fmt.Println(string(m.Payload))
		result = result.Append(sink.ResponseOK(m.ID))
	}
	return result, nil
}

func main() {
	sinksdk.Start(context.Background(), handle)
}
