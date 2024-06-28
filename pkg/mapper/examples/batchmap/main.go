package main

import (
	"context"
	"log"
	"strings"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

func mapFn(_ context.Context, datums []mapper.Datum) mapper.BatchResponses {
	batchResponses := mapper.BatchResponsesBuilder()
	log.Println("MYDEBUG: length of input ", len(datums))
	for _, d := range datums {
		msg := d.Value()
		_ = d.EventTime() // Event time is available
		_ = d.Watermark() // Watermark is available
		// Split the msg into an array with comma.
		strs := strings.Split(string(msg), ",")
		results := mapper.NewBatchResponse(d.Id())
		for _, s := range strs {
			results = results.Append(mapper.NewMessage([]byte(s)))
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
