package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"go.uber.org/atomic"

	"github.com/numaproj/numaflow-go/pkg/sessionreducer"
)

// Counter is a simple session reducer which counts the number of events in a session.
type Counter struct {
	count *atomic.Int32
}

func (c *Counter) SessionReduce(ctx context.Context, keys []string, input <-chan sessionreducer.Datum, outputCh chan<- sessionreducer.Message) {
	for range input {
		c.count.Inc()
	}
	outputCh <- sessionreducer.NewMessage([]byte(fmt.Sprintf("%d", c.count.Load()))).WithKeys(keys)
}

func (c *Counter) Accumulator(ctx context.Context) []byte {
	return []byte(strconv.Itoa(int(c.count.Load())))
}

func (c *Counter) MergeAccumulator(ctx context.Context, accumulator []byte) {
	val, err := strconv.Atoi(string(accumulator))
	if err != nil {
		log.Println("unable to convert the accumulator value to int: ", err.Error())
		return
	}
	c.count.Add(int32(val))
}

func NewSessionCounter() sessionreducer.SessionReducer {
	return &Counter{
		count: atomic.NewInt32(0),
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
