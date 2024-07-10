package main

import (
	"context"
	"log"
	"strings"

	"github.com/numaproj/numaflow-go/pkg/batchmapper"
)

func batchMapFn(_ context.Context, datums <-chan batchmapper.Datum) batchmapper.BatchResponses {
	batchResponses := batchmapper.BatchResponsesBuilder()
	for d := range datums {
		msg := d.Value()
		_ = d.EventTime() // Event time is available
		_ = d.Watermark() // Watermark is available
		batchResponse := batchmapper.NewBatchResponse(d.Id())
		strs := strings.Split(string(msg), ",")
		for _, s := range strs {
			batchResponse = batchResponse.Append(batchmapper.NewMessage([]byte(s)))
		}

		batchResponses = batchResponses.Append(batchResponse)
	}
	return batchResponses
}

func main() {
	err := batchmapper.NewServer(batchmapper.BatchMapperFunc(batchMapFn)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start batch map function server: ", err)
	}
}
