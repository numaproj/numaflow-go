package main

import (
	"context"
	"strings"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func handle(_ context.Context, key string, d functionsdk.Datum) functionsdk.Messages {
	msg := d.Value()
	_ = d.EventTime() // Event time is available
	_ = d.Watermark() // Watermark is available
	// Split the msg into an array with comma.
	strs := strings.Split(string(msg), ",")
	results := functionsdk.MessagesBuilder()
	for _, s := range strs {
		results = results.Append(functionsdk.MessageToAll([]byte(s)))
	}
	return results
}

func main() {
	server.New().RegisterMapper(functionsdk.MapFunc(handle)).Start(context.Background())
}
