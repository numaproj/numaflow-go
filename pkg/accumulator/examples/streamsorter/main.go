package main

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/numaproj/numaflow-go/pkg/accumulator"
)

// streamSorter buffers the incoming unordered (unordered by event-time but honors watermark) stream(s) and stores
// in the sorted sortedBuffer.
type streamSorter struct {
	latestWm time.Time
	// sortedBuffer stores the sorted data it has seen so far. The data stored will be completely sorted for
	// (current watermark - 1) because it honors watermark.
	sortedBuffer []accumulator.Datum
}

// Accumulate accumulates the read data in the sorted buffer and will emit the sorted results whenever the watermark
// progresses.
func (s *streamSorter) Accumulate(ctx context.Context, input <-chan accumulator.Datum, output chan<- accumulator.Message) {
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
			log.Println("Received datum with event time: ", datum.EventTime().UnixMilli())
			// watermark has moved, let's flush
			if datum.Watermark().After(s.latestWm) {
				s.latestWm = datum.Watermark()
				s.flushBuffer(output)
			}
			// store the data into the internal buffer
			s.insertSorted(datum)
		}
	}
}

// insertSorted will do a binary-search and inserts the Datum into the internal sorted buffer.
func (s *streamSorter) insertSorted(datum accumulator.Datum) {
	index := sort.Search(len(s.sortedBuffer), func(i int) bool {
		return s.sortedBuffer[i].EventTime().After(datum.EventTime())
	})
	s.sortedBuffer = append(s.sortedBuffer, datum)
	copy(s.sortedBuffer[index+1:], s.sortedBuffer[index:])
	s.sortedBuffer[index] = datum
}

// flushBuffer flushes the sorted results from the buffer into the outbound stream. While flushing it will
// truncate the buffer up till <= (current watermark - 1).
func (s *streamSorter) flushBuffer(output chan<- accumulator.Message) {
	log.Println("Watermark updated, flushing sortedBuffer: ", s.latestWm.UnixMilli())
	i := 0
	for _, datum := range s.sortedBuffer {
		if datum.EventTime().After(s.latestWm) {
			break
		}
		output <- accumulator.MessageFromDatum(datum)
		log.Println("Sent datum with event time: ", datum.EventTime().UnixMilli())
		i++
	}

	// truncate
	s.sortedBuffer = s.sortedBuffer[i:]
}

type streamSorterCreator struct {
}

// Create creates an Accumulator for every key. It will be closed only when the timeout has expired.
func (s *streamSorterCreator) Create() accumulator.Accumulator {
	return &streamSorter{latestWm: time.UnixMilli(-1)}
}

func main() {
	log.Println("Starting server")
	err := accumulator.NewServer(&streamSorterCreator{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start server: ", err)
	}
}
