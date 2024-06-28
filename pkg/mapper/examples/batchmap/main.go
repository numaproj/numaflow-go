package main

import (
	"context"
	"log"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

func mapFn(_ context.Context, datums []mapper.Datum) mapper.BatchResponses {
	batchResponses := mapper.BatchResponsesBuilder()
	log.Println("MYDEBUG: length of input ", len(datums))
	for _, d := range datums {
		msg := d.Value()
		_ = d.EventTime() // Event time is available
		_ = d.Watermark() // Watermark is available
		results := mapper.NewBatchResponse(d.Id())
		for i := 0; i < 2; i++ {
			results = results.Append(mapper.NewMessage(msg))
		}
		batchResponses = batchResponses.Append(results)
	}
	return batchResponses
}

func main() {
	err := mapper.NewBatchMapServer(mapper.BatchMapperFunc(mapFn)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}
