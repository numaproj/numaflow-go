package sinker

import (
	"os"

	"github.com/numaproj/numaflow-go/internal/shared"
)

type options struct {
	sockAddr           string
	maxMessageSize     int
	serverInfoFilePath string
}

// Option configures server startup options.
type Option func(*options)

func defaultOptions() *options {
	defaultPath := serverInfoFilePath
	defaultAddress := address

	// If the container type is fallback sink, then use the fallback sink address and path.
	if os.Getenv(shared.EnvUDContainerType) == UDContainerFallbackSink {
		defaultPath = fbServerInfoFilePath
		defaultAddress = fbAddress
	} else if os.Getenv(shared.EnvUDContainerType) == UDContainerOnSuccessSink {
		defaultPath = onSuccessServerInfoFilePath
		defaultAddress = onSuccessAddress
	}

	return &options{
		sockAddr:           defaultAddress,
		maxMessageSize:     defaultMaxMessageSize,
		serverInfoFilePath: defaultPath,
	}
}

// WithMaxMessageSize sets the sinkServer max receive message size and the sinkServer max send message size to the given size.
func WithMaxMessageSize(size int) Option {
	return func(opts *options) {
		opts.maxMessageSize = size
	}
}

// WithSockAddr starts the server on the given socket address. This is mainly used for testing purposes.
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
