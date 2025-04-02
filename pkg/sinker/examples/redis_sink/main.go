package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

type RedisTestSink struct {
	hashKey          string
	messageCount     int
	inflightMessages []sinksdk.Datum
	client           *redis.Client
	checkOrder       bool
}

// NewRedisTestSink creates a new instance of redisTestSink with a Redis client.
func NewRedisTestSink() *RedisTestSink {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	hashKey := os.Getenv("SINK_HASH_KEY")
	if hashKey == "" {
		log.Panicf("SINK_HASH_KEY environment variable is not set.")
	}

	messageCountStr := os.Getenv("MESSAGE_COUNT")
	messageCount := 0
	if messageCountStr != "" {
		messageCount, _ = strconv.Atoi(messageCountStr)
	}

	checkOrderStr := os.Getenv("CHECK_ORDER")
	checkOrder := false
	if checkOrderStr != "" {
		checkOrder, _ = strconv.ParseBool(checkOrderStr)
	}

	return &RedisTestSink{
		client:           client,
		hashKey:          hashKey,
		messageCount:     messageCount,
		inflightMessages: make([]sinksdk.Datum, 0, messageCount),
		checkOrder:       checkOrder,
	}
}

// Sink This redis UDSink is created for numaflow e2e tests. This handle function assumes that
// a redis instance listening on address redis:6379 has already be up and running.
func (rds *RedisTestSink) Sink(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		if rds.checkOrder {
			rds.inflightMessages = append(rds.inflightMessages, d)
			if len(rds.inflightMessages) == rds.messageCount {
				var ordered = true
				for i := 1; i < len(rds.inflightMessages); i++ {
					if rds.inflightMessages[i].EventTime().Before(rds.inflightMessages[i-1].EventTime()) {
						ordered = false
						break
					}
				}

				var resultMessage string
				if ordered {
					resultMessage = "ordered"
				} else {
					resultMessage = "not ordered"
				}

				err := rds.client.HIncrBy(ctx, rds.hashKey, resultMessage, 1).Err()
				if err != nil {
					log.Println("Set Error - ", err)
				} else {
					log.Printf("Incremented by 1 the no. of occurrences of %s under hash key %s\n", resultMessage, rds.hashKey)
				}
				rds.inflightMessages = make([]sinksdk.Datum, 0, rds.messageCount)
			}
		}

		// watermark and event time of the message can be accessed
		_ = d.EventTime()
		_ = d.Watermark()

		// We use redis hashes to store messages.
		// Each field of a hash is the content of a message and value of the field is the no. of occurrences of the message.
		err := rds.client.HIncrBy(ctx, rds.hashKey, string(d.Value()), 1).Err()
		if err != nil {
			log.Println("Set Error - ", err)
		} else {
			log.Printf("Incremented by 1 the no. of occurrences of %s under hash key %s\n", string(d.Value()), rds.hashKey)
		}
		id := d.ID()
		result = result.Append(sinksdk.ResponseOK(id))
	}
	return result
}

func main() {
	sink := NewRedisTestSink()
	err := sinksdk.NewServer(sink).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start sink function server: ", err)
	}
}
