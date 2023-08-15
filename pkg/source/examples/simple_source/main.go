package main

import (
	"context"

	"github.com/numaproj/numaflow-go/pkg/source/examples/simple_source/impl"
	"github.com/numaproj/numaflow-go/pkg/source/server"
)

func main() {
	simpleSource := impl.NewSimpleSource()
	server.New(simpleSource).Start(context.Background())
}
