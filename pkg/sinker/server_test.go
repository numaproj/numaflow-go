package sinker

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSink_Start(t *testing.T) {
	socketFile, err := os.CreateTemp("/tmp", "numaflow-test.sock")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(socketFile.Name())
		assert.NoError(t, err)
	}()

	serverInfoFile, err := os.CreateTemp("/tmp", "numaflow-test-info")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(serverInfoFile.Name())
		assert.NoError(t, err)
	}()

	sinkHandler := SinkerFunc(func(ctx context.Context, datumStreamCh <-chan Datum) Responses {
		result := ResponsesBuilder()
		for d := range datumStreamCh {
			id := d.ID()
			if strings.Contains(string(d.Value()), "err") {
				result = result.Append(ResponseFailure(id, "mock sink message error"))
			} else {
				result = result.Append(ResponseOK(id))
			}

		}
		return result
	})
	// note: using actual uds connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = NewServer(sinkHandler, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name())).Start(ctx)
	assert.NoError(t, err)
}
