package sourcer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithMaxMessageSize(t *testing.T) {
	var (
		testSize = 1024 * 1024 * 10
		opts     = &options{
			maxMessageSize: defaultMaxMessageSize,
		}
	)
	WithMaxMessageSize(testSize)(opts)
	assert.Equal(t, testSize, opts.maxMessageSize)
}

func TestWithSockAddr(t *testing.T) {
	var (
		testSocketAddr = "test-socket-address"
		opts           = &options{
			sockAddr: address,
		}
	)
	WithSockAddr(testSocketAddr)(opts)
	assert.Equal(t, testSocketAddr, opts.sockAddr)
}

func TestWithServerInfoFilePath(t *testing.T) {
	var (
		testServerInfoFilePath = "test-server-info-file-path"
		opts                   = &options{
			maxMessageSize: defaultMaxMessageSize,
		}
	)
	WithServerInfoFilePath(testServerInfoFilePath)(opts)
	assert.Equal(t, testServerInfoFilePath, opts.serverInfoFilePath)
}
