package main

import (
	"context"
	"fmt"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

var counts map[string]int

const SUCCESS_ITERATION = 3

// This imitates a function which fails the first 2 times it sees a message, but then
// succeeds on the third time.
// When it fails, it applies a "retry" Tag to the outgoing message, which causes the message to cycle back to it
// (this is designed to be used as a Map Vertex which Conditionally Forwards back to
// itself when the "retry" Tag is set, and forwards onto the next Vertex when it's not set).
func handle(_ context.Context, keys []string, d functionsdk.Datum) functionsdk.Messages {
	msgBytes := d.Value()
	msg := string(msgBytes)

	_, found := counts[msg]
	if !found {
		counts[msg] = 1
	} else {
		counts[msg]++
	}
	fmt.Printf("count for %q=%d\n", msg, counts[msg])
	if counts[msg] >= SUCCESS_ITERATION {
		// imitate successful outgoing message here
		delete(counts, msg)
		return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(msgBytes))
	} else {
		return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(msgBytes).WithTags([]string{"retry"}))
	}

}

func main() {

	counts = make(map[string]int)

	server.New().RegisterMapper(functionsdk.MapFunc(handle)).Start(context.Background())
}
