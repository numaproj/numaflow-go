# Numaflow Golang SDK

This SDK provides the interfaces to implement [Numaflow](https://github.com/numaproj/numaflow) User Defined Functions or Sinks in Golang.

## Implement User Defined Functions

```golang
package main

import (
	"context"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

// Simply return the same msg
func handle(ctx context.Context, key string, data functionsdk.Datum) functionsdk.Messages {
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

	sinksdk "github.com/numaproj/numaflow-go/pkg/sink"
	"github.com/numaproj/numaflow-go/pkg/sink/server"
)

func handle(ctx context.Context, datumList []sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for _, datum := range datumList {
		fmt.Println(string(datum.Value()))
		result = result.Append(sinksdk.ResponseOK(datum.ID()))
	}
	return result
}

func main() {
	server.New().RegisterSinker(sinksdk.SinkFunc(handle)).Start(context.Background())
}

```
