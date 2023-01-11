package main

import (
	"context"
	"github.com/numaproj/numaflow-go/pkg/function/types"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func mapHandle(_ context.Context, key string, d functionsdk.Datum) types.Messages {
	// directly forward the input to the output
	val := d.Value()
	eventTime := d.EventTime()
	_ = eventTime
	watermark := d.Watermark()
	_ = watermark

	var resultKey = key
	var resultVal = val
	return types.MessagesBuilder().Append(types.MessageTo(resultKey, resultVal))
}

func main() {
	server.New().RegisterMapper(functionsdk.MapFunc(mapHandle)).Start(context.Background())
}
