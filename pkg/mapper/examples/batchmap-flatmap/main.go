package main

import (
	"context"
	"log"
	"strings"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

func mapFn(_ context.Context, datums []mapper.Datum) mapper.BatchResponses {
	batchResponses := mapper.BatchResponsesBuilder()
	for _, d := range datums {
		msg := d.Value()
		_ = d.EventTime() // Event time is available
		_ = d.Watermark() // Watermark is available
		results := mapper.NewBatchResponse(d.Id())
		strs := strings.Split(string(msg), ",")
		for _, s := range strs {
			results = results.Append(mapper.NewMessage([]byte(s)))
		}
		batchResponses = batchResponses.Append(results)
	}
	return batchResponses
}

func main() {
	err := mapper.NewBatchServer(mapper.BatchMapperFunc(mapFn)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}
