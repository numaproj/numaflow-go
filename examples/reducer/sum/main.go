package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/numaproj/numaflow-go/pkg/reducer"
)

// SumReducerCreator implements the reducer.ReducerCreator interface which creates a reducer
type SumReducerCreator struct {
}

func (s *SumReducerCreator) Create() reducer.Reducer {
	return &Sum{}
}

// Sum is a reducer that sum up the values for the given keys
type Sum struct {
	sum int
}

func (s *Sum) Reduce(ctx context.Context, keys []string, inputCh <-chan reducer.Datum, md reducer.Metadata) reducer.Messages {
	// sum up values for the same keys
	intervalWindow := md.IntervalWindow()
	_ = intervalWindow
	var resultKeys = keys
	var resultVal []byte
	// sum up the values
	for d := range inputCh {
		val := d.Value()

		// event time and watermark can be fetched from the datum
		eventTime := d.EventTime()
		_ = eventTime
		watermark := d.Watermark()
		_ = watermark

		v, err := strconv.Atoi(string(val))
		if err != nil {
			fmt.Printf("unable to convert the value to int: %v\n", err)
			continue
		}
		s.sum += v
	}
	resultVal = []byte(strconv.Itoa(s.sum))
	return reducer.MessagesBuilder().Append(reducer.NewMessage(resultVal).WithKeys(resultKeys))
}

func main() {
	err := reducer.NewServer(&SumReducerCreator{}).Start(context.Background())
	if err != nil {
		log.Panic("unable to start the server due to: ", err)
	}
}
