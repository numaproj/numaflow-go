package server

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sideinputsdk "github.com/numaproj/numaflow-go/pkg/sideinput"
)

// TestServer_Start tests the Start method to check whether the grpc server is correctly started.
func TestServer_Start(t *testing.T) {
	socketFile, _ := os.CreateTemp("/tmp", "numaflow-test.sock")
	defer func() {
		_ = os.RemoveAll(socketFile.Name())
	}()

	serverInfoFile, _ := os.CreateTemp("/tmp", "numaflow-test-info")
	defer func() {
		_ = os.RemoveAll(serverInfoFile.Name())
	}()
	var retrieveHandler = sideinputsdk.RetrieveSideInput(func(ctx context.Context) sideinputsdk.MessageSI {
		return sideinputsdk.NewMessageSI([]byte("test"))
	})
	// note: using actual UDS connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	err := NewSideInputServer().RegisterRetriever(retrieveHandler).Start(ctx, WithSockAddr(socketFile.Name()))
	assert.NoError(t, err)
}
