package main

import (
	"context"
	"strings"

	funcsdk "github.com/numaproj/numaflow-go/function"
)

// Split the msg into an array with comma.
func handle(ctx context.Context, key, msg []byte) (funcsdk.Messages, error) {
	strs := strings.Split(string(msg), ",")
	results := funcsdk.MessagesBuilder()
	for _, s := range strs {
		results = results.Append(funcsdk.MessageToAll([]byte(s)))
	}
	return results, nil
}

func main() {
	funcsdk.Start(context.Background(), handle)
}
