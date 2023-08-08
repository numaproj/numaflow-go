package server

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sinksdk "github.com/KeranYang/numaflow-go/pkg/sink"
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

	sinkHandler := sinksdk.SinkFunc(func(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
		result := sinksdk.ResponsesBuilder()
		for d := range datumStreamCh {
			id := d.ID()
			if strings.Contains(string(d.Value()), "err") {
				result = result.Append(sinksdk.ResponseFailure(id, "mock sink message error"))
			} else {
				result = result.Append(sinksdk.ResponseOK(id))
			}

		}
		return result
	})
	// note: using actual UDS connection
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	err = New().RegisterSinker(sinkHandler).Start(ctx, WithSockAddr(socketFile.Name()), WithServerInfoFilePath(serverInfoFile.Name()))
	assert.NoError(t, err)
}

func TestSink_RegisterSinker(t *testing.T) {

	sinkHandler := sinksdk.SinkFunc(func(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
		result := sinksdk.ResponsesBuilder()
		for d := range datumStreamCh {
			id := d.ID()
			if strings.Contains(string(d.Value()), "err") {
				result = result.Append(sinksdk.ResponseFailure(id, "mock sink message error"))
			} else {
				result = result.Append(sinksdk.ResponseOK(id))
			}

		}
		return result
	})

	serv := New()
	assert.Nil(t, serv.svc.Sinker)
	serv.RegisterSinker(sinkHandler)
	assert.NotNil(t, serv.svc.Sinker)
}
