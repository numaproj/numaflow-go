package main

import (
	"context"
	"log"
	"strings"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

func mapFn(_ context.Context, keys []string, d mapper.Datum) mapper.Messages {
	msg := d.Value()
	_ = d.EventTime() // Operation time is available
	_ = d.Watermark() // Watermark is available
	// Split the msg into an array with comma.
	strs := strings.Split(string(msg), ",")
	results := mapper.MessagesBuilder()
	for _, s := range strs {
		results = results.Append(mapper.NewMessage([]byte(s)))
	}
	return results
}

func main() {
	err := mapper.NewServer(mapper.MapperFunc(mapFn)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}
