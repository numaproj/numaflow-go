package client

import (
	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1/funcmock"
)

type options struct {
	mockClient gRPClientOption
}

// Option is the interface to apply options.
type Option interface {
	apply(*options)
}

type gRPClientOption struct {
	mockClnt *funcmock.MockUserDefinedFunctionClient
}

func (c gRPClientOption) apply(opts *options) {
	opts.mockClient = gRPClientOption{
		c.mockClnt,
	}
}

// WithMockGRPCClient creates a new client object with the given mock client for mock testing.
func WithMockGRPCClient(c *funcmock.MockUserDefinedFunctionClient) Option {
	return gRPClientOption{mockClnt: c}
}
