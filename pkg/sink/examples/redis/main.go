package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	sinksdk "github.com/numaproj/numaflow-go/pkg/sink"
	"github.com/numaproj/numaflow-go/pkg/sink/server"
)

// This redis UDSink is created for numaflow e2e tests. This handle function assumes that
// a redis instance listening on address redis:6379 has already be up and running.
func handle(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		_ = d.EventTime()
		_ = d.Watermark()

		// Our E2E tests time out after 10 minutes. Set redis message TTL to the same.
		const msgTTL = 10 * time.Minute
		err := client.Set(ctx, string(d.Value()), 1, msgTTL).Err()

		if err != nil {
			log.Println("Set Error - ", err)
		} else {
			log.Printf("Added key %s\n", string(d.Value()))
		}

		id := d.ID()
		result = result.Append(sinksdk.ResponseOK(id))
	}
	return result
}

func main() {
	server.New().RegisterSinker(sinksdk.SinkFunc(handle)).Start(context.Background())
}
