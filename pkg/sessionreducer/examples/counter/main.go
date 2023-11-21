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
	log.Println("SessionReduce called")
	for range input {
		println("incrementing the counter")
		c.count.Inc()
	}
	log.Println("SessionReduce done")
	return sessionreducer.MessagesBuilder().Append(sessionreducer.NewMessage([]byte(fmt.Sprintf("%d", c.count.Load()))).WithKeys(keys))
}

func (c *Counter) Accumulator(ctx context.Context) []byte {
	log.Println("Accumulator called")
	println("returning the accumulator value - ", c.count.Load())
	return []byte(strconv.Itoa(int(c.count.Load())))
}

func (c *Counter) MergeAccumulator(ctx context.Context, accumulator []byte) {
	log.Println("MergeAccumulator called")
	val, err := strconv.Atoi(string(accumulator))
	println("merging the accumulator with value - ", val)
	if err != nil {
		log.Println("unable to convert the accumulator value to int: ", err)
		return
	}
	c.count.Add(int32(val))
}

func reduceCounter() sessionreducer.SessionReducer {
	return &Counter{
		count: atomic.NewInt32(0),
	}
}

func main() {
	sessionreducer.NewServer(reduceCounter).Start(context.Background())
}
