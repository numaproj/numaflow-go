package main

import (
	"context"
	"time"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func mapTHandle(_ context.Context, keys []string, d functionsdk.Datum) functionsdk.MessageTs {
	// Update message event time to time.Now()
	eventTime := time.Now()
	return functionsdk.MessageTsBuilder().Append(functionsdk.NewMessageT(d.Value()).WithKeys(keys).WithEventTime(eventTime))
}

func main() {
	server.New().RegisterMapperT(functionsdk.MapTFunc(mapTHandle)).Start(context.Background())
}
