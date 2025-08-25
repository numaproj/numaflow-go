package main

import (
	"context"
	"log"

	"simple_source/impl"

	"github.com/numaproj/numaflow-go/pkg/sourcer"
)

func main() {
	simpleSource := impl.NewSimpleSource()
	err := sourcer.NewServer(simpleSource).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start source server : ", err)
	}
}
