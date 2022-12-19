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

  func handle(ctx context.Context, key string, data functionsdk.Datum) functionsdk.Messages {
    _ = data.EventTime() // Event time is available
    _ = data.Watermark() // Watermark is available
    return functionsdk.MessagesBuilder().Append(functionsdk.MessageToAll(data.Value()))
  }

  func main() {
    server.New().RegisterMapper(functionsdk.MapFunc(handle)).Start(context.Background())
  }
*/
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

  func handle(_ context.Context, key string, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) functionsdk.Messages {
    var resultKey = key
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
