// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v4.25.1
// source: pkg/apis/proto/sink/v1/sink.proto

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
	Sink_SinkFn_FullMethodName  = "/sink.v1.Sink/SinkFn"
	Sink_IsReady_FullMethodName = "/sink.v1.Sink/IsReady"
)

// SinkClient is the client API for Sink service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SinkClient interface {
	// SinkFn writes the request to a user defined sink.
	SinkFn(ctx context.Context, opts ...grpc.CallOption) (Sink_SinkFnClient, error)
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error)
}

type sinkClient struct {
	cc grpc.ClientConnInterface
}

func NewSinkClient(cc grpc.ClientConnInterface) SinkClient {
	return &sinkClient{cc}
}

func (c *sinkClient) SinkFn(ctx context.Context, opts ...grpc.CallOption) (Sink_SinkFnClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Sink_ServiceDesc.Streams[0], Sink_SinkFn_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &sinkSinkFnClient{ClientStream: stream}
	return x, nil
}

type Sink_SinkFnClient interface {
	Send(*SinkRequest) error
	CloseAndRecv() (*SinkResponse, error)
	grpc.ClientStream
}

type sinkSinkFnClient struct {
	grpc.ClientStream
}

func (x *sinkSinkFnClient) Send(m *SinkRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *sinkSinkFnClient) CloseAndRecv() (*SinkResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(SinkResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sinkClient) IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadyResponse)
	err := c.cc.Invoke(ctx, Sink_IsReady_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SinkServer is the server API for Sink service.
// All implementations must embed UnimplementedSinkServer
// for forward compatibility
type SinkServer interface {
	// SinkFn writes the request to a user defined sink.
	SinkFn(Sink_SinkFnServer) error
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error)
	mustEmbedUnimplementedSinkServer()
}

// UnimplementedSinkServer must be embedded to have forward compatible implementations.
type UnimplementedSinkServer struct {
}

func (UnimplementedSinkServer) SinkFn(Sink_SinkFnServer) error {
	return status.Errorf(codes.Unimplemented, "method SinkFn not implemented")
}
func (UnimplementedSinkServer) IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReady not implemented")
}
func (UnimplementedSinkServer) mustEmbedUnimplementedSinkServer() {}

// UnsafeSinkServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SinkServer will
// result in compilation errors.
type UnsafeSinkServer interface {
	mustEmbedUnimplementedSinkServer()
}

func RegisterSinkServer(s grpc.ServiceRegistrar, srv SinkServer) {
	s.RegisterService(&Sink_ServiceDesc, srv)
}

func _Sink_SinkFn_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SinkServer).SinkFn(&sinkSinkFnServer{ServerStream: stream})
}

type Sink_SinkFnServer interface {
	SendAndClose(*SinkResponse) error
	Recv() (*SinkRequest, error)
	grpc.ServerStream
}

type sinkSinkFnServer struct {
	grpc.ServerStream
}

func (x *sinkSinkFnServer) SendAndClose(m *SinkResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *sinkSinkFnServer) Recv() (*SinkRequest, error) {
	m := new(SinkRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Sink_IsReady_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SinkServer).IsReady(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sink_IsReady_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SinkServer).IsReady(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Sink_ServiceDesc is the grpc.ServiceDesc for Sink service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sink_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sink.v1.Sink",
	HandlerType: (*SinkServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReady",
			Handler:    _Sink_IsReady_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SinkFn",
			Handler:       _Sink_SinkFn_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/apis/proto/sink/v1/sink.proto",
}
