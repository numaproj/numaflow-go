package main

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/numaproj/numaflow-go/pkg/accumulator"
)

type streamSorter struct {
	latestWm time.Time
	buffer   []accumulator.Datum
}

func (s *streamSorter) Accumulate(ctx context.Context, input <-chan accumulator.Datum, output chan<- accumulator.Message) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Context done")
			return
		case datum := <-input:
			log.Println("Received datum with event time: ", datum.EventTime().UnixMilli())
			if datum.Watermark().After(s.latestWm) {
				s.latestWm = datum.Watermark()
				s.flushBuffer(output)
			}
			s.insertSorted(datum)
		}
	}
}

func (s *streamSorter) insertSorted(datum accumulator.Datum) {
	index := sort.Search(len(s.buffer), func(i int) bool {
		return s.buffer[i].EventTime().After(datum.EventTime())
	})
	s.buffer = append(s.buffer, datum)
	copy(s.buffer[index+1:], s.buffer[index:])
	s.buffer[index] = datum
}

func (s *streamSorter) flushBuffer(output chan<- accumulator.Message) {
	log.Println("Watermark updated, flushing buffer: ", s.latestWm.UnixMilli())
	i := 0
	for _, datum := range s.buffer {
		if datum.EventTime().After(s.latestWm) {
			break
		}
		output <- accumulator.MessageFromDatum(datum)
		println("Sent datum with event time: ", datum.EventTime().UnixMilli())
		i++
	}
	s.buffer = s.buffer[i:]
}

type streamSorterCreator struct {
}

func (s *streamSorterCreator) Create() accumulator.Accumulator {
	return &streamSorter{latestWm: time.UnixMilli(-1)}
}

func main() {
	println("Starting server")
	err := accumulator.NewServer(&streamSorterCreator{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start server: ", err)
	}
}
