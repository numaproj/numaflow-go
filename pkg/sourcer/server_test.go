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
}

func (ts TestNoopSource) Ack(ctx context.Context, request AckRequest) {
}

func (ts TestNoopSource) Pending(ctx context.Context) int64 {
	return 0
}

func (ts TestNoopSource) Partitions(ctx context.Context) []int32 {
	return []int32{0}
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

	// note: using actual uds connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err := NewServer(TestNoopSource{}, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name())).Start(ctx)
	assert.NoError(t, err)
}
