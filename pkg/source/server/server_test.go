package server

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sourcesdk "github.com/KeranYang/numaflow-go/pkg/source"
)

func TestServer_Start(t *testing.T) {
	socketFile, _ := os.CreateTemp("/tmp", "numaflow-test.sock")
	defer func() {
		_ = os.RemoveAll(socketFile.Name())
	}()

	serverInfoFile, _ := os.CreateTemp("/tmp", "numaflow-test-info")
	defer func() {
		_ = os.RemoveAll(serverInfoFile.Name())
	}()

	var pendingHandler = sourcesdk.PendingFunc(func(ctx context.Context) uint64 {
		return 1
	})
	// note: using actual UDS connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	err := New().RegisterPendingHandler(pendingHandler).Start(ctx, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name()))
	assert.NoError(t, err)
}
