package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

const defaultFailAfterSeconds = 60

// eotSink is a sinker that works normally during a warm-up period so the
// vertex can pass health / readiness checks, then starts producing an EOT
// failure by skipping responses for individual datums. That causes numaflow to
// see EOT before all expected responses have arrived
// (monovtx_critical_error_total{reason="eot_received_from_sink"} /
// forwarder_critical_error_total{reason="eot_received_from_sink"}).
//
// Configuration (environment variables):
//   - FAIL_AFTER_SECONDS  – start failing after this many seconds (default 60).
//     Set to 0 to disable the time-based trigger.
//   - FAIL_AFTER_MESSAGES – start failing after this many messages have been
//     successfully sunk. When unset or 0, only the time-based trigger is used.
//     When both are set, failure starts when either condition is met.
type eotSink struct {
	start             time.Time
	failAfterSeconds  int
	failAfterMessages int64
	processed         atomic.Int64
}

func newEotSink() *eotSink {
	return &eotSink{
		start:             time.Now(),
		failAfterSeconds:  envInt("FAIL_AFTER_SECONDS", defaultFailAfterSeconds),
		failAfterMessages: int64(envInt("FAIL_AFTER_MESSAGES", 0)),
	}
}

func (s *eotSink) shouldFail() bool {
	if s.failAfterMessages > 0 && s.processed.Load() >= s.failAfterMessages {
		return true
	}
	if s.failAfterSeconds > 0 && time.Since(s.start) >= time.Duration(s.failAfterSeconds)*time.Second {
		return true
	}
	return false
}

func (s *eotSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		if s.shouldFail() {
			// Intentionally skip the response for this datum while still
			// draining the channel so the SDK can send EOT with missing replies.
			continue
		}
		result = result.Append(sinksdk.ResponseOK(d.ID()))
		s.processed.Add(1)
	}
	return result
}

func main() {
	sink := newEotSink()
	log.Printf("eot sink starting: fail_after_seconds=%d fail_after_messages=%d",
		sink.failAfterSeconds, sink.failAfterMessages)
	err := sinksdk.NewServer(sink).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start eot sink server: ", err)
	}
}

func envInt(key string, defaultVal int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("%s=%q is not an int, using default %d", key, raw, defaultVal)
		return defaultVal
	}
	return val
}
