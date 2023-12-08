package main

import (
	"context"
	"log"
	"strings"

	"github.com/numaproj/numaflow-go/pkg/mapstreamer"
)

// FlatMap is a MapStreamer that split the input message into multiple messages and stream them.
type FlatMap struct {
}

func (f *FlatMap) MapStream(ctx context.Context, keys []string, d mapstreamer.Datum, messageCh chan<- mapstreamer.Message) {
	defer close(messageCh)
	msg := d.Value()
	_ = d.EventTime() // Event time is available
	_ = d.Watermark() // Watermark is available
	// Split the msg into an array with comma.
	strs := strings.Split(string(msg), ",")
	for _, s := range strs {
		messageCh <- mapstreamer.NewMessage([]byte(s))
	}
}

func main() {
	err := mapstreamer.NewServer(&FlatMap{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map stream function server: ", err)
	}
}
