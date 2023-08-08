package main

import (
	"context"
	"time"

	"github.com/KeranYang/numaflow-go/pkg/function/server"

	functionsdk "github.com/KeranYang/numaflow-go/pkg/function"
)

func mapTHandle(_ context.Context, keys []string, d functionsdk.Datum) functionsdk.MessageTs {
	// Update message event time to time.Now()
	eventTime := time.Now()
	return functionsdk.MessageTsBuilder().Append(functionsdk.NewMessageT(d.Value(), eventTime).WithKeys(keys))
}

func main() {
	server.New().RegisterMapperT(functionsdk.MapTFunc(mapTHandle)).Start(context.Background())
}
