// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.29.1
// source: pkg/apis/proto/map/v1/map.proto

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
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Map_MapFn_FullMethodName   = "/map.v1.Map/MapFn"
	Map_IsReady_FullMethodName = "/map.v1.Map/IsReady"
)

// MapClient is the client API for Map service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MapClient interface {
	// MapFn applies a function to each map request element.
	MapFn(ctx context.Context, opts ...grpc.CallOption) (Map_MapFnClient, error)
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error)
}

type mapClient struct {
	cc grpc.ClientConnInterface
}

func NewMapClient(cc grpc.ClientConnInterface) MapClient {
	return &mapClient{cc}
}

func (c *mapClient) MapFn(ctx context.Context, opts ...grpc.CallOption) (Map_MapFnClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Map_ServiceDesc.Streams[0], Map_MapFn_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &mapMapFnClient{ClientStream: stream}
	return x, nil
}

type Map_MapFnClient interface {
	Send(*MapRequest) error
	Recv() (*MapResponse, error)
	grpc.ClientStream
}

type mapMapFnClient struct {
	grpc.ClientStream
}

func (x *mapMapFnClient) Send(m *MapRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *mapMapFnClient) Recv() (*MapResponse, error) {
	m := new(MapResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mapClient) IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadyResponse)
	err := c.cc.Invoke(ctx, Map_IsReady_FullMethodName, in, out, cOpts...)
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
	MapFn(Map_MapFnServer) error
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error)
	mustEmbedUnimplementedMapServer()
}

// UnimplementedMapServer must be embedded to have forward compatible implementations.
type UnimplementedMapServer struct {
}

func (UnimplementedMapServer) MapFn(Map_MapFnServer) error {
	return status.Errorf(codes.Unimplemented, "method MapFn not implemented")
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

func _Map_MapFn_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MapServer).MapFn(&mapMapFnServer{ServerStream: stream})
}

type Map_MapFnServer interface {
	Send(*MapResponse) error
	Recv() (*MapRequest, error)
	grpc.ServerStream
}

type mapMapFnServer struct {
	grpc.ServerStream
}

func (x *mapMapFnServer) Send(m *MapResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *mapMapFnServer) Recv() (*MapRequest, error) {
	m := new(MapRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
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
		FullMethod: Map_IsReady_FullMethodName,
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
	ServiceName: "map.v1.Map",
	HandlerType: (*MapServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReady",
			Handler:    _Map_IsReady_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MapFn",
			Handler:       _Map_MapFn_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/apis/proto/map/v1/map.proto",
}
