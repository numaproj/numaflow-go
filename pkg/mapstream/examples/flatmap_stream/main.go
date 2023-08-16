package main

import (
	"context"
	"strings"

	"github.com/numaproj/numaflow-go/pkg/mapstream"
)

func handle(_ context.Context, _ []string, d mapstream.Datum, messageCh chan<- mapstream.Message) {
	defer close(messageCh)
	msg := d.Value()
	_ = d.EventTime() // Event time is available
	_ = d.Watermark() // Watermark is available
	// Split the msg into an array with comma.
	strs := strings.Split(string(msg), ",")
	for _, s := range strs {
		messageCh <- mapstream.NewMessage([]byte(s))
	}
}

func main() {
	mapstream.NewServer(context.Background(), mapstream.MapStreamFunc(handle)).Start(context.Background())
}
