package main

import (
	"context"
	"strconv"

	"github.com/numaproj/numaflow-go/pkg/reducer"
)

func reduceCounter(_ context.Context, keys []string, reduceCh <-chan reducer.Datum, md reducer.Metadata) reducer.Messages {
	// count the incoming events
	var resultKeys = keys
	var resultVal []byte
	var counter = 0
	for range reduceCh {
		counter++
	}
	resultVal = []byte(strconv.Itoa(counter))
	return reducer.MessagesBuilder().Append(reducer.NewMessage(resultVal).WithKeys(resultKeys))
}

func main() {
	reducer.NewServer(reducer.SimpleCreatorWithReduceFn(reduceCounter)).Start(context.Background())
}
