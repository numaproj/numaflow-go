package main

import (
	"context"
	"log"
	"strconv"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

type EvenOdd struct {
}

func (e *EvenOdd) Map(ctx context.Context, keys []string, d mapper.Datum) mapper.Messages {
	// users can access system metadata using d.metadata().systemMetadata() -> read only methods on this type
	// users can access user metadata using d.metadata().UserMetadata() -> read and write methods on this type

	// UserMetadata data := UserMetadata::new()
	// mapper.NewMessage().WithMetadata(mapper.NewMetadata(userMetadata UserMetadata))
	msg := d.Value()
	_ = d.EventTime() // Event time is available
	_ = d.Watermark() // Watermark is available
	// If msg is not an integer, drop it, otherwise return it with "even" or "odd" key.
	if num, err := strconv.Atoi(string(msg)); err != nil {
		return mapper.MessagesBuilder().Append(mapper.MessageToDrop())
	} else if num%2 == 0 {
		return mapper.MessagesBuilder().Append(mapper.NewMessage(msg).WithKeys([]string{"even"}).WithTags([]string{"even-tag"}))
	} else {
		return mapper.MessagesBuilder().Append(mapper.NewMessage(msg).WithKeys([]string{"odd"}).WithTags([]string{"odd-tag"}))
	}
}

func main() {
	err := mapper.NewServer(&EvenOdd{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}
