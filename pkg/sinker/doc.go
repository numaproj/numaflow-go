// Package sinker implements the server code for user defined sink.
//
// Example:
/*
	 package main

	import (
		"context"
		"fmt"
		"log"

		sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
	)

	// logSink is a sinker implementation that logs the input to stdout
	type logSink struct {
	}

	func (l *logSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
		result := sinksdk.ResponsesBuilder()
		for d := range datumStreamCh {
			_ = d.EventTime()
			_ = d.Watermark()
			fmt.Println("User Defined Sink:", string(d.Value()))
			id := d.ID()
			result = result.Append(sinksdk.ResponseOK(id))
		}
		return result
	}

	func main() {
		err := sinksdk.NewSinkServer(&logSink{}).Start(context.Background())
		if err != nil {
			log.Panic("Failed to start sink function server: ", err)
		}
	}

*/
package sinker
