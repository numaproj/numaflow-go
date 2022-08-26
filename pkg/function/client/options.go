package client

type options struct {
	sockAddr string
}

// Option is the interface to apply options.
type Option func(*options)

// WithSockAddr start the server with the given sock addr. This is mainly used for testing purpose.
func WithSockAddr(addr string) Option {
	return func(opts *options) {
		opts.sockAddr = addr
	}
}
