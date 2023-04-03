package client

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/numaproj/numaflow-go/pkg/info"
	"github.com/numaproj/numaflow-go/pkg/info/server"
	"github.com/stretchr/testify/assert"
)

func Test_client(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	file, err := os.CreateTemp("/tmp", "test-info.sock")
	assert.NoError(t, err)
	defer func() {
		err = os.RemoveAll(file.Name())
		assert.NoError(t, err)
	}()

	go func() {
		if err := server.Start(ctx, server.WithSocketAddress(file.Name())); err != nil {
			t.Errorf("Start() error = %v", err)
		}
	}()

	c := NewInfoClient(WithSocketAddress(file.Name()))
	cctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()
	err = c.waitUntilReady(cctx)
	assert.NoError(t, err)
	si, err := c.GetServerInfo(ctx)
	assert.NoError(t, err)
	assert.Equal(t, si.Language, info.Go)
}
