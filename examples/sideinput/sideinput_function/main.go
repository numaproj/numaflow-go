package main

import (
	"context"
	sideinputsdk "github.com/numaproj/numaflow-go/pkg/sideinput"
	"log"
)

var counter = 0

// handle is the side input handler function.
func handle(_ context.Context) sideinputsdk.Message {
	// BroadcastMessage() is used to broadcast the message with the given value to other side input vertices.
	// val must be converted to []byte.
	counter = (counter + 1) % 10
	if counter%2 == 0 {
		return sideinputsdk.BroadcastMessage([]byte(`e2e-even`))
	}

	return sideinputsdk.BroadcastMessage([]byte(`e2e-odd`))
}

func main() {
	// Start the side input server.
	err := sideinputsdk.NewSideInputServer(sideinputsdk.RetrieveFunc(handle)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start side input server: ", err)
	}
}
