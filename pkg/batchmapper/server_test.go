package batchmapper

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBatchMapServer_Start(t *testing.T) {
	socketFile, _ := os.CreateTemp("/tmp", "numaflow-test.sock")
	defer func() {
		_ = os.RemoveAll(socketFile.Name())
	}()

	serverInfoFile, _ := os.CreateTemp("/tmp", "numaflow-test-info")
	defer func() {
		_ = os.RemoveAll(serverInfoFile.Name())
	}()

	var batchMapHandler = BatchMapperFunc(func(ctx context.Context, datums <-chan Datum) BatchResponses {
		batchResponses := BatchResponsesBuilder()
		for d := range datums {
			results := NewBatchResponse(d.ID())
			results.Append(NewMessage(d.Value()).WithKeys([]string{d.Keys()[0] + "_test"}))
			batchResponses.Append(results)
		}
		return batchResponses
	})
	// note: using actual uds connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err := NewServer(batchMapHandler, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name())).Start(ctx)
	assert.NoError(t, err)
}
