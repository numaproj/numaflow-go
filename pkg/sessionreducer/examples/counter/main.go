package main

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/atomic"

	"github.com/numaproj/numaflow-go/pkg/sessionreducer"
)

type Counter struct {
	count atomic.Int64
}

func (c *Counter) SessionReduce(ctx context.Context, keys []string, input <-chan sessionreducer.Datum) sessionreducer.Messages {
	log.Println("SessionReduce called")
	for range input {
		c.count.Inc()
	}
	log.Println("SessionReduce done")
	return sessionreducer.MessagesBuilder().Append(sessionreducer.NewMessage([]byte(fmt.Sprintf("%d", c.count.Load()))))
}

func (c *Counter) Accumulator(ctx context.Context) []byte {
	log.Println("Accumulator called")
	return []byte(fmt.Sprintf("%d", c.count.Load()))
}

func (c *Counter) MergeAccumulator(ctx context.Context, accumulator []byte) {
	log.Println("MergeAccumulator called")
	c.count.Add(int64(len(accumulator)))
}

func reduceCounter() sessionreducer.SessionReducer {
	return &Counter{}
}

func main() {
	sessionreducer.NewServer(reduceCounter).Start(context.Background())
}
