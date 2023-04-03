package main

import (
	"context"
	"fmt"
	"strconv"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func handle(_ context.Context, key string, d functionsdk.Datum) functionsdk.Messages {
	msg := d.Value()
	_ = d.EventTime() // Event time is available
	_ = d.Watermark() // Watermark is available
	// If msg is not an integer, drop it, otherwise return it with "even" or "odd" key.
	if num, err := strconv.Atoi(string(key)); err != nil {
		fmt.Println("dropping message (not an integer)")
		return functionsdk.MessagesBuilder().Append(functionsdk.MessageToDrop())
	} else if num%2 == 0 {
		return functionsdk.MessagesBuilder().Append(functionsdk.MessageTo("even", msg))
	} else {
		return functionsdk.MessagesBuilder().Append(functionsdk.MessageTo("odd", msg))
	}
}

func main() {
	server.New().RegisterMapper(functionsdk.MapFunc(handle)).Start(context.Background())
}
