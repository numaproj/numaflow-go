package main

import (
	"context"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func handle(ctx context.Context, key string, msg []byte) (functionsdk.Messages, error) {
	// directly forward the input to the output
	return functionsdk.MessagesBuilder().Append(functionsdk.MessageTo(key, msg)), nil
}

func main() {
	server.New().RegisterMapper(functionsdk.DoFunc(handle)).Start()
}
