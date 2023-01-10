package sharedutils

import (
	"testing"

	"github.com/numaproj/numaflow-go/pkg/function"
	"github.com/stretchr/testify/assert"
)

func TestWithMaxMessageSize(t *testing.T) {
	var (
		size = 1024 * 1024 * 10
		opts = &Options{
			MaxMessageSize: function.DefaultMaxMessageSize,
		}
	)
	WithMaxMessageSize(1024 * 1024 * 10)(opts)
	assert.Equal(t, size, opts.MaxMessageSize)
}
