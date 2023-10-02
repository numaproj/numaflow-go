package main

import (
	"context"
	"log"

	"github.com/numaproj/numaflow-go/pkg/sourcer"

	"github.com/numaproj/numaflow-go/pkg/sourcer/examples/simple_source/impl"
)

func main() {
	simpleSource := impl.NewSimpleSource()
	err := sourcer.NewServer(simpleSource).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start source server : ", err)
	}
}
