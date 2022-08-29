package main

import (
	"context"
	"log"
	"os"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func handle(ctx context.Context, key string, msg []byte) (functionsdk.Messages, error) {
	// directly forward the input to the output
	return functionsdk.MessagesBuilder().Append(functionsdk.MessageTo(key, msg)), nil
}

func main() {
	file, _ := os.CreateTemp("/tmp", "numaflow-test.sock")
	defer func() {
		os.Remove(file.Name())
		log.Println("clean up sock")
	}()

	server.New().RegisterMapper(functionsdk.DoFunc(handle)).Start(server.WithSockAddr("/tmp/numaflow-test.sock"))

	log.Println("Exit")
}
