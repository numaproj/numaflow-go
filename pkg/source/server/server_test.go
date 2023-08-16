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

type TestSource struct{}

func (ts TestSource) Read(ctx context.Context, readRequest sourcesdk.ReadRequest, messageCh chan<- model.Message) {
	return
}

func (ts TestSource) Ack(ctx context.Context, request sourcesdk.AckRequest) {
	return
}

func (ts TestSource) Pending(ctx context.Context) uint64 {
	return 0
}

func TestServer_Start(t *testing.T) {
	socketFile, _ := os.CreateTemp("/tmp", "numaflow-test.sock")
	defer func() {
		_ = os.RemoveAll(socketFile.Name())
	}()

	serverInfoFile, _ := os.CreateTemp("/tmp", "numaflow-test-info")
	defer func() {
		_ = os.RemoveAll(serverInfoFile.Name())
	}()

	// note: using actual UDS connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	err := New(TestSource{}).Start(ctx, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name()))
	assert.NoError(t, err)
}
