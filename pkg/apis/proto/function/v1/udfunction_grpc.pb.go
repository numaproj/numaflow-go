// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: pkg/apis/proto/function/v1/udfunction.proto

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

// UserDefinedFunctionClient is the client API for UserDefinedFunction service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserDefinedFunctionClient interface {
	// MapFn applies a function to each datum element without modifying event time.
	MapFn(ctx context.Context, in *Datum, opts ...grpc.CallOption) (*DatumList, error)
	// MapTFn applies a function to each datum element.
	// In addition to map function, MapTFn also supports assigning a new event time to datum.
	MapTFn(ctx context.Context, in *Datum, opts ...grpc.CallOption) (*DatumList, error)
	// ReduceFn applies a reduce function to a datum stream.
	ReduceFn(ctx context.Context, opts ...grpc.CallOption) (UserDefinedFunction_ReduceFnClient, error)
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error)
}

type userDefinedFunctionClient struct {
	cc grpc.ClientConnInterface
}

func NewUserDefinedFunctionClient(cc grpc.ClientConnInterface) UserDefinedFunctionClient {
	return &userDefinedFunctionClient{cc}
}

func (c *userDefinedFunctionClient) MapFn(ctx context.Context, in *Datum, opts ...grpc.CallOption) (*DatumList, error) {
	out := new(DatumList)
	err := c.cc.Invoke(ctx, "/function.v1.UserDefinedFunction/MapFn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userDefinedFunctionClient) MapTFn(ctx context.Context, in *Datum, opts ...grpc.CallOption) (*DatumList, error) {
	out := new(DatumList)
	err := c.cc.Invoke(ctx, "/function.v1.UserDefinedFunction/MapTFn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userDefinedFunctionClient) ReduceFn(ctx context.Context, opts ...grpc.CallOption) (UserDefinedFunction_ReduceFnClient, error) {
	stream, err := c.cc.NewStream(ctx, &UserDefinedFunction_ServiceDesc.Streams[0], "/function.v1.UserDefinedFunction/ReduceFn", opts...)
	if err != nil {
		return nil, err
	}
	x := &userDefinedFunctionReduceFnClient{stream}
	return x, nil
}

type UserDefinedFunction_ReduceFnClient interface {
	Send(*Datum) error
	CloseAndRecv() (*DatumList, error)
	grpc.ClientStream
}

type userDefinedFunctionReduceFnClient struct {
	grpc.ClientStream
}

func (x *userDefinedFunctionReduceFnClient) Send(m *Datum) error {
	return x.ClientStream.SendMsg(m)
}

func (x *userDefinedFunctionReduceFnClient) CloseAndRecv() (*DatumList, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(DatumList)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *userDefinedFunctionClient) IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ReadyResponse, error) {
	out := new(ReadyResponse)
	err := c.cc.Invoke(ctx, "/function.v1.UserDefinedFunction/IsReady", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserDefinedFunctionServer is the server API for UserDefinedFunction service.
// All implementations must embed UnimplementedUserDefinedFunctionServer
// for forward compatibility
type UserDefinedFunctionServer interface {
	// MapFn applies a function to each datum element without modifying event time.
	MapFn(context.Context, *Datum) (*DatumList, error)
	// MapTFn applies a function to each datum element.
	// In addition to map function, MapTFn also supports assigning a new event time to datum.
	MapTFn(context.Context, *Datum) (*DatumList, error)
	// ReduceFn applies a reduce function to a datum stream.
	ReduceFn(UserDefinedFunction_ReduceFnServer) error
	// IsReady is the heartbeat endpoint for gRPC.
	IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error)
	mustEmbedUnimplementedUserDefinedFunctionServer()
}

// UnimplementedUserDefinedFunctionServer must be embedded to have forward compatible implementations.
type UnimplementedUserDefinedFunctionServer struct {
}

func (UnimplementedUserDefinedFunctionServer) MapFn(context.Context, *Datum) (*DatumList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MapFn not implemented")
}
func (UnimplementedUserDefinedFunctionServer) MapTFn(context.Context, *Datum) (*DatumList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MapTFn not implemented")
}
func (UnimplementedUserDefinedFunctionServer) ReduceFn(UserDefinedFunction_ReduceFnServer) error {
	return status.Errorf(codes.Unimplemented, "method ReduceFn not implemented")
}
func (UnimplementedUserDefinedFunctionServer) IsReady(context.Context, *emptypb.Empty) (*ReadyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReady not implemented")
}
func (UnimplementedUserDefinedFunctionServer) mustEmbedUnimplementedUserDefinedFunctionServer() {}

// UnsafeUserDefinedFunctionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserDefinedFunctionServer will
// result in compilation errors.
type UnsafeUserDefinedFunctionServer interface {
	mustEmbedUnimplementedUserDefinedFunctionServer()
}

func RegisterUserDefinedFunctionServer(s grpc.ServiceRegistrar, srv UserDefinedFunctionServer) {
	s.RegisterService(&UserDefinedFunction_ServiceDesc, srv)
}

func _UserDefinedFunction_MapFn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Datum)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDefinedFunctionServer).MapFn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/function.v1.UserDefinedFunction/MapFn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDefinedFunctionServer).MapFn(ctx, req.(*Datum))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserDefinedFunction_MapTFn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Datum)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDefinedFunctionServer).MapTFn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/function.v1.UserDefinedFunction/MapTFn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDefinedFunctionServer).MapTFn(ctx, req.(*Datum))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserDefinedFunction_ReduceFn_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UserDefinedFunctionServer).ReduceFn(&userDefinedFunctionReduceFnServer{stream})
}

type UserDefinedFunction_ReduceFnServer interface {
	SendAndClose(*DatumList) error
	Recv() (*Datum, error)
	grpc.ServerStream
}

type userDefinedFunctionReduceFnServer struct {
	grpc.ServerStream
}

func (x *userDefinedFunctionReduceFnServer) SendAndClose(m *DatumList) error {
	return x.ServerStream.SendMsg(m)
}

func (x *userDefinedFunctionReduceFnServer) Recv() (*Datum, error) {
	m := new(Datum)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _UserDefinedFunction_IsReady_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDefinedFunctionServer).IsReady(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/function.v1.UserDefinedFunction/IsReady",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDefinedFunctionServer).IsReady(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// UserDefinedFunction_ServiceDesc is the grpc.ServiceDesc for UserDefinedFunction service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserDefinedFunction_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "function.v1.UserDefinedFunction",
	HandlerType: (*UserDefinedFunctionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MapFn",
			Handler:    _UserDefinedFunction_MapFn_Handler,
		},
		{
			MethodName: "MapTFn",
			Handler:    _UserDefinedFunction_MapTFn_Handler,
		},
		{
			MethodName: "IsReady",
			Handler:    _UserDefinedFunction_IsReady_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReduceFn",
			Handler:       _UserDefinedFunction_ReduceFn_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/apis/proto/function/v1/udfunction.proto",
}
