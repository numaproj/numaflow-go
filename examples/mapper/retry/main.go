package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	_map "github.com/numaproj/numaflow-go/pkg/mapper"
)

var counts sync.Map

const SUCCESS_ITERATION = 3

// This imitates a function which fails the first 2 times it sees a message, but then
// succeeds on the third time.
// When it fails, it applies a "retry" Tag to the outgoing message, which causes the message to cycle back to it
// (this is designed to be used as a Map Vertex which Conditionally Forwards back to
// itself when the "retry" Tag is set, and forwards onto the next Vertex when it's not set).
func mapFn(_ context.Context, keys []string, d _map.Datum) _map.Messages {
	msgBytes := d.Value()
	msg := string(msgBytes)

	currCounts, found := counts.Load(msg)
	var newCounts int
	if !found {
		newCounts = 1
	} else {
		newCounts = currCounts.(int) + 1
	}
	counts.Store(msg, newCounts)
	fmt.Printf("count for %q=%d\n", msg, newCounts)
	if newCounts >= SUCCESS_ITERATION {
		// imitate successful outgoing message here
		counts.Delete(msg)
		return _map.MessagesBuilder().Append(_map.NewMessage(msgBytes))
	} else {
		return _map.MessagesBuilder().Append(_map.NewMessage(msgBytes).WithTags([]string{"retry"}))
	}

}

func main() {
	err := _map.NewServer(_map.MapperFunc(mapFn)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}
