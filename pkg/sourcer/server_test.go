package sourcer

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestNoopSource struct{}

func (ts TestNoopSource) Read(ctx context.Context, readRequest ReadRequest, messageCh chan<- Message) {
	return
}

func (ts TestNoopSource) Ack(ctx context.Context, request AckRequest) {
	return
}

func (ts TestNoopSource) Pending(ctx context.Context) uint64 {
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
	err := NewServer(TestNoopSource{}, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name())).Start(ctx)
	assert.NoError(t, err)
}
