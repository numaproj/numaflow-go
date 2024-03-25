package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"go.uber.org/atomic"

	"github.com/numaproj/numaflow-go/pkg/sessionreducer"
)

// Sum is a simple session reducer which computes sum of events in a session.
type Sum struct {
	sum *atomic.Int32
}

func (c *Sum) SessionReduce(ctx context.Context, keys []string, input <-chan sessionreducer.Datum, outputCh chan<- sessionreducer.Message) {
	for d := range input {
		val, err := strconv.Atoi(string(d.Value()))
		if err != nil {
			log.Panic("unable to convert the value to int: ", err.Error())
		} else {
			c.sum.Add(int32(val))
		}
	}
	outputCh <- sessionreducer.NewMessage([]byte(fmt.Sprintf("%d", c.sum.Load()))).WithKeys(keys)
}

func (c *Sum) Accumulator(ctx context.Context) []byte {
	return []byte(strconv.Itoa(int(c.sum.Load())))
}

func (c *Sum) MergeAccumulator(ctx context.Context, accumulator []byte) {
	val, err := strconv.Atoi(string(accumulator))
	if err != nil {
		log.Println("unable to convert the accumulator value to int: ", err.Error())
		return
	}
	c.sum.Add(int32(val))
}

func NewSessionCounter() sessionreducer.SessionReducer {
	return &Sum{
		sum: atomic.NewInt32(0),
	}
}

// SessionCounterCreator is the creator for the session reducer.
type SessionCounterCreator struct{}

func (s *SessionCounterCreator) Create() sessionreducer.SessionReducer {
	return NewSessionCounter()
}

func main() {
	sessionreducer.NewServer(&SessionCounterCreator{}).Start(context.Background())
}
