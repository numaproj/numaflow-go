package server

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sourcesdk "github.com/numaproj/numaflow-go/pkg/source"
	"github.com/numaproj/numaflow-go/pkg/source/model"
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

	pendingHandler := sourcesdk.PendingFunc(func(ctx context.Context) uint64 {
		return 1
	})
	readHandler := sourcesdk.ReadFunc(func(ctx context.Context, readRequest sourcesdk.ReadRequest, messageCh chan<- model.Message) {
		return
	})

	// note: using actual UDS connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	err := New(pendingHandler, readHandler).Start(ctx, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name()))
	assert.NoError(t, err)
}
