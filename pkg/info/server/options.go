package server

import (
	"time"

	"github.com/numaproj/numaflow-go/pkg/info"
)

type options struct {
	socketAdress string
	drainTimeout time.Duration
}

func defaultOptions() *options {
	return &options{
		socketAdress: info.SocketAddress,
		drainTimeout: 60 * time.Second,
	}
}

type Option func(*options)

// WithSocketAddress sets the socket address for the info server.
func WithSocketAddress(addr string) Option {
	return func(o *options) {
		o.socketAdress = addr
	}
}

// WithDrainTimeout sets the timeout for draining the info server.
func WithDrainTimeout(d time.Duration) Option {
	return func(o *options) {
		o.drainTimeout = d
	}
}
