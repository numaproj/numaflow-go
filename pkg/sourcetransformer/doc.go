// Package source implements the server code for Source Transformer in golang.
//
// Example Transform (extracting event time from the datum payload)
// Transform includes both Map and EventTime assignment functionalities.
// Although the input datum already contains EventTime and Watermark, it's up to the Transform implementor to
// decide on whether to use them for generating new EventTime.
// Transform can be used only at source vertex by source data transformer.
/*
	package main

	import (
		"context"
		"log"
		"time"

		"github.com/numaproj/numaflow-go/pkg/sourcetransform"
	)

	// AssignEventTime is a source transformer that assigns event time to the message.
	type AssignEventTime struct {
	}

	func (a *AssignEventTime) Transform(ctx context.Context, keys []string, d sourcetransform.Datum) sourcetransform.Messages {
		// Update message event time to time.Now()
		eventTime := time.Now()
		return sourcetransform.MessagesBuilder().Append(sourcetransform.NewMessage(d.Value(), eventTime).WithKeys(keys))
	}

	func main() {
		err := sourcetransform.NewServer(&AssignEventTime{}).Start(context.Background())
		if err != nil {
			log.Panic("Failed to start map function server: ", err)
		}
	}
*/

package sourcetransformer
