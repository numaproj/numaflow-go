package main

import (
	"context"
	"strings"

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
