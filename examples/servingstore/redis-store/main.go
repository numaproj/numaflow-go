package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/numaproj/numaflow-go/pkg/servingstore"
	"github.com/redis/go-redis/v9"
)

const DEFAULT_REDIS_URL = "redis:6379"
const DEFAULT_REDIS_TTL_SECONDS = 7200 // 2 hours

type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
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
	slog.Info("Received Put request", "key", id)
	payloads := putDatum.Payloads()
	encodedPayloads := make([]any, 0, len(payloads))
	for _, payload := range payloads {
		encoded, err := encodeServingPayload(payload)
		if err != nil {
			slog.Error("Encoding redis payload", "error", err)
			os.Exit(1)
		}
		encodedPayloads = append(encodedPayloads, encoded)
	}
	_, err := rs.client.LPush(ctx, id, encodedPayloads...).Result()
	if err != nil {
		slog.Error("Saving payloads with LPUSH", "key", id, "error", err)
		os.Exit(1)
	}
	if _, err := rs.client.Expire(ctx, id, rs.ttl).Result(); err != nil {
		slog.Error("Setting expiry for redis key", "key", id, "error", err)
		os.Exit(1)
	}
	slog.Info("Saved payloads", "key", id, "count", len(encodedPayloads))
}

func (rs *RedisStore) Get(ctx context.Context, getDatum servingstore.GetDatum) servingstore.StoredResult {
	id := getDatum.ID()
	slog.Info("Received Get request", "key", id)
	values, err := rs.client.LRange(ctx, id, 0, -1).Result()
	if err != nil {
		slog.Error("Retrieving results", "key", id, "error", err)
		os.Exit(1)
	}
	if len(values) == 0 {
		slog.Info("Returning empty results", "id", id)
		return servingstore.NewStoredResult(id, nil)
	}
	payloads := make([]servingstore.Payload, 0, len(values))
	for _, value := range values {
		payload, err := decodeServingPayload(value)
		if err != nil {
			slog.Error("Decoding redis payload", "error", err)
			os.Exit(1)
		}
		payloads = append(payloads, servingstore.NewPayload(payload.Origin, payload.Value))
	}
	slog.Info("Returning results", "key", id, "count", len(payloads))
	return servingstore.NewStoredResult(id, payloads)
}

func NewRedisStore(addr string, ttl time.Duration) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisStore{client: rdb, ttl: ttl}
}

func main() {
	redisURL := DEFAULT_REDIS_URL
	if addr, exists := os.LookupEnv("REDIS_ADDR"); exists {
		redisURL = addr
	}
	redisTTL := DEFAULT_REDIS_TTL_SECONDS * time.Second
	if ttl, exists := os.LookupEnv("REDIS_TTL_SECONDS"); exists {
		ttlSecs, err := strconv.ParseInt(ttl, 10, 64)
		if err != nil {
			slog.Error("Converting value of env variable REDIS_TTL_SECONDS to integer:", "error", err)
		} else {
			redisTTL = time.Duration(ttlSecs) * time.Second
		}
	}

	slog.Info("Starting Redis serving store", "redis_url", redisURL, "ttl", redisTTL)
	err := servingstore.NewServer(NewRedisStore(redisURL, redisTTL)).Start(context.Background())
	if err != nil {
		slog.Error("Failed to serving store function server", "error", err)
		os.Exit(1)
	}
}
