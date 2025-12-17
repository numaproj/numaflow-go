package main

import (
	"context"
	"log"
	"time"

	"github.com/numaproj/numaflow-go/pkg/sourcetransformer"
)

// AssignEventTime is a source transformer that assigns event time to the message.
type MetadataEventTime struct {
}

func (m *MetadataEventTime) Transform(ctx context.Context, keys []string, d sourcetransformer.Datum) sourcetransformer.Messages {
	eventTime := time.Now()
	umd := d.UserMetadata()
	if umd != nil {
		umd.AddKV("event-time-group", "event-time", []byte(eventTime.Format(time.RFC3339)))
	}
	return sourcetransformer.MessagesBuilder().Append(sourcetransformer.NewMessage(d.Value(), eventTime).WithKeys(keys).WithUserMetadata(umd))
}

func main() {
	err := sourcetransformer.NewServer(&MetadataEventTime{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start metadata event time server: ", err)
	}
}
