package main

import (
	"context"
	"github.com/numaproj/numaflow-go/pkg/function/types"
	"strconv"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func handle(_ context.Context, key string, d functionsdk.Datum) types.Messages {
	msg := d.Value()
	_ = d.EventTime() // Event time is available
	_ = d.Watermark() // Watermark is available
	// If msg is not an integer, drop it, otherwise return it with "even" or "odd" key.
	if num, err := strconv.Atoi(string(msg)); err != nil {
		return types.MessagesBuilder().Append(types.MessageToDrop())
	} else if num%2 == 0 {
		return types.MessagesBuilder().Append(types.MessageTo("even", msg))
	} else {
		return types.MessagesBuilder().Append(types.MessageTo("odd", msg))
	}
}

func main() {
	server.New().RegisterMapper(functionsdk.MapFunc(handle)).Start(context.Background())
}
