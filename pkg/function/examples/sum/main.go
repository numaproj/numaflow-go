package main

import (
	"context"
	"fmt"
	"github.com/numaproj/numaflow-go/pkg/function/types"
	"strconv"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func reduceHandle(_ context.Context, key string, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) types.Messages {
	// sum up values for the same key
	intervalWindow := md.IntervalWindow()
	_ = intervalWindow
	var resultKey = key
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
	return types.MessagesBuilder().Append(types.MessageTo(resultKey, resultVal))
}

func main() {
	server.New().RegisterReducer(functionsdk.ReduceFunc(reduceHandle)).Start(context.Background())
}
