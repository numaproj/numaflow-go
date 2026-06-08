package main

import (
	"context"
	"log"

	"github.com/numaproj/numaflow-go/pkg/accumulator"
)

// blackhole is an accumulator that intentionally discards every datum it receives without
// forwarding any data downstream.
//
// A naive implementation would simply read the input stream and emit nothing. However, an
// accumulator that never emits anything for the datums it consumes leaves the framework unable to
// release the per-datum tracked state, leading to unbounded memory growth.
//
// Instead, this example emits a drop message for every datum using accumulator.MessageToDrop. A
// drop message is not forwarded to the next vertex, but it still allows the framework to advance
// the watermark and release the tracked state (WAL) for that datum - giving us "blackhole"
// semantics without leaking memory. This pattern is useful for multiplexer-, cross-join-, or
// filter-style accumulators that legitimately need to omit some (or all) of their inputs.
type blackhole struct {
}

// Accumulate drains the input stream and drops every datum, emitting a drop message so the
// framework can release the datum's tracked state.
func (b *blackhole) Accumulate(ctx context.Context, input <-chan accumulator.Datum, output chan<- accumulator.Message) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Exiting the Accumulator")
			return
		case datum, ok := <-input:
			// this case happens due to timeout
			if !ok {
				log.Println("Input channel closed")
				return
			}
			log.Println("Dropping datum with event time: ", datum.EventTime().UnixMilli(), " watermark: ", datum.Watermark().UnixMilli())
			// Emit a drop message: nothing is forwarded downstream, but the framework still
			// advances the watermark and releases the tracked state for this datum.
			output <- accumulator.MessageToDrop(datum)
		}
	}
}

type blackholeCreator struct {
}

// Create creates an Accumulator for every key. It will be closed only when the timeout has expired.
func (b *blackholeCreator) Create() accumulator.Accumulator {
	return &blackhole{}
}

func main() {
	log.Println("Starting server")
	err := accumulator.NewServer(&blackholeCreator{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start server: ", err)
	}
}
