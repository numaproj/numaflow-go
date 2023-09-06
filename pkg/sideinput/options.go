package sideinput

// options is the struct to hold the server options.
type options struct {
	sockAddr       string
	maxMessageSize int
}

// Option is the interface to apply options.
type Option func(*options)

// DefaultOptions returns the default options.
func DefaultOptions() *options {
	return &options{
		sockAddr:       Address,
		maxMessageSize: DefaultMaxMessageSize,
	}
}

// WithMaxMessageSize sets the server max receive message size and the server max send message size to the given size.
func WithMaxMessageSize(size int) Option {
	return func(opts *options) {
		opts.maxMessageSize = size
	}
}

// WithSockAddr start the server with the given sock addr. This is mainly used for testing purpose.
func WithSockAddr(addr string) Option {
	return func(opts *options) {
		opts.sockAddr = addr
	}
}
