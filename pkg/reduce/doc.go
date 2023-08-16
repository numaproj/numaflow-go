// Package reduce implements the server code for reduce operation.
// Example Reduce
/*
	package main

	import (
		"context"
		"fmt"
		"log"
		"strconv"

		"github.com/numaproj/numaflow-go/pkg/reduce"
	)

	type Sum struct {
	}

	func (s *Sum) Reduce(ctx context.Context, keys []string, reduceCh <-chan reduce.Datum, md reduce.Metadata) reduce.Messages {
		// sum up values for the same keys
		intervalWindow := md.IntervalWindow()
		_ = intervalWindow
		var resultKeys = keys
		var resultVal []byte
		var sum = 0
		// sum up the values
		for d := range reduceCh {
			val := d.Value()
			eventTime := d.EventTime()
			_ = eventTime
			watermark := d.Watermark()
			_ = watermark

			v, err := strconv.Atoi(string(val))
			if err != nil {
				fmt.Printf("unable to convert the value to int: %v\n", err)
				continue
			}
			sum += v
		}
		resultVal = []byte(strconv.Itoa(sum))
		return reduce.MessagesBuilder().Append(reduce.NewMessage(resultVal).WithKeys(resultKeys))
	}

	func main() {
		err := reduce.NewServer(&Sum{}).Start(context.Background())
		if err != nil {
			log.Panic("unable to start the server due to: ", err)
		}
	}
*/

package reduce
