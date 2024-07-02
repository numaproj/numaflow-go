// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: pkg/apis/proto/batchmap/v1/batchmap.proto

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

// BatchMapClient is the client API for BatchMap service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BatchMapClient interface {
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error)
	// BatchMapFn is a bi-directional streaming rpc which applies a
	// Map function on each BatchMapRequest element of the stream and then returns streams
	// back MapResponse elements.
	BatchMapFn(ctx context.Context, opts ...grpc.CallOption) (BatchMap_BatchMapFnClient, error)
}

type batchMapClient struct {
	cc grpc.ClientConnInterface
}

func NewBatchMapClient(cc grpc.ClientConnInterface) BatchMapClient {
	return &batchMapClient{cc}
}

func (c *batchMapClient) IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error) {
	out := new(ReadyResponse)
	err := c.cc.Invoke(ctx, "/batchmap.v1.BatchMap/IsReady", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *batchMapClient) BatchMapFn(ctx context.Context, opts ...grpc.CallOption) (BatchMap_BatchMapFnClient, error) {
	stream, err := c.cc.NewStream(ctx, &BatchMap_ServiceDesc.Streams[0], "/batchmap.v1.BatchMap/BatchMapFn", opts...)
	if err != nil {
		return nil, err
	}
	x := &batchMapBatchMapFnClient{stream}
	return x, nil
}

type BatchMap_BatchMapFnClient interface {
	Send(*BatchMapRequest) error
	Recv() (*BatchMapResponse, error)
	grpc.ClientStream
}

type batchMapBatchMapFnClient struct {
	grpc.ClientStream
}

func (x *batchMapBatchMapFnClient) Send(m *BatchMapRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *batchMapBatchMapFnClient) Recv() (*BatchMapResponse, error) {
	m := new(BatchMapResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BatchMapServer is the server API for BatchMap service.
// All implementations must embed UnimplementedBatchMapServer
// for forward compatibility
type BatchMapServer interface {
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error)
	// BatchMapFn is a bi-directional streaming rpc which applies a
	// Map function on each BatchMapRequest element of the stream and then returns streams
	// back MapResponse elements.
	BatchMapFn(BatchMap_BatchMapFnServer) error
	mustEmbedUnimplementedBatchMapServer()
}

// UnimplementedBatchMapServer must be embedded to have forward compatible implementations.
type UnimplementedBatchMapServer struct {
}

func (UnimplementedBatchMapServer) IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReady not implemented")
}
func (UnimplementedBatchMapServer) BatchMapFn(BatchMap_BatchMapFnServer) error {
	return status.Errorf(codes.Unimplemented, "method BatchMapFn not implemented")
}
func (UnimplementedBatchMapServer) mustEmbedUnimplementedBatchMapServer() {}

// UnsafeBatchMapServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BatchMapServer will
// result in compilation errors.
type UnsafeBatchMapServer interface {
	mustEmbedUnimplementedBatchMapServer()
}

func RegisterBatchMapServer(s grpc.ServiceRegistrar, srv BatchMapServer) {
	s.RegisterService(&BatchMap_ServiceDesc, srv)
}

func _BatchMap_IsReady_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BatchMapServer).IsReady(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/batchmap.v1.BatchMap/IsReady",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BatchMapServer).IsReady(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BatchMap_BatchMapFn_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BatchMapServer).BatchMapFn(&batchMapBatchMapFnServer{stream})
}

type BatchMap_BatchMapFnServer interface {
	Send(*BatchMapResponse) error
	Recv() (*BatchMapRequest, error)
	grpc.ServerStream
}

type batchMapBatchMapFnServer struct {
	grpc.ServerStream
}

func (x *batchMapBatchMapFnServer) Send(m *BatchMapResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *batchMapBatchMapFnServer) Recv() (*BatchMapRequest, error) {
	m := new(BatchMapRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BatchMap_ServiceDesc is the grpc.ServiceDesc for BatchMap service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BatchMap_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "batchmap.v1.BatchMap",
	HandlerType: (*BatchMapServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReady",
			Handler:    _BatchMap_IsReady_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BatchMapFn",
			Handler:       _BatchMap_BatchMapFn_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/apis/proto/batchmap/v1/batchmap.proto",
}
