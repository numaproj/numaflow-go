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

  func mapHandle(_ context.Context, keys []string, d functionsdk.Datum) functionsdk.Messages {
    // directly forward the input to the output
	val := d.Value()

	var resultKeys = keys
	var resultVal = val
	return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(resultVal).WithKeys(resultKeys))
  }

	func main() {
		server.NewServer(functionsdk.MapFunc(mapHandle)).Start(context.Background())
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

  func handle(_ context.Context, keys []string, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) functionsdk.Messages {
    var resultKeys = keys
    var resultVal []byte
    var counter = 0
    for _ = range reduceCh {
        counter++
    }
    resultVal = []byte(strconv.Itoa(counter))
    return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(resultVal).WithKeys(resultKey))
  }

  func main() {
    server.NewReduceServer(functionsdk.ReduceFunc(handle)).Start(context.Background())
  }
*/
// The Datum object contains the message payload and metadata. Currently, there are two fields in metadata:
// the message ID, the message delivery count to indicate how many times the message has been delivered.
// You can use these metadata to implement customized logic. For example,
/*
...
func handle(ctx context.Context, keys []string, data functionsdk.Datum) functionsdk.Messages {
    deliveryCount := data.Metadata().NumDelivered()
    // Choose to do specific actions, if the message delivery count reaches a certain threshold.
    if deliveryCount > 3:
        ...
}

*/
package _map
