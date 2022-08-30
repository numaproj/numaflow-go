package main

import (
	"context"
	"strconv"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func reduceHandle(ctx context.Context, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) (functionsdk.Messages, error) {
	// sum up values for the same key
	intervalWindow := md.IntervalWindow()
	_ = intervalWindow
	var resultKey string
	var resultVal []byte
	var sum = 0
	var firstDatum = true
	// sum up the values
	for d := range reduceCh {
		if firstDatum {
			// key is the same for this sum use case
			key := d.Key()
			resultKey = key
			firstDatum = !firstDatum
		}
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
