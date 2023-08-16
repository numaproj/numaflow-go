package main

import (
	"context"
	"log"
	"strings"

	"github.com/numaproj/numaflow-go/pkg/map"
)

func mapFn(_ context.Context, keys []string, d _map.Datum) _map.Messages {
	msg := d.Value()
	_ = d.EventTime() // Event time is available
	_ = d.Watermark() // Watermark is available
	// Split the msg into an array with comma.
	strs := strings.Split(string(msg), ",")
	results := _map.MessagesBuilder()
	for _, s := range strs {
		results = results.Append(_map.NewMessage([]byte(s)))
	}
	return results
}

func main() {
	err := _map.NewServer(_map.MapperFunc(mapFn)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}
