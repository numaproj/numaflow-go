package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	sideinputsdk "github.com/numaproj/numaflow-go/pkg/sideinput"
)

// handle is the side input handler function.
func handle(_ context.Context) sideinputsdk.Message {
	t := time.Now()
	// val is the side input message value. This would be the value that the side input vertex receives.
	val := "an example: " + string(t.String())
	// randomly drop side input message. Note that the side input message is not retried.
	// NoBroadcastMessage() is used to drop the message and not to
	// broadcast it to other side input vertices.
	if rand.Int()%2 == 0 {
		return sideinputsdk.NoBroadcastMessage()
	}
	// BroadcastMessage() is used to broadcast the message with the given value to other side input vertices.
	// val must be converted to []byte.
	return sideinputsdk.BroadcastMessage([]byte(val))
}

func main() {
	// Start the side input server.
	err := sideinputsdk.NewSideInputServer(sideinputsdk.RetrieveFunc(handle)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start side input server: ", err)
	}
}
