package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

type redisTestSink struct{}

// Sink This redis UDSink is created for numaflow e2e tests. This handle function assumes that
// a redis instance listening on address redis:6379 has already be up and running.
func (rds *redisTestSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		_ = d.EventTime()
		_ = d.Watermark()

		// We use redis hashes to store messages.
		// Each field of a hash is the content of a message and value of the field is the no. of occurrences of the message.
		var hashKey string
		if os.Getenv("NUMAFLOW_MONO_VERTEX_NAME") != "" {
			// if the mono vertex name environment variable is set, we are in a mono vertex
			// in this case, we use the mono vertex name as the hash key
			hashKey = os.Getenv("NUMAFLOW_MONO_VERTEX_NAME")
		} else {
			// if the mono vertex name environment variable is not set, we are in a pipeline
			// in this case, the name of a hash is pipelineName:sinkName.
			hashKey = fmt.Sprintf("%s:%s", os.Getenv("NUMAFLOW_PIPELINE_NAME"), os.Getenv("NUMAFLOW_VERTEX_NAME"))
		}
		err := client.HIncrBy(ctx, hashKey, string(d.Value()), 1).Err()
		if err != nil {
			log.Println("Set Error - ", err)
		} else {
			log.Printf("Incremented by 1 the no. of occurrences of %s under hash key %s\n", string(d.Value()), hashKey)
		}

		id := d.ID()
		result = result.Append(sinksdk.ResponseOK(id))
	}
	return result
}

func main() {
	err := sinksdk.NewServer(&redisTestSink{}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start sink function server: ", err)
	}
}
