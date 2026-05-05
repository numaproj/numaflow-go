package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

// crashSink succeeds for each message but panics periodically (wall clock),
// for chaos / failure testing of Numaflow and the Go SDK.
type crashSink struct {
	mu       sync.Mutex
	interval time.Duration
	// lastPanic is the time of the last panic (or the anchor time before any panic has fired).
	lastPanic time.Time
}

func crashIntervalFromEnv() time.Duration {
	s := strings.TrimSpace(os.Getenv("CRASH_INTERVAL"))
	if s == "" {
		return 5 * time.Minute
	}
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		log.Fatalf("invalid CRASH_INTERVAL %q: %v", s, err)
	}
	return d
}

func (c *crashSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		if ctx.Err() != nil {
			return result
		}

		now := time.Now()
		c.mu.Lock()
		if c.lastPanic.IsZero() {
			c.lastPanic = now
		} else if now.Sub(c.lastPanic) >= c.interval {
			c.lastPanic = now
			c.mu.Unlock()
			panic(fmt.Sprintf("crash_at_frequency: periodic panic (CRASH_INTERVAL=%s)", c.interval.String()))
		}
		c.mu.Unlock()

		result = result.Append(sinksdk.ResponseOK(d.ID()))
	}
	return result
}

func main() {
	s := &crashSink{interval: crashIntervalFromEnv()}
	err := sinksdk.NewServer(s).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start crash-at-frequency sink server: ", err)
	}
}
