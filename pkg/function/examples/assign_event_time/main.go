package main

import (
	"context"
	"time"

	"github.com/numaproj/numaflow-go/pkg/source"
	"github.com/numaproj/numaflow-go/pkg/source/server"
)

func mapTHandle(_ context.Context, keys []string, d source.Datum) source.MessageTs {
	// Update message event time to time.Now()
	eventTime := time.Now()
	return source.MessageTsBuilder().Append(source.NewMessageT(d.Value(), eventTime).WithKeys(keys))
}

func main() {
	server.NewSourceTransformerServer(context.Background(), source.MapTFunc(mapTHandle)).Start(context.Background())
}
