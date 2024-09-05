// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v4.25.1
// source: pkg/apis/proto/mapstream/v1/mapstream.proto

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
	MapStream_MapStreamFn_FullMethodName = "/mapstream.v1.MapStream/MapStreamFn"
	MapStream_IsReady_FullMethodName     = "/mapstream.v1.MapStream/IsReady"
)

// MapStreamClient is the client API for MapStream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MapStreamClient interface {
	// MapStreamFn applies a function to each request element and returns a stream.
	MapStreamFn(ctx context.Context, in *MapStreamRequest, opts ...grpc.CallOption) (MapStream_MapStreamFnClient, error)
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error)
}

type mapStreamClient struct {
	cc grpc.ClientConnInterface
}

func NewMapStreamClient(cc grpc.ClientConnInterface) MapStreamClient {
	return &mapStreamClient{cc}
}

func (c *mapStreamClient) MapStreamFn(ctx context.Context, in *MapStreamRequest, opts ...grpc.CallOption) (MapStream_MapStreamFnClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &MapStream_ServiceDesc.Streams[0], MapStream_MapStreamFn_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &mapStreamMapStreamFnClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MapStream_MapStreamFnClient interface {
	Recv() (*MapStreamResponse, error)
	grpc.ClientStream
}

type mapStreamMapStreamFnClient struct {
	grpc.ClientStream
}

func (x *mapStreamMapStreamFnClient) Recv() (*MapStreamResponse, error) {
	m := new(MapStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mapStreamClient) IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadyResponse)
	err := c.cc.Invoke(ctx, MapStream_IsReady_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MapStreamServer is the server API for MapStream service.
// All implementations must embed UnimplementedMapStreamServer
// for forward compatibility
type MapStreamServer interface {
	// MapStreamFn applies a function to each request element and returns a stream.
	MapStreamFn(*MapStreamRequest, MapStream_MapStreamFnServer) error
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error)
	mustEmbedUnimplementedMapStreamServer()
}

// UnimplementedMapStreamServer must be embedded to have forward compatible implementations.
type UnimplementedMapStreamServer struct {
}

func (UnimplementedMapStreamServer) MapStreamFn(*MapStreamRequest, MapStream_MapStreamFnServer) error {
	return status.Errorf(codes.Unimplemented, "method MapStreamFn not implemented")
}
func (UnimplementedMapStreamServer) IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReady not implemented")
}
func (UnimplementedMapStreamServer) mustEmbedUnimplementedMapStreamServer() {}

// UnsafeMapStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MapStreamServer will
// result in compilation errors.
type UnsafeMapStreamServer interface {
	mustEmbedUnimplementedMapStreamServer()
}

func RegisterMapStreamServer(s grpc.ServiceRegistrar, srv MapStreamServer) {
	s.RegisterService(&MapStream_ServiceDesc, srv)
}

func _MapStream_MapStreamFn_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(MapStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MapStreamServer).MapStreamFn(m, &mapStreamMapStreamFnServer{ServerStream: stream})
}

type MapStream_MapStreamFnServer interface {
	Send(*MapStreamResponse) error
	grpc.ServerStream
}

type mapStreamMapStreamFnServer struct {
	grpc.ServerStream
}

func (x *mapStreamMapStreamFnServer) Send(m *MapStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _MapStream_IsReady_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MapStreamServer).IsReady(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MapStream_IsReady_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MapStreamServer).IsReady(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// MapStream_ServiceDesc is the grpc.ServiceDesc for MapStream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MapStream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mapstream.v1.MapStream",
	HandlerType: (*MapStreamServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReady",
			Handler:    _MapStream_IsReady_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MapStreamFn",
			Handler:       _MapStream_MapStreamFn_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/apis/proto/mapstream/v1/mapstream.proto",
}
