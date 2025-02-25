package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/numaproj/numaflow-go/pkg/servingstore"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

type payload struct {
	Origin string `json:"origin"`
	Value  []byte `json:"value"`
}

func encodeServingPayload(data servingstore.Payload) ([]byte, error) {
	payloadData := payload{
		Origin: data.Origin(),
		Value:  data.Value(),
	}
	return json.Marshal(&payloadData)
}

func decodeServingPayload(data string) (payload, error) {
	var payloadData payload
	if err := json.Unmarshal([]byte(data), &payloadData); err != nil {
		return payloadData, fmt.Errorf("unmarshaling payload: %w", err)
	}
	return payloadData, nil
}

func (rs *RedisStore) Put(ctx context.Context, putDatum servingstore.PutDatum) {
	id := putDatum.ID()
	log.Printf("Received Put request for %s", id)
	payloads := putDatum.Payloads()
	encodedPayloads := make([]any, 0, len(payloads))
	for _, payload := range payloads {
		encoded, err := encodeServingPayload(payload)
		if err != nil {
			log.Fatal(err)
		}
		encodedPayloads = append(encodedPayloads, encoded)
	}
	_, err := rs.client.LPush(ctx, id, encodedPayloads...).Result()
	if err != nil {
		log.Fatalf("Saving payloads with LPUSH: key=%s, err=%v", id, err)
	}
	log.Printf("Saved payloads. key=%s, count=%d", id, len(encodedPayloads))
}

func (rs *RedisStore) Get(ctx context.Context, getDatum servingstore.GetDatum) servingstore.StoredResult {
	id := getDatum.ID()
	log.Printf("Received Get request for %s", id)
	values, err := rs.client.LRange(ctx, id, 0, -1).Result()
	if err != nil {
		log.Fatalf("Retrieving results id=%s: %v", id, err)
	}
	if len(values) == 0 {
		log.Printf("Returning empty results id=%s", id)
		return servingstore.NewStoredResult(id, nil)
	}
	payloads := make([]servingstore.Payload, 0, len(values))
	for _, value := range values {
		payload, err := decodeServingPayload(value)
		if err != nil {
			log.Fatal(err)
		}
		payloads = append(payloads, servingstore.NewPayload(payload.Origin, payload.Value))
	}
	log.Printf("Returning results id=%s count=%d", id, len(payloads))
	return servingstore.NewStoredResult(id, payloads)
}

func NewRedisStore() *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisStore{client: rdb}
}

func main() {
	err := servingstore.NewServer(NewRedisStore()).Start(context.Background())
	if err != nil {
		log.Panic("Failed to serving store function server: ", err)
	}
}
