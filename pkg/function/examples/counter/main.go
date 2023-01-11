package main

import (
	"context"
	"github.com/numaproj/numaflow-go/pkg/function/types"
	"strconv"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
)

func reduceHandle(_ context.Context, key string, reduceCh <-chan functionsdk.Datum, md functionsdk.Metadata) types.Messages {
	// count the incoming events
	var resultKey = key
	var resultVal []byte
	var counter = 0
	for _ = range reduceCh {
		counter++
	}
	resultVal = []byte(strconv.Itoa(counter))
	return types.MessagesBuilder().Append(types.MessageTo(resultKey, resultVal))
}

func main() {
	server.New().RegisterReducer(functionsdk.ReduceFunc(reduceHandle)).Start(context.Background())
}
