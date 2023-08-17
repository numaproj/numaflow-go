// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: pkg/apis/proto/reduce/v1/reduce.proto

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

// ReduceClient is the client API for Reduce service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReduceClient interface {
	// ReduceFn applies a reduce function to a request stream.
	ReduceFn(ctx context.Context, opts ...grpc.CallOption) (Reduce_ReduceFnClient, error)
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error)
}

type reduceClient struct {
	cc grpc.ClientConnInterface
}

func NewReduceClient(cc grpc.ClientConnInterface) ReduceClient {
	return &reduceClient{cc}
}

func (c *reduceClient) ReduceFn(ctx context.Context, opts ...grpc.CallOption) (Reduce_ReduceFnClient, error) {
	stream, err := c.cc.NewStream(ctx, &Reduce_ServiceDesc.Streams[0], "/reduce.v1.Reduce/ReduceFn", opts...)
	if err != nil {
		return nil, err
	}
	x := &reduceReduceFnClient{stream}
	return x, nil
}

type Reduce_ReduceFnClient interface {
	Send(*ReduceRequest) error
	Recv() (*ReduceResponse, error)
	grpc.ClientStream
}

type reduceReduceFnClient struct {
	grpc.ClientStream
}

func (x *reduceReduceFnClient) Send(m *ReduceRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *reduceReduceFnClient) Recv() (*ReduceResponse, error) {
	m := new(ReduceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *reduceClient) IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error) {
	out := new(ReadyResponse)
	err := c.cc.Invoke(ctx, "/reduce.v1.Reduce/IsReady", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReduceServer is the server API for Reduce service.
// All implementations must embed UnimplementedReduceServer
// for forward compatibility
type ReduceServer interface {
	// ReduceFn applies a reduce function to a request stream.
	ReduceFn(Reduce_ReduceFnServer) error
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error)
	mustEmbedUnimplementedReduceServer()
}

// UnimplementedReduceServer must be embedded to have forward compatible implementations.
type UnimplementedReduceServer struct {
}

func (UnimplementedReduceServer) ReduceFn(Reduce_ReduceFnServer) error {
	return status.Errorf(codes.Unimplemented, "method ReduceFn not implemented")
}
func (UnimplementedReduceServer) IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReady not implemented")
}
func (UnimplementedReduceServer) mustEmbedUnimplementedReduceServer() {}

// UnsafeReduceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReduceServer will
// result in compilation errors.
type UnsafeReduceServer interface {
	mustEmbedUnimplementedReduceServer()
}

func RegisterReduceServer(s grpc.ServiceRegistrar, srv ReduceServer) {
	s.RegisterService(&Reduce_ServiceDesc, srv)
}

func _Reduce_ReduceFn_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ReduceServer).ReduceFn(&reduceReduceFnServer{stream})
}

type Reduce_ReduceFnServer interface {
	Send(*ReduceResponse) error
	Recv() (*ReduceRequest, error)
	grpc.ServerStream
}

type reduceReduceFnServer struct {
	grpc.ServerStream
}

func (x *reduceReduceFnServer) Send(m *ReduceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *reduceReduceFnServer) Recv() (*ReduceRequest, error) {
	m := new(ReduceRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Reduce_IsReady_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReduceServer).IsReady(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/reduce.v1.Reduce/IsReady",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReduceServer).IsReady(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Reduce_ServiceDesc is the grpc.ServiceDesc for Reduce service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Reduce_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "reduce.v1.Reduce",
	HandlerType: (*ReduceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReady",
			Handler:    _Reduce_IsReady_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReduceFn",
			Handler:       _Reduce_ReduceFn_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/apis/proto/reduce/v1/reduce.proto",
}
