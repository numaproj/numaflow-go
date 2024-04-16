package sinker

import "os"

type options struct {
	sockAddr           string
	maxMessageSize     int
	serverInfoFilePath string
}

// Option is the interface to apply options.
type Option func(*options)

func defaultOptions() *options {
	defaultPath := serverInfoFilePath
	defaultAddress := address

	// If the container type is fallback sink, then use the fallback sink address and path.
	if os.Getenv(EnvUDContainerType) == UDContainerFallbackSink {
		defaultPath = fbServerInfoFilePath
		defaultAddress = fbAddress
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
