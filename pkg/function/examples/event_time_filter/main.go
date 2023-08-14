package main

import (
	"context"

	"github.com/numaproj/numaflow-go/pkg/function/examples/event_time_filter/impl"
	"github.com/numaproj/numaflow-go/pkg/source"
	"github.com/numaproj/numaflow-go/pkg/source/server"
)

func mapTHandle(_ context.Context, keys []string, d source.Datum) source.MessageTs {
	return impl.FilterEventTime(keys, d)
}

func main() {
	server.NewSourceTransformerServer(source.MapTFunc(mapTHandle)).Start(context.Background())
}
