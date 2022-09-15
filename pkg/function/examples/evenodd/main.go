package main

import (
	"context"
	"strconv"

	"github.com/numaproj/numaflow-go/pkg/datum"
	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func handle(_ context.Context, key string, d datum.Datum) functionsdk.Messages {
	msg := d.Value()
	eventTime := d.EventTime() // Event time is available
	_ = eventTime
	watermark := d.Watermark() // Watermark is available
	_ = watermark
	// If msg is not an integer, drop it, otherwise return it with "even" or "odd" key.
	if num, err := strconv.Atoi(string(msg)); err != nil {
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
