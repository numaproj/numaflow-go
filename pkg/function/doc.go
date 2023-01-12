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
//
// Example MapT (assigning a random event time to a message)
/*
  package main

  import (
	"context"
	"math/rand"
	"time"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
  )

  func mapTHandle(_ context.Context, key string, d functionsdk.Datum) functionsdk.MessageTs {
	// assign a random event time to the message, then forward the input to the output.
	randomEventTime := generateRandomTime
	return types.MessageTsBuilder().Append(types.MessageTTo(randomEventTime(), key, d.Value()))
  }

  // generateRandomTime generates a random timestamp within date range [1970-01-01 to 2023-01-01]
  func generateRandomTime() time.Time {
	min := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
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
