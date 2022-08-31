package main

import (
	"context"
	"strconv"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func reduceHandle(ctx context.Context, key string, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) (functionsdk.Messages, error) {
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
			continue
		}
		sum += v
	}
	resultVal = []byte(strconv.Itoa(sum))
	return functionsdk.MessagesBuilder().Append(functionsdk.MessageTo(resultKey, resultVal)), nil
}

func main() {
	server.New().RegisterReducer(functionsdk.ReduceFunc(reduceHandle)).Start(context.Background())
}
