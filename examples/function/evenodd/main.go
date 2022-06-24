package main

import (
	"context"
	"strconv"

	funcsdk "github.com/numaproj/numaflow-go/function"
)

// If msg is not an integer, drop it, otherwise return it with "even" or "odd" key.
func handle(ctx context.Context, key, msg []byte) (funcsdk.Messages, error) {
	if num, err := strconv.Atoi(string(msg)); err != nil {
		return funcsdk.MessagesBuilder().Append(funcsdk.MessageToDrop()), nil
	} else if num%2 == 0 {
		return funcsdk.MessagesBuilder().Append(funcsdk.MessageTo("even", msg)), nil
	} else {
		return funcsdk.MessagesBuilder().Append(funcsdk.MessageTo("odd", msg)), nil
	}
}

func main() {
	funcsdk.Start(context.Background(), handle)
}
