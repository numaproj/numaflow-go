package reducer

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

	var createReduceHandler = func() Reducer {
		return ReducerFunc(func(ctx context.Context, keys []string, rch <-chan Datum, md Metadata) Messages {
			sum := 0
			for val := range rch {
				msgVal, _ := strconv.Atoi(string(val.Value()))
				sum += msgVal
			}
			return MessagesBuilder().Append(NewMessage([]byte(strconv.Itoa(sum))).WithKeys([]string{keys[0] + "_test"}))
		})
	}
	// note: using actual uds connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err := NewServer(ReducerCreatorFn(createReduceHandler), WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name())).Start(ctx)
	assert.NoError(t, err)
}
