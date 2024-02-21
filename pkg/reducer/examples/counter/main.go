package main

import (
	"context"
	"log"
	"strconv"

	"github.com/numaproj/numaflow-go/pkg/reducer"
)

func reduceCounter(_ context.Context, keys []string, inputCh <-chan reducer.Datum, md reducer.Metadata) reducer.Messages {
	// count the incoming events
	var resultKeys = keys
	var resultVal []byte
	var counter = 0
	for range inputCh {
		counter++
	}
	resultVal = []byte(strconv.Itoa(counter))
	return reducer.MessagesBuilder().Append(reducer.NewMessage(resultVal).WithKeys(resultKeys))
}

func main() {
	err := reducer.NewServer(reducer.SimpleCreatorWithReduceFn(reduceCounter)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start server: ", err)
	}
}
