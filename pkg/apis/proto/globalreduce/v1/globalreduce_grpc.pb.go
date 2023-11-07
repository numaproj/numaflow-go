// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: pkg/apis/proto/globalreduce/v1/globalreduce.proto

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

// GlobalReduceClient is the client API for GlobalReduce service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GlobalReduceClient interface {
	// GlobalReduceFn applies a reduce function to a request stream.
	GlobalReduceFn(ctx context.Context, opts ...grpc.CallOption) (GlobalReduce_GlobalReduceFnClient, error)
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error)
}

type globalReduceClient struct {
	cc grpc.ClientConnInterface
}

func NewGlobalReduceClient(cc grpc.ClientConnInterface) GlobalReduceClient {
	return &globalReduceClient{cc}
}

func (c *globalReduceClient) GlobalReduceFn(ctx context.Context, opts ...grpc.CallOption) (GlobalReduce_GlobalReduceFnClient, error) {
	stream, err := c.cc.NewStream(ctx, &GlobalReduce_ServiceDesc.Streams[0], "/globalreduce.v1.GlobalReduce/GlobalReduceFn", opts...)
	if err != nil {
		return nil, err
	}
	x := &globalReduceGlobalReduceFnClient{stream}
	return x, nil
}

type GlobalReduce_GlobalReduceFnClient interface {
	Send(*GlobalReduceRequest) error
	Recv() (*GlobalReduceResponse, error)
	grpc.ClientStream
}

type globalReduceGlobalReduceFnClient struct {
	grpc.ClientStream
}

func (x *globalReduceGlobalReduceFnClient) Send(m *GlobalReduceRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *globalReduceGlobalReduceFnClient) Recv() (*GlobalReduceResponse, error) {
	m := new(GlobalReduceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *globalReduceClient) IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error) {
	out := new(ReadyResponse)
	err := c.cc.Invoke(ctx, "/globalreduce.v1.GlobalReduce/IsReady", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GlobalReduceServer is the server API for GlobalReduce service.
// All implementations must embed UnimplementedGlobalReduceServer
// for forward compatibility
type GlobalReduceServer interface {
	// GlobalReduceFn applies a reduce function to a request stream.
	GlobalReduceFn(GlobalReduce_GlobalReduceFnServer) error
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error)
	mustEmbedUnimplementedGlobalReduceServer()
}

// UnimplementedGlobalReduceServer must be embedded to have forward compatible implementations.
type UnimplementedGlobalReduceServer struct {
}

func (UnimplementedGlobalReduceServer) GlobalReduceFn(GlobalReduce_GlobalReduceFnServer) error {
	return status.Errorf(codes.Unimplemented, "method GlobalReduceFn not implemented")
}
func (UnimplementedGlobalReduceServer) IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReady not implemented")
}
func (UnimplementedGlobalReduceServer) mustEmbedUnimplementedGlobalReduceServer() {}

// UnsafeGlobalReduceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GlobalReduceServer will
// result in compilation errors.
type UnsafeGlobalReduceServer interface {
	mustEmbedUnimplementedGlobalReduceServer()
}

func RegisterGlobalReduceServer(s grpc.ServiceRegistrar, srv GlobalReduceServer) {
	s.RegisterService(&GlobalReduce_ServiceDesc, srv)
}

func _GlobalReduce_GlobalReduceFn_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GlobalReduceServer).GlobalReduceFn(&globalReduceGlobalReduceFnServer{stream})
}

type GlobalReduce_GlobalReduceFnServer interface {
	Send(*GlobalReduceResponse) error
	Recv() (*GlobalReduceRequest, error)
	grpc.ServerStream
}

type globalReduceGlobalReduceFnServer struct {
	grpc.ServerStream
}

func (x *globalReduceGlobalReduceFnServer) Send(m *GlobalReduceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *globalReduceGlobalReduceFnServer) Recv() (*GlobalReduceRequest, error) {
	m := new(GlobalReduceRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GlobalReduce_IsReady_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GlobalReduceServer).IsReady(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/globalreduce.v1.GlobalReduce/IsReady",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GlobalReduceServer).IsReady(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// GlobalReduce_ServiceDesc is the grpc.ServiceDesc for GlobalReduce service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GlobalReduce_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "globalreduce.v1.GlobalReduce",
	HandlerType: (*GlobalReduceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReady",
			Handler:    _GlobalReduce_IsReady_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GlobalReduceFn",
			Handler:       _GlobalReduce_GlobalReduceFn_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/apis/proto/globalreduce/v1/globalreduce.proto",
}
