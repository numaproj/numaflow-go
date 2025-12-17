package main

import (
	"context"
	"fmt"
	"log"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

// logMetadataSink is a sinker implementation that logs the user metadata of the input
type logMetadataSink struct {
}

func (l *logMetadataSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		umd := d.UserMetadata()
		fmt.Println("User Metadata: ", umd)
		id := d.ID()
		result = result.Append(sinksdk.ResponseOK(id))
	}
	return result
}

func main() {
	err := sinksdk.NewServer(&logMetadataSink{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start sink function server: ", err)
	}
}
