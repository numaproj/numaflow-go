package main

import (
	"context"
	"log"
	"strconv"

	"github.com/numaproj/numaflow-go/pkg/map"
)

type EvenOdd struct {
}

func (e *EvenOdd) Map(ctx context.Context, keys []string, d _map.Datum) _map.Messages {
	msg := d.Value()
	_ = d.EventTime() // Event time is available
	_ = d.Watermark() // Watermark is available
	// If msg is not an integer, drop it, otherwise return it with "even" or "odd" key.
	if num, err := strconv.Atoi(string(msg)); err != nil {
		return _map.MessagesBuilder().Append(_map.MessageToDrop())
	} else if num%2 == 0 {
		return _map.MessagesBuilder().Append(_map.NewMessage(msg).WithKeys([]string{"even"}).WithTags([]string{"even-tag"}))
	} else {
		return _map.MessagesBuilder().Append(_map.NewMessage(msg).WithKeys([]string{"odd"}).WithTags([]string{"odd-tag"}))
	}
}

func main() {
	err := _map.NewServer(&EvenOdd{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}
