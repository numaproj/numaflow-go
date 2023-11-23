package main

import (
	"context"
	"strconv"

	"github.com/numaproj/numaflow-go/pkg/reducer"
)

func reduceCounter(_ context.Context, keys []string, inputCh <-chan reducer.Datum, outputCh chan<- reducer.Message, md reducer.Metadata) {
	// count the incoming events
	var resultKeys = keys
	var resultVal []byte
	var counter = 0
	for range inputCh {
		counter++
	}
	resultVal = []byte(strconv.Itoa(counter))
	outputCh <- reducer.NewMessage(resultVal).WithKeys(resultKeys)
}

func main() {
	reducer.NewServer(reducer.ReducerFunc(reduceCounter)).Start(context.Background())
}
