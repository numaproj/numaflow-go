package main

import (
	"context"
	"math/rand"
	"time"

	functionsdk "github.com/numaproj/numaflow-go/pkg/function"
	"github.com/numaproj/numaflow-go/pkg/function/server"
	"github.com/numaproj/numaflow-go/pkg/function/types"
)

func mapTHandle(_ context.Context, key string, d functionsdk.Datum) types.MessageTs {
	// assign a random event time then forward the input to the output.
	randomEventTime := generateRandomTime
	return types.MessageTsBuilder().Append(types.MessageTTo(randomEventTime(), key, d.Value()))
}

// generateRandomTime generates a random timestamp within date range [1970-01-01 to 2023-01-01]
func generateRandomTime() time.Time {
	min := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func main() {
	server.New().RegisterMapperT(functionsdk.MapTFunc(mapTHandle)).Start(context.Background())
}
