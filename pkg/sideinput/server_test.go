package sideinput

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestServer_Start tests the Start method to check whether the grpc server is correctly started.
func TestServer_Start(t *testing.T) {
	socketFile, _ := os.CreateTemp("/tmp", "numaflow-sideinput")
	defer func() {
		_ = os.RemoveAll(socketFile.Name())
	}()

	var retrieveHandler = RetrieveFunc(func(ctx context.Context) Message {
		return NewMessage([]byte("test"))
	})
	// note: using actual UDS connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err := NewSideInputServer(retrieveHandler, WithSockAddr(socketFile.Name())).Start(ctx)
	assert.NoError(t, err)
}
