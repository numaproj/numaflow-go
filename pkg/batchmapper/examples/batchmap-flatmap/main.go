package main

import (
	"context"
	"log"
	"strings"

	"github.com/numaproj/numaflow-go/pkg/batchmapper"
)

func batchMapFn(_ context.Context, datums []batchmapper.Datum) batchmapper.BatchResponses {
	batchResponses := batchmapper.BatchResponsesBuilder()
	for _, d := range datums {
		msg := d.Value()
		_ = d.EventTime() // Event time is available
		_ = d.Watermark() // Watermark is available
		results := batchmapper.NewBatchResponse(d.Id())
		strs := strings.Split(string(msg), ",")
		for _, s := range strs {
			results = results.Append(batchmapper.NewMessage([]byte(s)))
		}
		batchResponses = batchResponses.Append(results)
	}
	return batchResponses
}

func main() {
	err := batchmapper.NewServer(batchmapper.BatchMapperFunc(batchMapFn)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}
