package server

import (
	"testing"

	sinksdk "github.com/numaproj/numaflow-go/pkg/sinkv2"
	"github.com/stretchr/testify/assert"
)

func TestWithMaxMessageSize(t *testing.T) {
	var (
		size = 1024 * 1024 * 10
		opts = &options{
			maxMessageSize: sinksdk.DefaultMaxMessageSize,
		}
	)
	WithMaxMessageSize(1024 * 1024 * 10)(opts)
	assert.Equal(t, size, opts.maxMessageSize)
}
