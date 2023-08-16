// Package source implements the server code for User Defined Source and Transformer in golang.
//
// Example MapT (extracting event time from the datum payload)
// MapT includes both Map and EventTime assignment functionalities.
// Although the input datum already contains EventTime and Watermark, it's up to the MapT implementor to
// decide on whether to use them for generating new EventTime.
// MapT can be used only at source vertex by source data transformer.
/*
  package main

  import (
	"context"
	"time"

    "github.com/numaproj/numaflow-go/pkg/source"
	"github.com/numaproj/numaflow-go/pkg/function/server"
  )

  func mapTHandle(_ context.Context, keys []string, d source.Datum) source.Messages {
	eventTime := getEventTime(d.Value())
	return source.MessagesBuilder().Append(source.NewMessage(eventTime, d.Value()).WithKeys(keys)))
  }

  func getEventTime(val []byte) time.Time {
	...
  }

  func main() {
	server.NewServer(source.MapTFunc(mapTHandle)).Start(context.Background())
  }
*/

package sourcetransformer
