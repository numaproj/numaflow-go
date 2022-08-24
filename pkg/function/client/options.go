package client

import (
	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1/funcmock"
)

type options struct {
	mockClient *funcmock.MockUserDefinedFunctionClient
}

// Option is the interface to apply options.
type Option func(*options)

// WithMockGRPCClient creates a new client object with the given mock client for mock testing.
func WithMockGRPCClient(c *funcmock.MockUserDefinedFunctionClient) Option {
	return func(opts *options) {
		opts.mockClient = c
	}
}
