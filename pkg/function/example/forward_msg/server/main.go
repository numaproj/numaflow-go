package main

import (
	"context"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func handle(ctx context.Context, key string, msg []byte) (functionsdk.Messages, error) {
	return functionsdk.MessagesBuilder().Append(functionsdk.MessageToAll(msg)), nil
}

func main() {
	server.NewServer().RegisterMapper(functionsdk.DoFunc(handle)).Start()
}
