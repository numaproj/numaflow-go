package main

import (
	"context"
	"strings"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func handle(_ context.Context, _ []string, d functionsdk.Datum, messageCh chan<- functionsdk.Message) {
	defer close(messageCh)
	msg := d.Value()
	_ = d.EventTime() // Event time is available
	_ = d.Watermark() // Watermark is available
	// Split the msg into an array with comma.
	strs := strings.Split(string(msg), ",")
	for _, s := range strs {
		messageCh <- functionsdk.NewMessage([]byte(s))
	}
}

func main() {
	server.NewMapStreamServer(context.Background(), functionsdk.MapStreamFunc(handle)).Start(context.Background())
}
