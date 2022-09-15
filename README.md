# Numaflow Golang SDK

This SDK provides the interfaces to implement [Numaflow](https://github.com/numaproj/numaflow) User Defined Functions or Sinks in Golang.

## Implement User Defined Functions

```golang
package main

import (
	"context"

	"github.com/numaproj/numaflow-go/pkg/datum"
	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

// Simply return the same msg
func handle(ctx context.Context, key string, data datum.Datum) functionsdk.Messages {
	_ = data.EventTime() // Event time is available
	_ = data.Watermark() // Watermark is available
	return functionsdk.MessagesBuilder().Append(functionsdk.MessageToAll(data.Value()))
}

func main() {
	server.New().RegisterMapper(functionsdk.MapFunc(handle)).Start(context.Background())
}

```

## Implement User Defined Sinks

```golang
package main

import (
	"context"
	"fmt"

	sinksdk "github.com/numaproj/numaflow-go/sink"
)

func handle(ctx context.Context, msgs []sinksdk.Message) (sinksdk.Responses, error) {
	result := sinksdk.ResponsesBuilder()
	for _, m := range msgs {
		fmt.Println(string(m.Payload))
		result = result.Append(sinksdk.ResponseOK(m.ID))
	}
	return result, nil
}

func main() {
	sinksdk.Start(context.Background(), handle)
}

```
