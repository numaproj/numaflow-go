// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v2 (interfaces: UserDefinedSinkClient,UserDefinedSink_SinkFnClient)

// Package sinkmock is a generated GoMock package.
package sinkmock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v2 "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockUserDefinedSinkClient is a mock of UserDefinedSinkClient interface.
type MockUserDefinedSinkClient struct {
	ctrl     *gomock.Controller
	recorder *MockUserDefinedSinkClientMockRecorder
}

// MockUserDefinedSinkClientMockRecorder is the mock recorder for MockUserDefinedSinkClient.
type MockUserDefinedSinkClientMockRecorder struct {
	mock *MockUserDefinedSinkClient
}

// NewMockUserDefinedSinkClient creates a new mock instance.
func NewMockUserDefinedSinkClient(ctrl *gomock.Controller) *MockUserDefinedSinkClient {
	mock := &MockUserDefinedSinkClient{ctrl: ctrl}
	mock.recorder = &MockUserDefinedSinkClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserDefinedSinkClient) EXPECT() *MockUserDefinedSinkClientMockRecorder {
	return m.recorder
}

// IsReady mocks base method.
func (m *MockUserDefinedSinkClient) IsReady(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*v2.ReadyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IsReady", varargs...)
	ret0, _ := ret[0].(*v2.ReadyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsReady indicates an expected call of IsReady.
func (mr *MockUserDefinedSinkClientMockRecorder) IsReady(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsReady", reflect.TypeOf((*MockUserDefinedSinkClient)(nil).IsReady), varargs...)
}

// SinkFn mocks base method.
func (m *MockUserDefinedSinkClient) SinkFn(arg0 context.Context, arg1 ...grpc.CallOption) (v2.UserDefinedSink_SinkFnClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SinkFn", varargs...)
	ret0, _ := ret[0].(v2.UserDefinedSink_SinkFnClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SinkFn indicates an expected call of SinkFn.
func (mr *MockUserDefinedSinkClientMockRecorder) SinkFn(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SinkFn", reflect.TypeOf((*MockUserDefinedSinkClient)(nil).SinkFn), varargs...)
}

// MockUserDefinedSink_SinkFnClient is a mock of UserDefinedSink_SinkFnClient interface.
type MockUserDefinedSink_SinkFnClient struct {
	ctrl     *gomock.Controller
	recorder *MockUserDefinedSink_SinkFnClientMockRecorder
}

// MockUserDefinedSink_SinkFnClientMockRecorder is the mock recorder for MockUserDefinedSink_SinkFnClient.
type MockUserDefinedSink_SinkFnClientMockRecorder struct {
	mock *MockUserDefinedSink_SinkFnClient
}

// NewMockUserDefinedSink_SinkFnClient creates a new mock instance.
func NewMockUserDefinedSink_SinkFnClient(ctrl *gomock.Controller) *MockUserDefinedSink_SinkFnClient {
	mock := &MockUserDefinedSink_SinkFnClient{ctrl: ctrl}
	mock.recorder = &MockUserDefinedSink_SinkFnClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserDefinedSink_SinkFnClient) EXPECT() *MockUserDefinedSink_SinkFnClientMockRecorder {
	return m.recorder
}

// CloseAndRecv mocks base method.
func (m *MockUserDefinedSink_SinkFnClient) CloseAndRecv() (*v2.ResponseList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseAndRecv")
	ret0, _ := ret[0].(*v2.ResponseList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloseAndRecv indicates an expected call of CloseAndRecv.
func (mr *MockUserDefinedSink_SinkFnClientMockRecorder) CloseAndRecv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseAndRecv", reflect.TypeOf((*MockUserDefinedSink_SinkFnClient)(nil).CloseAndRecv))
}

// CloseSend mocks base method.
func (m *MockUserDefinedSink_SinkFnClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockUserDefinedSink_SinkFnClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockUserDefinedSink_SinkFnClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockUserDefinedSink_SinkFnClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockUserDefinedSink_SinkFnClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockUserDefinedSink_SinkFnClient)(nil).Context))
}

// Header mocks base method.
func (m *MockUserDefinedSink_SinkFnClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockUserDefinedSink_SinkFnClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockUserDefinedSink_SinkFnClient)(nil).Header))
}

// RecvMsg mocks base method.
func (m *MockUserDefinedSink_SinkFnClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockUserDefinedSink_SinkFnClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockUserDefinedSink_SinkFnClient)(nil).RecvMsg), arg0)
}

// Send mocks base method.
func (m *MockUserDefinedSink_SinkFnClient) Send(arg0 *v2.Datum) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockUserDefinedSink_SinkFnClientMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockUserDefinedSink_SinkFnClient)(nil).Send), arg0)
}

// SendMsg mocks base method.
func (m *MockUserDefinedSink_SinkFnClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockUserDefinedSink_SinkFnClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockUserDefinedSink_SinkFnClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method.
func (m *MockUserDefinedSink_SinkFnClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockUserDefinedSink_SinkFnClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockUserDefinedSink_SinkFnClient)(nil).Trailer))
}
