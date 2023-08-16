package reduce

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReduceServer_Start(t *testing.T) {
	socketFile, _ := os.CreateTemp("/tmp", "numaflow-test.sock")
	defer func() {
		_ = os.RemoveAll(socketFile.Name())
	}()

	serverInfoFile, _ := os.CreateTemp("/tmp", "numaflow-test-info")
	defer func() {
		_ = os.RemoveAll(serverInfoFile.Name())
	}()

	var reduceHandler = ReducerFunc(func(ctx context.Context, keys []string, rch <-chan Datum, md Metadata) Messages {
		sum := 0
		for val := range rch {
			msgVal, _ := strconv.Atoi(string(val.Value()))
			sum += msgVal
		}
		return MessagesBuilder().Append(NewMessage([]byte(strconv.Itoa(sum))).WithKeys([]string{keys[0] + "_test"}))
	})
	// note: using actual UDS connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	err := NewServer(reduceHandler, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name())).Start(ctx)
	assert.NoError(t, err)
}
