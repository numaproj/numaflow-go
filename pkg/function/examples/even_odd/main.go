package main

import (
	"context"
	"strconv"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func handle(_ context.Context, keys []string, d functionsdk.Datum) functionsdk.Messages {
	msg := d.Value()
	_ = d.EventTime() // Event time is available
	_ = d.Watermark() // Watermark is available
	// If msg is not an integer, drop it, otherwise return it with "even" or "odd" key.
	if num, err := strconv.Atoi(string(msg)); err != nil {
		return functionsdk.MessagesBuilder().Append(functionsdk.MessageToDrop())
	} else if num%2 == 0 {
		return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(msg).WithKeys([]string{"even"}).WithTags([]string{"even-tag"}))
	} else {
		return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(msg).WithKeys([]string{"odd"}).WithTags([]string{"odd-tag"}))
	}
}

func main() {
	server.New().RegisterMapper(functionsdk.MapFunc(handle)).Start(context.Background())
}
