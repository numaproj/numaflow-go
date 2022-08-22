package client

import (
	"github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1/funcmock"
)

type options struct {
	mockClient gRPClientOption
}

type Option interface {
	apply(*options)
}

type gRPClientOption struct {
	isMock   bool
	mockClnt *funcmock.MockUserDefinedFunctionClient
}

func (c gRPClientOption) apply(opts *options) {
	opts.mockClient = gRPClientOption{
		c.isMock,
		c.mockClnt,
	}
}

func WithMockGRPCClient(c *funcmock.MockUserDefinedFunctionClient) Option {
	return gRPClientOption{isMock: true, mockClnt: c}
}
