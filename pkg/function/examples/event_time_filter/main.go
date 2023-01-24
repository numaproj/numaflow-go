package main

import (
	"context"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/examples/event_time_filter/impl"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func mapTHandle(_ context.Context, key string, d functionsdk.Datum) functionsdk.MessageTs {
	return impl.FilterEventTime(key, d)
}

func main() {
	server.New().RegisterMapperT(functionsdk.MapTFunc(mapTHandle)).Start(context.Background())
}
