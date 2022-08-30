package main

import (
	"context"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func mapHandle(_ context.Context, d functionsdk.Datum) (functionsdk.Messages, error) {
	// directly forward the input to the output
	key := d.Key()
	val := d.Value()
	eventTime := d.EventTime()
	_ = eventTime
	watermark := d.Watermark()
	_ = watermark

	var resultKey = key
	var resultVal = val
	return functionsdk.MessagesBuilder().Append(functionsdk.MessageTo(resultKey, resultVal)), nil
}

func main() {
	server.New().RegisterMapper(functionsdk.MapFunc(mapHandle)).Start(context.Background())
}
