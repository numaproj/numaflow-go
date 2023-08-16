package main

import (
	"context"
	"log"

	"github.com/numaproj/numaflow-go/pkg/source"
	"github.com/numaproj/numaflow-go/pkg/source/examples/simple_source/impl"
)

func main() {
	simpleSource := impl.NewSimpleSource()
	err := source.NewServer(simpleSource).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start source server : ", err)
	}
}
