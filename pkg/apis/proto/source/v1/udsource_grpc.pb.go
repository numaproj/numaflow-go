// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: pkg/apis/proto/source/v1/udsource.proto

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

// UserDefinedSourceClient is the client API for UserDefinedSource service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserDefinedSourceClient interface {
	// Read returns a stream of datum responses.
	// The size of the returned ReadResponse is less than or equal to the num_records specified in ReadRequest.
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (UserDefinedSource_ReadClient, error)
	// Ack acknowledges a list of datum offsets.
	// It indicates that the datum stream has been processed by the source vertex.
	Ack(ctx context.Context, in *AckRequest, opts ...grpc.CallOption) (*AckResponse, error)
	// Pending returns the number of pending records at the user defined source.
	Pending(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PendingResponse, error)
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error)
}

type userDefinedSourceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserDefinedSourceClient(cc grpc.ClientConnInterface) UserDefinedSourceClient {
	return &userDefinedSourceClient{cc}
}

func (c *userDefinedSourceClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (UserDefinedSource_ReadClient, error) {
	stream, err := c.cc.NewStream(ctx, &UserDefinedSource_ServiceDesc.Streams[0], "/source.v1.UserDefinedSource/Read", opts...)
	if err != nil {
		return nil, err
	}
	x := &userDefinedSourceReadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type UserDefinedSource_ReadClient interface {
	Recv() (*ReadResponse, error)
	grpc.ClientStream
}

type userDefinedSourceReadClient struct {
	grpc.ClientStream
}

func (x *userDefinedSourceReadClient) Recv() (*ReadResponse, error) {
	m := new(ReadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *userDefinedSourceClient) Ack(ctx context.Context, in *AckRequest, opts ...grpc.CallOption) (*AckResponse, error) {
	out := new(AckResponse)
	err := c.cc.Invoke(ctx, "/source.v1.UserDefinedSource/Ack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userDefinedSourceClient) Pending(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PendingResponse, error) {
	out := new(PendingResponse)
	err := c.cc.Invoke(ctx, "/source.v1.UserDefinedSource/Pending", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userDefinedSourceClient) IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error) {
	out := new(ReadyResponse)
	err := c.cc.Invoke(ctx, "/source.v1.UserDefinedSource/IsReady", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserDefinedSourceServer is the server API for UserDefinedSource service.
// All implementations must embed UnimplementedUserDefinedSourceServer
// for forward compatibility
type UserDefinedSourceServer interface {
	// Read returns a stream of datum responses.
	// The size of the returned ReadResponse is less than or equal to the num_records specified in ReadRequest.
	Read(*ReadRequest, UserDefinedSource_ReadServer) error
	// Ack acknowledges a list of datum offsets.
	// It indicates that the datum stream has been processed by the source vertex.
	Ack(context.Context, *AckRequest) (*AckResponse, error)
	// Pending returns the number of pending records at the user defined source.
	Pending(context.Context, *emptypb.Empty) (*PendingResponse, error)
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error)
	mustEmbedUnimplementedUserDefinedSourceServer()
}

// UnimplementedUserDefinedSourceServer must be embedded to have forward compatible implementations.
type UnimplementedUserDefinedSourceServer struct {
}

func (UnimplementedUserDefinedSourceServer) Read(*ReadRequest, UserDefinedSource_ReadServer) error {
	return status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedUserDefinedSourceServer) Ack(context.Context, *AckRequest) (*AckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ack not implemented")
}
func (UnimplementedUserDefinedSourceServer) Pending(context.Context, *emptypb.Empty) (*PendingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Pending not implemented")
}
func (UnimplementedUserDefinedSourceServer) IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReady not implemented")
}
func (UnimplementedUserDefinedSourceServer) mustEmbedUnimplementedUserDefinedSourceServer() {}

// UnsafeUserDefinedSourceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserDefinedSourceServer will
// result in compilation errors.
type UnsafeUserDefinedSourceServer interface {
	mustEmbedUnimplementedUserDefinedSourceServer()
}

func RegisterUserDefinedSourceServer(s grpc.ServiceRegistrar, srv UserDefinedSourceServer) {
	s.RegisterService(&UserDefinedSource_ServiceDesc, srv)
}

func _UserDefinedSource_Read_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ReadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UserDefinedSourceServer).Read(m, &userDefinedSourceReadServer{stream})
}

type UserDefinedSource_ReadServer interface {
	Send(*ReadResponse) error
	grpc.ServerStream
}

type userDefinedSourceReadServer struct {
	grpc.ServerStream
}

func (x *userDefinedSourceReadServer) Send(m *ReadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _UserDefinedSource_Ack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDefinedSourceServer).Ack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/source.v1.UserDefinedSource/Ack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDefinedSourceServer).Ack(ctx, req.(*AckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserDefinedSource_Pending_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDefinedSourceServer).Pending(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/source.v1.UserDefinedSource/Pending",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDefinedSourceServer).Pending(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserDefinedSource_IsReady_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDefinedSourceServer).IsReady(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/source.v1.UserDefinedSource/IsReady",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDefinedSourceServer).IsReady(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// UserDefinedSource_ServiceDesc is the grpc.ServiceDesc for UserDefinedSource service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserDefinedSource_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "source.v1.UserDefinedSource",
	HandlerType: (*UserDefinedSourceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ack",
			Handler:    _UserDefinedSource_Ack_Handler,
		},
		{
			MethodName: "Pending",
			Handler:    _UserDefinedSource_Pending_Handler,
		},
		{
			MethodName: "IsReady",
			Handler:    _UserDefinedSource_IsReady_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Read",
			Handler:       _UserDefinedSource_Read_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/apis/proto/source/v1/udsource.proto",
}
