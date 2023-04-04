package main

import (
	"context"
	"strconv"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func reduceHandle(_ context.Context, keys []string, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) functionsdk.Messages {
	// count the incoming events
	var resultKeys = keys
	var resultVal []byte
	var counter = 0
	for _ = range reduceCh {
		counter++
	}
	resultVal = []byte(strconv.Itoa(counter))
	return functionsdk.MessagesBuilder().Append(functionsdk.MessageTo(resultKeys, resultVal))
}

func main() {
	server.New().RegisterReducer(functionsdk.ReduceFunc(reduceHandle)).Start(context.Background())
}
