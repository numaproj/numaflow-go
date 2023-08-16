// Package source implements the server code for User Defined Source Transformer in golang.
//
// Example SourceTransformer (extracting event time from the datum payload)
// MapT includes both Map and EventTime assignment functionalities.
// Although the input datum already contains EventTime and Watermark, it's up to the MapT implementor to
// decide on whether to use them for generating new EventTime.
// MapT can be used only at source vertex by source data transformer.
/*
  package main

  import (
	"context"
	"time"

	"github.com/numaproj/numaflow-go/pkg/sourcetransformer"
  )

  func mapTHandle(_ context.Context, keys []string, d sourcetransformer.Datum) sourcetransformer.Messages {
	eventTime := getEventTime(d.Value())
	return sourcetransformer.MessagesBuilder().Append(sourcetransformer.NewMessage(eventTime, d.Value()).WithKeys(keys)))
  }

  func getEventTime(val []byte) time.Time {
	...
  }

  func main() {
	sourcetransformer.NewServer(source.SourceTransformFunc(mapTHandle)).Start(context.Background())
  }
*/

package sourcetransformer
