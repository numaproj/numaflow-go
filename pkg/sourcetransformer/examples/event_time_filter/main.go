package main

import (
	"context"

	"github.com/numaproj/numaflow-go/pkg/sourcetransformer"
	"github.com/numaproj/numaflow-go/pkg/sourcetransformer/examples/event_time_filter/impl"
)

func mapTHandle(_ context.Context, keys []string, d sourcetransformer.Datum) sourcetransformer.Messages {
	return impl.FilterEventTime(keys, d)
}

func main() {
	sourcetransformer.NewServer(sourcetransformer.MapTFunc(mapTHandle)).Start(context.Background())
}
