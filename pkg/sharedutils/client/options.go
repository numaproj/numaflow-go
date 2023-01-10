package sharedutils

type Options struct {
	SockAddr string
}

// Option is the interface to apply options.
type Option func(*Options)

// WithSockAddr start the client with the given sock addr. This is mainly used for testing purpose.
func WithSockAddr(addr string) Option {
	return func(opts *Options) {
		opts.SockAddr = addr
	}
}
