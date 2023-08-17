package main

import (
	"context"
	"log"

	"github.com/numaproj/numaflow-go/pkg/map"
)

func mapFn(_ context.Context, keys []string, d _map.Datum) _map.Messages {
	// directly forward the input to the output
	val := d.Value()
	eventTime := d.EventTime()
	_ = eventTime
	watermark := d.Watermark()
	_ = watermark

	var resultKeys = keys
	var resultVal = val
	return _map.MessagesBuilder().Append(_map.NewMessage(resultVal).WithKeys(resultKeys))
}

func main() {
	err := _map.NewServer(_map.MapperFunc(mapFn)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}