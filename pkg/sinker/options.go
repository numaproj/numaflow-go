package sinker

import (
	"github.com/numaproj/numaflow-go/pkg/info"
	"github.com/numaproj/numaflow-go/pkg/shared"
)

type options struct {
	sockAddr           string
	maxMessageSize     int
	serverInfoFilePath string
}

// Option is the interface to apply options.
type Option func(*options)

func DefaultOptions() *options {
	return &options{
		sockAddr:           shared.SinkAddr,
		maxMessageSize:     shared.DefaultMaxMessageSize,
		serverInfoFilePath: info.ServerInfoFilePath,
	}
}

// WithMaxMessageSize sets the sinkServer max receive message size and the sinkServer max send message size to the given size.
func WithMaxMessageSize(size int) Option {
	return func(opts *options) {
		opts.maxMessageSize = size
	}
}

// WithSockAddr start the sinkServer with the given sock addr. This is mainly used for testing purpose.
func WithSockAddr(addr string) Option {
	return func(opts *options) {
		opts.sockAddr = addr
	}
}

// WithServerInfoFilePath sets the sinkServer info file path.
func WithServerInfoFilePath(path string) Option {
	return func(opts *options) {
		opts.serverInfoFilePath = path
	}
}
