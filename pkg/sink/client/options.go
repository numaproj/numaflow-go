package client

type options struct {
	sockAddr       string
	maxMessageSize int
}

// Option is the interface to apply options.
type Option func(*options)

// WithMaxMessageSize sets the client max receive message size and the client max send message size to the given size.
func WithMaxMessageSize(size int) Option {
	return func(opts *options) {
		opts.maxMessageSize = size
	}
}

// WithSockAddr start the client with the given sock addr. This is mainly used for testing purpose.
func WithSockAddr(addr string) Option {
	return func(opts *options) {
		opts.sockAddr = addr
	}
}
