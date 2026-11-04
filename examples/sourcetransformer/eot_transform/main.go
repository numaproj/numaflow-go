package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/numaproj/numaflow-go/pkg/sourcetransformer"
)

const defaultFailAfterSeconds = 60

// eotTransform is a source transformer that works normally during a warm-up
// period so the vertex can pass health / readiness checks, then starts
// producing an EOT failure by returning an empty response (no messages and no
// DROP tag). That causes numaflow to treat the reply as a missing response
// (monovtx_critical_error_total{reason="eot_received_from_transformer"} /
// forwarder_critical_error_total{reason="eot_received_from_transformer"}).
//
// Configuration (environment variables):
//   - FAIL_AFTER_SECONDS  – start failing after this many seconds (default 60).
//     Set to 0 to disable the time-based trigger.
//   - FAIL_AFTER_MESSAGES – start failing after this many messages have been
//     successfully transformed. When unset or 0, only the time-based trigger
//     is used. When both are set, failure starts when either condition is met.
type eotTransform struct {
	start             time.Time
	failAfterSeconds  int
	failAfterMessages int64
	processed         atomic.Int64
}

func newEotTransform() *eotTransform {
	return &eotTransform{
		start:             time.Now(),
		failAfterSeconds:  envInt("FAIL_AFTER_SECONDS", defaultFailAfterSeconds),
		failAfterMessages: int64(envInt("FAIL_AFTER_MESSAGES", 0)),
	}
}

func (t *eotTransform) shouldFail() bool {
	if t.failAfterMessages > 0 && t.processed.Load() >= t.failAfterMessages {
		return true
	}
	if t.failAfterSeconds > 0 && time.Since(t.start) >= time.Duration(t.failAfterSeconds)*time.Second {
		return true
	}
	return false
}

func (t *eotTransform) Transform(_ context.Context, keys []string, d sourcetransformer.Datum) sourcetransformer.Messages {
	if t.shouldFail() {
		// Empty response (not MessageToDrop) triggers eot_received_from_transformer.
		return sourcetransformer.MessagesBuilder()
	}

	t.processed.Add(1)
	return sourcetransformer.MessagesBuilder().Append(
		sourcetransformer.NewMessage(d.Value(), d.EventTime()).WithKeys(keys),
	)
}

func main() {
	transformer := newEotTransform()
	log.Printf("eot transform starting: fail_after_seconds=%d fail_after_messages=%d",
		transformer.failAfterSeconds, transformer.failAfterMessages)
	err := sourcetransformer.NewServer(transformer).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start eot transform server: ", err)
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
