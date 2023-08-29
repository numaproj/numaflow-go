package mapstreamer

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMapStreamServer_Start(t *testing.T) {
	socketFile, _ := os.CreateTemp("/tmp", "numaflow-test.sock")
	defer func() {
		_ = os.RemoveAll(socketFile.Name())
	}()

	serverInfoFile, _ := os.CreateTemp("/tmp", "numaflow-test-info")
	defer func() {
		_ = os.RemoveAll(serverInfoFile.Name())
	}()

	var mapStreamHandler = MapStreamerFunc(func(ctx context.Context, keys []string, datum Datum, messageCh chan<- Message) {
		msg := datum.Value()
		messageCh <- NewMessage(msg).WithKeys([]string{keys[0] + "_test"})
		close(messageCh)
	})
	// note: using actual uds connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err := NewServer(mapStreamHandler, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name())).Start(ctx)
	assert.NoError(t, err)
}
