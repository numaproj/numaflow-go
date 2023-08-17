package mapstreamer

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
		sockAddr:           shared.MapStreamAddr,
		maxMessageSize:     shared.DefaultMaxMessageSize,
		serverInfoFilePath: info.ServerInfoFilePath,
	}
}

// WithMaxMessageSize sets the server max receive message size and the server max send message size to the given size.
func WithMaxMessageSize(size int) Option {
	return func(opts *options) {
		opts.maxMessageSize = size
	}
}

// WithSockAddr start the server with the given sock addr. This is mainly used for testing purposes.
func WithSockAddr(addr string) Option {
	return func(opts *options) {
		opts.sockAddr = addr
	}
}

// WithServerInfoFilePath sets the server info file path to the given path.
func WithServerInfoFilePath(f string) Option {
	return func(opts *options) {
		opts.serverInfoFilePath = f
	}
}
