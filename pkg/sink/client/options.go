package client

type options struct {
	sockAddr        string
	infoSvrSockAddr string
}

// Option is the interface to apply options.
type Option func(*options)

// WithSockAddr start the client with the given sock addr. This is mainly used for testing purpose.
func WithSockAddr(addr string) Option {
	return func(opts *options) {
		opts.sockAddr = addr
	}
}

// WithInfoServerSocketAddr start the client with the given info server sock addr. This is mainly used for testing purpose.
func WithInfoServerSocketAddr(addr string) Option {
	return func(opts *options) {
		opts.infoSvrSockAddr = addr
	}
}
