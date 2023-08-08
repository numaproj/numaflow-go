package main

import (
	"context"
	"fmt"
	"strconv"

	functionsdk "github.com/KeranYang/numaflow-go/pkg/function"
	"github.com/KeranYang/numaflow-go/pkg/function/server"
)

func reduceHandle(_ context.Context, keys []string, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) functionsdk.Messages {
	// sum up values for the same keys
	intervalWindow := md.IntervalWindow()
	_ = intervalWindow
	var resultKeys = keys
	var resultVal []byte
	var sum = 0
	// sum up the values
	for d := range reduceCh {
		val := d.Value()
		eventTime := d.EventTime()
		_ = eventTime
		watermark := d.Watermark()
		_ = watermark

		v, err := strconv.Atoi(string(val))
		if err != nil {
			fmt.Printf("unable to convert the value to int: %v\n", err)
			continue
		}
		sum += v
	}
	resultVal = []byte(strconv.Itoa(sum))
	return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(resultVal).WithKeys(resultKeys))
}

func main() {
	server.New().RegisterReducer(functionsdk.ReduceFunc(reduceHandle)).Start(context.Background())
}
