package main

import (
	"context"
	"strconv"

	"github.com/numaproj/numaflow-go/pkg/reduceStreamer"
)

func reduceCounter(_ context.Context, keys []string, inputCh <-chan reduceStreamer.Datum, outputCh chan<- reduceStreamer.Message, md reduceStreamer.Metadata) {
	// count the incoming events
	var resultKeys = keys
	var resultVal []byte
	var counter = 0
	for range inputCh {
		counter++
		if counter >= 10 {
			resultVal = []byte(strconv.Itoa(counter))
			outputCh <- reduceStreamer.NewMessage(resultVal).WithKeys(resultKeys)
			counter = 0
		}
	}
	resultVal = []byte(strconv.Itoa(counter))
	outputCh <- reduceStreamer.NewMessage(resultVal).WithKeys(resultKeys)
}

func main() {
	reduceStreamer.NewServer(reduceStreamer.ReducerFunc(reduceCounter)).Start(context.Background())
}
