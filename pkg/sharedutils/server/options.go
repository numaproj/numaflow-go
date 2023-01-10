package sharedutils

type Options struct {
	SockAddr       string
	MaxMessageSize int
}

// Option is the interface to apply options.
type Option func(*Options)

// WithMaxMessageSize sets the server max receive message size and the server max send message size to the given size.
func WithMaxMessageSize(size int) Option {
	return func(opts *Options) {
		opts.MaxMessageSize = size
	}
}

// WithSockAddr start the server with the given sock addr. This is mainly used for testing purpose.
func WithSockAddr(addr string) Option {
	return func(opts *Options) {
		opts.SockAddr = addr
	}
}
