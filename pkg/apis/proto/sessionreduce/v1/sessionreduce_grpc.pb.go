// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: pkg/apis/proto/sessionreduce/v1/sessionreduce.proto

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

// SessionReduceClient is the client API for SessionReduce service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SessionReduceClient interface {
	// SessionReduceFn applies a reduce function to a request stream.
	SessionReduceFn(ctx context.Context, opts ...grpc.CallOption) (SessionReduce_SessionReduceFnClient, error)
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error)
}

type sessionReduceClient struct {
	cc grpc.ClientConnInterface
}

func NewSessionReduceClient(cc grpc.ClientConnInterface) SessionReduceClient {
	return &sessionReduceClient{cc}
}

func (c *sessionReduceClient) SessionReduceFn(ctx context.Context, opts ...grpc.CallOption) (SessionReduce_SessionReduceFnClient, error) {
	stream, err := c.cc.NewStream(ctx, &SessionReduce_ServiceDesc.Streams[0], "/sessionreduce.v1.SessionReduce/SessionReduceFn", opts...)
	if err != nil {
		return nil, err
	}
	x := &sessionReduceSessionReduceFnClient{stream}
	return x, nil
}

type SessionReduce_SessionReduceFnClient interface {
	Send(*SessionReduceRequest) error
	Recv() (*SessionReduceResponse, error)
	grpc.ClientStream
}

type sessionReduceSessionReduceFnClient struct {
	grpc.ClientStream
}

func (x *sessionReduceSessionReduceFnClient) Send(m *SessionReduceRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *sessionReduceSessionReduceFnClient) Recv() (*SessionReduceResponse, error) {
	m := new(SessionReduceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sessionReduceClient) IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error) {
	out := new(ReadyResponse)
	err := c.cc.Invoke(ctx, "/sessionreduce.v1.SessionReduce/IsReady", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SessionReduceServer is the server API for SessionReduce service.
// All implementations must embed UnimplementedSessionReduceServer
// for forward compatibility
type SessionReduceServer interface {
	// SessionReduceFn applies a reduce function to a request stream.
	SessionReduceFn(SessionReduce_SessionReduceFnServer) error
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error)
	mustEmbedUnimplementedSessionReduceServer()
}

// UnimplementedSessionReduceServer must be embedded to have forward compatible implementations.
type UnimplementedSessionReduceServer struct {
}

func (UnimplementedSessionReduceServer) SessionReduceFn(SessionReduce_SessionReduceFnServer) error {
	return status.Errorf(codes.Unimplemented, "method SessionReduceFn not implemented")
}
func (UnimplementedSessionReduceServer) IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReady not implemented")
}
func (UnimplementedSessionReduceServer) mustEmbedUnimplementedSessionReduceServer() {}

// UnsafeSessionReduceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SessionReduceServer will
// result in compilation errors.
type UnsafeSessionReduceServer interface {
	mustEmbedUnimplementedSessionReduceServer()
}

func RegisterSessionReduceServer(s grpc.ServiceRegistrar, srv SessionReduceServer) {
	s.RegisterService(&SessionReduce_ServiceDesc, srv)
}

func _SessionReduce_SessionReduceFn_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SessionReduceServer).SessionReduceFn(&sessionReduceSessionReduceFnServer{stream})
}

type SessionReduce_SessionReduceFnServer interface {
	Send(*SessionReduceResponse) error
	Recv() (*SessionReduceRequest, error)
	grpc.ServerStream
}

type sessionReduceSessionReduceFnServer struct {
	grpc.ServerStream
}

func (x *sessionReduceSessionReduceFnServer) Send(m *SessionReduceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *sessionReduceSessionReduceFnServer) Recv() (*SessionReduceRequest, error) {
	m := new(SessionReduceRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SessionReduce_IsReady_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionReduceServer).IsReady(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sessionreduce.v1.SessionReduce/IsReady",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionReduceServer).IsReady(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// SessionReduce_ServiceDesc is the grpc.ServiceDesc for SessionReduce service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SessionReduce_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sessionreduce.v1.SessionReduce",
	HandlerType: (*SessionReduceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReady",
			Handler:    _SessionReduce_IsReady_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SessionReduceFn",
			Handler:       _SessionReduce_SessionReduceFn_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/apis/proto/sessionreduce/v1/sessionreduce.proto",
}
