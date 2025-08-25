package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/numaproj/numaflow-go/pkg/reducestreamer"
)

// Sum is a reducestreamer that sum up the values for the given keys and output the sum when the sum is greater than 100.
type Sum struct {
}

func (s *Sum) ReduceStream(ctx context.Context, keys []string, inputCh <-chan reducestreamer.Datum, outputCh chan<- reducestreamer.Message, md reducestreamer.Metadata) {
	// sum up values for the same keys
	intervalWindow := md.IntervalWindow()
	_ = intervalWindow
	var resultKeys = keys
	var resultVal []byte
	var sum = 0
	// sum up the values
	for d := range inputCh {
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

		if sum >= 100 {
			resultVal = []byte(strconv.Itoa(sum))
			outputCh <- reducestreamer.NewMessage(resultVal).WithKeys(resultKeys)
			sum = 0
		}
	}
	resultVal = []byte(strconv.Itoa(sum))
	outputCh <- reducestreamer.NewMessage(resultVal).WithKeys(resultKeys)
}

// SumCreator is the creator for the sum reducestreamer.
type SumCreator struct{}

func (s *SumCreator) Create() reducestreamer.ReduceStreamer {
	return &Sum{}
}

func main() {
	err := reducestreamer.NewServer(&SumCreator{}).Start(context.Background())
	if err != nil {
		log.Panic("unable to start the server due to: ", err)
	}
}
