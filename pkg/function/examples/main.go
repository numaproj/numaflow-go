package main

import (
	"context"

	"github.com/numaproj/numaflow-go/pkg/function"
)

func handle(ctx context.Context, key string, msg []byte) (function.Messages, error) {
	return function.MessagesBuilder().Append(function.MessageToAll(msg)), nil
}

func main() {
	function.Initialize().RegisterMapper(function.DoFunc(handle)).Start()
}
