package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"go.uber.org/atomic"

	"github.com/numaproj/numaflow-go/pkg/sessionreducer"
)

type Counter struct {
	count *atomic.Int32
}

func (c *Counter) SessionReduce(ctx context.Context, keys []string, input <-chan sessionreducer.Datum) sessionreducer.Messages {
	log.Println("SessionReduce invoked")
	for range input {
		c.count.Inc()
	}
	log.Println("SessionReduce done")
	return sessionreducer.MessagesBuilder().Append(sessionreducer.NewMessage([]byte(fmt.Sprintf("%d", c.count.Load()))).WithKeys(keys))
}

func (c *Counter) Accumulator(ctx context.Context) []byte {
	log.Println("Accumulator invoked")
	return []byte(strconv.Itoa(int(c.count.Load())))
}

func (c *Counter) MergeAccumulator(ctx context.Context, accumulator []byte) {
	log.Println("MergeAccumulator invoked")
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

func main() {
	sessionreducer.NewServer(NewSessionCounter).Start(context.Background())
}
