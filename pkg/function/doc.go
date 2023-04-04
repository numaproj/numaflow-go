// Package function implements the server code for User Defined Function in golang.
//
// Example Map
/*
  package main

  import (

    "context"

    functionsdk "github.com/numaproj/numaflow-go/pkg/function"
    "github.com/numaproj/numaflow-go/pkg/function/server"

  )

  // Simply return the same msg

  func handle(ctx context.Context, keys string, data functionsdk.Datum) functionsdk.Messages {
    _ = data.EventTime() // Event time is available
    _ = data.Watermark() // Watermark is available
    return functionsdk.MessagesBuilder().Append(functionsdk.MessageToAll(data.Value()))
  }

  func main() {
    server.New().RegisterMapper(functionsdk.MapFunc(handle)).Start(context.Background())
  }
*/
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

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
  )

  func mapTHandle(_ context.Context, keys string, d functionsdk.Datum) functionsdk.MessageTs {
	eventTime := getEventTime(d.Value())
	return types.MessageTsBuilder().Append(types.MessageTTo(eventTime, keys, d.Value()))
  }

  func getEventTime(val []byte) time.Time {
	...
  }

  func main() {
	server.New().RegisterMapperT(functionsdk.MapTFunc(mapTHandle)).Start(context.Background())
  }
*/
//
// Example Reduce
/*
  package main

  import (
    "context"
    "strconv"

    functionsdk "github.com/numaproj/numaflow-go/pkg/function"
    "github.com/numaproj/numaflow-go/pkg/function/server"
  )

  // Count the incoming events

  func handle(_ context.Context, keys string, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) functionsdk.Messages {
    var resultKey = keys
    var resultVal []byte
    var counter = 0
    for _ = range reduceCh {
        counter++
    }
    resultVal = []byte(strconv.Itoa(counter))
    return functionsdk.MessagesBuilder().Append(functionsdk.MessageTo(resultKey, resultVal))
  }

  func main() {
    server.New().RegisterReducer(functionsdk.ReduceFunc(handle)).Start(context.Background())
  }
*/
package function
