package main

import (
	"context"
	"strconv"

	"github.com/numaproj/numaflow-go/pkg/reduce"
)

func reduceCounter(_ context.Context, keys []string, reduceCh <-chan reduce.Datum, md reduce.Metadata) reduce.Messages {
	// count the incoming events
	var resultKeys = keys
	var resultVal []byte
	var counter = 0
	for range reduceCh {
		counter++
	}
	resultVal = []byte(strconv.Itoa(counter))
	return reduce.MessagesBuilder().Append(reduce.NewMessage(resultVal).WithKeys(resultKeys))
}

func main() {
	reduce.NewServer(reduce.ReducerFunc(reduceCounter)).Start(context.Background())
}
