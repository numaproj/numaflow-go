// Package mapstream implements the server code for Map Stream Operation.
//
// Example MapStream
/*
	package main

	import (
		"context"
		"log"
		"strings"

		"github.com/numaproj/numaflow-go/pkg/mapstream"
	)

	type FlatMap struct {
	}

	func (f *FlatMap) MapStream(ctx context.Context, keys []string, d mapstream.Datum, messageCh chan<- mapstream.Message) {
		defer close(messageCh)
		msg := d.Value()
		_ = d.EventTime() // Event time is available
		_ = d.Watermark() // Watermark is available
		// Split the msg into an array with comma.
		strs := strings.Split(string(msg), ",")
		for _, s := range strs {
			messageCh <- mapstream.NewMessage([]byte(s))
		}
	}

	func main() {
		err := mapstream.NewServer(&FlatMap{}).Start(context.Background())
		if err != nil {
			log.Panic("Failed to start map stream function server: ", err)
		}
	}
*/
package mapstream
