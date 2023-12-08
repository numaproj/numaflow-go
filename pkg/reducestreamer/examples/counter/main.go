package main

import (
	"context"
	"strconv"

	"github.com/numaproj/numaflow-go/pkg/reducestreamer"
)

func reduceCounter(_ context.Context, keys []string, inputCh <-chan reducestreamer.Datum, outputCh chan<- reducestreamer.Message, md reducestreamer.Metadata) {
	// count the incoming events
	var resultKeys = keys
	var resultVal []byte
	var counter = 0
	for range inputCh {
		counter++
		if counter >= 10 {
			resultVal = []byte(strconv.Itoa(counter))
			outputCh <- reducestreamer.NewMessage(resultVal).WithKeys(resultKeys)
			counter = 0
		}
	}
	resultVal = []byte(strconv.Itoa(counter))
	outputCh <- reducestreamer.NewMessage(resultVal).WithKeys(resultKeys)
}

func main() {
	reducestreamer.NewServer(reducestreamer.ReducerFunc(reduceCounter)).Start(context.Background())
}
