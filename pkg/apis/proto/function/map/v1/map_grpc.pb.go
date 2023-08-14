// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: pkg/apis/proto/function/map/v1/map.proto

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

// MapClient is the client API for Map service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MapClient interface {
	// MapFn applies a function to each map request element.
	MapFn(ctx context.Context, in *MapRequest, opts ...grpc.CallOption) (*MapResponseList, error)
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error)
}

type mapClient struct {
	cc grpc.ClientConnInterface
}

func NewMapClient(cc grpc.ClientConnInterface) MapClient {
	return &mapClient{cc}
}

func (c *mapClient) MapFn(ctx context.Context, in *MapRequest, opts ...grpc.CallOption) (*MapResponseList, error) {
	out := new(MapResponseList)
	err := c.cc.Invoke(ctx, "/function.map.v1.Map/MapFn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mapClient) IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error) {
	out := new(ReadyResponse)
	err := c.cc.Invoke(ctx, "/function.map.v1.Map/IsReady", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MapServer is the server API for Map service.
// All implementations must embed UnimplementedMapServer
// for forward compatibility
type MapServer interface {
	// MapFn applies a function to each map request element.
	MapFn(context.Context, *MapRequest) (*MapResponseList, error)
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error)
	mustEmbedUnimplementedMapServer()
}

// UnimplementedMapServer must be embedded to have forward compatible implementations.
type UnimplementedMapServer struct {
}

func (UnimplementedMapServer) MapFn(context.Context, *MapRequest) (*MapResponseList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MapFn not implemented")
}
func (UnimplementedMapServer) IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReady not implemented")
}
func (UnimplementedMapServer) mustEmbedUnimplementedMapServer() {}

// UnsafeMapServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MapServer will
// result in compilation errors.
type UnsafeMapServer interface {
	mustEmbedUnimplementedMapServer()
}

func RegisterMapServer(s grpc.ServiceRegistrar, srv MapServer) {
	s.RegisterService(&Map_ServiceDesc, srv)
}

func _Map_MapFn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MapRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MapServer).MapFn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/function.map.v1.Map/MapFn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MapServer).MapFn(ctx, req.(*MapRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Map_IsReady_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MapServer).IsReady(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/function.map.v1.Map/IsReady",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MapServer).IsReady(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Map_ServiceDesc is the grpc.ServiceDesc for Map service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Map_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "function.map.v1.Map",
	HandlerType: (*MapServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MapFn",
			Handler:    _Map_MapFn_Handler,
		},
		{
			MethodName: "IsReady",
			Handler:    _Map_IsReady_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/apis/proto/function/map/v1/map.proto",
}
