package main

import (
	"context"
	"log"

	"github.com/numaproj/numaflow-go/pkg/sourcetransformer"
	"github.com/numaproj/numaflow-go/pkg/sourcetransformer/examples/event_time_filter/impl"
)

func transform(_ context.Context, keys []string, d sourcetransformer.Datum) sourcetransformer.Messages {
	return impl.FilterEventTime(keys, d)
}

func main() {
	err := sourcetransformer.NewServer(sourcetransformer.SourceTransformFunc(transform)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start source transform server: ", err)
	}
}
