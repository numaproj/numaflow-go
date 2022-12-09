package server

import (
	"testing"

	"github.com/numaproj/numaflow-go/pkg/sink"
	"github.com/stretchr/testify/assert"
)

func TestWithMaxMessageSize(t *testing.T) {
	var (
		size = 1024 * 1024 * 10
		opts = &options{
			maxMessageSize: sink.DefaultMaxMessageSize,
		}
	)
	WithMaxMessageSize(1024 * 1024 * 10)(opts)
	assert.Equal(t, size, opts.maxMessageSize)
}
