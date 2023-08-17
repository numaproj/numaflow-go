package main

import (
	"context"
	"log"
	"time"

	"github.com/numaproj/numaflow-go/pkg/sourcetransformer"
)

// AssignEventTime is a source transformer that assigns event time to the message.
type AssignEventTime struct {
}

func (a *AssignEventTime) Transform(ctx context.Context, keys []string, d sourcetransformer.Datum) sourcetransformer.Messages {
	// Update message event time to time.Now()
	eventTime := time.Now()
	return sourcetransformer.MessagesBuilder().Append(sourcetransformer.NewMessage(d.Value(), eventTime).WithKeys(keys))
}

func main() {
	err := sourcetransformer.NewServer(&AssignEventTime{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}
