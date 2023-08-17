package mapstreamer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/numaproj/numaflow-go/pkg/shared"
)

func TestWithMaxMessageSize(t *testing.T) {
	var (
		size = 1024 * 1024 * 10
		opts = &options{
			maxMessageSize: shared.DefaultMaxMessageSize,
		}
	)
	WithMaxMessageSize(1024 * 1024 * 10)(opts)
	assert.Equal(t, size, opts.maxMessageSize)
}
