// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: pkg/apis/proto/sourcetransform/v1/transform.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SourceTransformClient is the client API for SourceTransform service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SourceTransformClient interface {
	// SourceTransformFn applies a function to each request element.
	// In addition to map function, SourceTransformFn also supports assigning a new event time to response.
	// SourceTransformFn can be used only at source vertex by source data transformer.
	SourceTransformFn(ctx context.Context, in *SourceTransformRequest, opts ...grpc.CallOption) (*SourceTransformResponse, error)
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error)
}

type sourceTransformClient struct {
	cc grpc.ClientConnInterface
}

func NewSourceTransformClient(cc grpc.ClientConnInterface) SourceTransformClient {
	return &sourceTransformClient{cc}
}

func (c *sourceTransformClient) SourceTransformFn(ctx context.Context, in *SourceTransformRequest, opts ...grpc.CallOption) (*SourceTransformResponse, error) {
	out := new(SourceTransformResponse)
	err := c.cc.Invoke(ctx, "/sourcetransformer.v1.SourceTransform/SourceTransformFn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sourceTransformClient) IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error) {
	out := new(ReadyResponse)
	err := c.cc.Invoke(ctx, "/sourcetransformer.v1.SourceTransform/IsReady", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SourceTransformServer is the server API for SourceTransform service.
// All implementations must embed UnimplementedSourceTransformServer
// for forward compatibility
type SourceTransformServer interface {
	// SourceTransformFn applies a function to each request element.
	// In addition to map function, SourceTransformFn also supports assigning a new event time to response.
	// SourceTransformFn can be used only at source vertex by source data transformer.
	SourceTransformFn(context.Context, *SourceTransformRequest) (*SourceTransformResponse, error)
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error)
	mustEmbedUnimplementedSourceTransformServer()
}

// UnimplementedSourceTransformServer must be embedded to have forward compatible implementations.
type UnimplementedSourceTransformServer struct {
}

func (UnimplementedSourceTransformServer) SourceTransformFn(context.Context, *SourceTransformRequest) (*SourceTransformResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SourceTransformFn not implemented")
}
func (UnimplementedSourceTransformServer) IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReady not implemented")
}
func (UnimplementedSourceTransformServer) mustEmbedUnimplementedSourceTransformServer() {}

// UnsafeSourceTransformServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SourceTransformServer will
// result in compilation errors.
type UnsafeSourceTransformServer interface {
	mustEmbedUnimplementedSourceTransformServer()
}

func RegisterSourceTransformServer(s grpc.ServiceRegistrar, srv SourceTransformServer) {
	s.RegisterService(&SourceTransform_ServiceDesc, srv)
}

func _SourceTransform_SourceTransformFn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SourceTransformRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SourceTransformServer).SourceTransformFn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sourcetransformer.v1.SourceTransform/SourceTransformFn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SourceTransformServer).SourceTransformFn(ctx, req.(*SourceTransformRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SourceTransform_IsReady_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SourceTransformServer).IsReady(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sourcetransformer.v1.SourceTransform/IsReady",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SourceTransformServer).IsReady(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// SourceTransform_ServiceDesc is the grpc.ServiceDesc for SourceTransform service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SourceTransform_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sourcetransformer.v1.SourceTransform",
	HandlerType: (*SourceTransformServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SourceTransformFn",
			Handler:    _SourceTransform_SourceTransformFn_Handler,
		},
		{
			MethodName: "IsReady",
			Handler:    _SourceTransform_IsReady_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/apis/proto/sourcetransform/v1/transform.proto",
}
