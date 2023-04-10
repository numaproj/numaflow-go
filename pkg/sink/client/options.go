package client

type options struct {
	sockAddr            string
	sereverInfoFilePath string
}

// Option is the interface to apply options.
type Option func(*options)

// WithSockAddr start the client with the given sock addr. This is mainly used for testing purpose.
func WithSockAddr(addr string) Option {
	return func(opts *options) {
		opts.sockAddr = addr
	}
}

// WithServerInfoFilePath start the client with the given server info file path. This is mainly used for testing purpose.
func WithServerInfoFilePath(f string) Option {
	return func(o *options) {
		o.sereverInfoFilePath = f
	}
}
