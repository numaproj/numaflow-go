package info

type options struct {
	svrInfoFilePath string
}

func defaultOptions() *options {
	return &options{
		svrInfoFilePath: ServerInfoFilePath,
	}
}

type Option func(*options)

// WithServerInfoFilePath sets the server info file path
func WithServerInfoFilePath(f string) Option {
	return func(o *options) {
		o.svrInfoFilePath = f
	}
}
