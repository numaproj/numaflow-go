package main

import (
	"context"
	"fmt"
	"sync"

	functionsdk "github.com/KeranYang/numaflow-go/pkg/function"
	"github.com/KeranYang/numaflow-go/pkg/function/server"
)

var counts sync.Map

const SUCCESS_ITERATION = 3

// This imitates a function which fails the first 2 times it sees a message, but then
// succeeds on the third time.
// When it fails, it applies a "retry" Tag to the outgoing message, which causes the message to cycle back to it
// (this is designed to be used as a Map Vertex which Conditionally Forwards back to
// itself when the "retry" Tag is set, and forwards onto the next Vertex when it's not set).
func handle(_ context.Context, keys []string, d functionsdk.Datum) functionsdk.Messages {
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
		return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(msgBytes))
	} else {
		return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(msgBytes).WithTags([]string{"retry"}))
	}

}

func main() {
	server.New().RegisterMapper(functionsdk.MapFunc(handle)).Start(context.Background())
}
