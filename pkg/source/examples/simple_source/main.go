package main

import (
	"context"

	sourcesdk "github.com/numaproj/numaflow-go/pkg/source"
	"github.com/numaproj/numaflow-go/pkg/source/examples/simple_source/impl"
	"github.com/numaproj/numaflow-go/pkg/source/server"
)

func main() {
	simpleSource := impl.NewSimpleSource()
	server.New(
		sourcesdk.PendingFunc(simpleSource.Pending),
		sourcesdk.ReadFunc(simpleSource.Read),
		sourcesdk.AckFunc(simpleSource.Ack)).
		Start(context.Background())
}
