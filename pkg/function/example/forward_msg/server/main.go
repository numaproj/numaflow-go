package main

import (
	"context"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func map_handle(_ context.Context, key string, msg []byte) (functionsdk.Messages, error) {
	return functionsdk.MessagesBuilder().Append(functionsdk.MessageTo(key, msg)), nil
}

func main() {
	server.New().RegisterMapper(functionsdk.DoFunc(map_handle)).Start(context.Background())
}
