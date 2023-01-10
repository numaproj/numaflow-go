package main

import (
	"context"
	sourcesdk "github.com/numaproj/numaflow-go/pkg/source"

	"github.com/numaproj/numaflow-go/pkg/source/server"
)

func transformHandle(_ context.Context, key string, d sourcesdk.Datum) sourcesdk.Messages {
	// directly forward the input to the output
	val := d.Value()
	eventTime := d.EventTime()
	_ = eventTime
	watermark := d.Watermark()
	_ = watermark

	var resultKey = key
	var resultVal = val
	return sourcesdk.MessagesBuilder().Append(sourcesdk.MessageTo(d.EventTime(), resultKey, resultVal))
}

func main() {
	server.New().RegisterTransformer(sourcesdk.TransformFunc(transformHandle)).Start(context.Background())
}
