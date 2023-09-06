package sideinput

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWithMaxMessageSize tests the WithMaxMessageSize option.
// It should set the max message size in the options struct.
func TestWithMaxMessageSize(t *testing.T) {
	var (
		size = 1024 * 1024 * 10
		opts = &options{
			maxMessageSize: DefaultMaxMessageSize,
		}
	)
	WithMaxMessageSize(1024 * 1024 * 10)(opts)
	assert.Equal(t, size, opts.maxMessageSize)
}
