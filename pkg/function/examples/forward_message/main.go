package main

import (
	"context"

	functionsdk "github.com/KeranYang/numaflow-go/pkg/function"
	"github.com/KeranYang/numaflow-go/pkg/function/server"
)

func mapHandle(_ context.Context, keys []string, d functionsdk.Datum) functionsdk.Messages {
	// directly forward the input to the output
	val := d.Value()
	eventTime := d.EventTime()
	_ = eventTime
	watermark := d.Watermark()
	_ = watermark

	var resultKeys = keys
	var resultVal = val
	return functionsdk.MessagesBuilder().Append(functionsdk.NewMessage(resultVal).WithKeys(resultKeys))
}

func main() {
	server.New().RegisterMapper(functionsdk.MapFunc(mapHandle)).Start(context.Background())
}
