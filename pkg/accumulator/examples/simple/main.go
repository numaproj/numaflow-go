package main

import (
	"context"
	"log"
	"strconv"

	"github.com/numaproj/numaflow-go/pkg/accumulator"
)

func simpleAccumulator(ctx context.Context, inputCh <-chan accumulator.Datum, outputCh chan<- accumulator.Datum) {
	counter := 0
	for {
		select {
		case <-ctx.Done():
			println("Context done")
			return
		case datum := <-inputCh:
			println("Received datum")
			counter += 1
			datum.UpdateValue([]byte(strconv.Itoa(counter)))
			outputCh <- datum
			println("Sent datum")
		}
	}
}

func main() {
	println("Starting server")
	err := accumulator.NewServer(accumulator.SimpleCreatorWithAccumulateFn(simpleAccumulator)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start server: ", err)
	}
}
