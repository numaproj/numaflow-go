// Package _map implements the server code for Map Operation.
//
// Example Map
/*
	package main

	import (
		"context"
		"log"

		"github.com/numaproj/numaflow-go/pkg/mapper"
	)

	// Forward is a mapper that directly forward the input to the output
	type Forward struct {
	}

	func (f *Forward) Map(ctx context.Context, keys []string, d mapper.Datum) mapper.Messages {
		// directly forward the input to the output
		val := d.Value()
		eventTime := d.EventTime()
		_ = eventTime
		watermark := d.Watermark()
		_ = watermark

		var resultKeys = keys
		var resultVal = val
		return mapper.MessagesBuilder().Append(mapper.NewMessage(resultVal).WithKeys(resultKeys))
	}

	func main() {
		err := mapper.NewServer(&Forward{}).Start(context.Background())
		if err != nil {
			log.Panic("Failed to start map function server: ", err)
		}
	}
*/

package mapper
