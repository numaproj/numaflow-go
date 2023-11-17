// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/numaproj/numaflow-go/pkg/apis/proto/globalreduce/v1 (interfaces: GlobalReduceClient,GlobalReduce_GlobalReduceFnClient)

// Package globalreducemock is a generated GoMock package.
package globalreducemock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/globalreduce/v1"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockGlobalReduceClient is a mock of GlobalReduceClient interface.
type MockGlobalReduceClient struct {
	ctrl     *gomock.Controller
	recorder *MockGlobalReduceClientMockRecorder
}

// MockGlobalReduceClientMockRecorder is the mock recorder for MockGlobalReduceClient.
type MockGlobalReduceClientMockRecorder struct {
	mock *MockGlobalReduceClient
}

// NewMockGlobalReduceClient creates a new mock instance.
func NewMockGlobalReduceClient(ctrl *gomock.Controller) *MockGlobalReduceClient {
	mock := &MockGlobalReduceClient{ctrl: ctrl}
	mock.recorder = &MockGlobalReduceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGlobalReduceClient) EXPECT() *MockGlobalReduceClientMockRecorder {
	return m.recorder
}

// GlobalReduceFn mocks base method.
func (m *MockGlobalReduceClient) GlobalReduceFn(arg0 context.Context, arg1 ...grpc.CallOption) (v1.GlobalReduce_GlobalReduceFnClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GlobalReduceFn", varargs...)
	ret0, _ := ret[0].(v1.GlobalReduce_GlobalReduceFnClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GlobalReduceFn indicates an expected call of GlobalReduceFn.
func (mr *MockGlobalReduceClientMockRecorder) GlobalReduceFn(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GlobalReduceFn", reflect.TypeOf((*MockGlobalReduceClient)(nil).GlobalReduceFn), varargs...)
}

// IsReady mocks base method.
func (m *MockGlobalReduceClient) IsReady(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*v1.ReadyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IsReady", varargs...)
	ret0, _ := ret[0].(*v1.ReadyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsReady indicates an expected call of IsReady.
func (mr *MockGlobalReduceClientMockRecorder) IsReady(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsReady", reflect.TypeOf((*MockGlobalReduceClient)(nil).IsReady), varargs...)
}

// MockGlobalReduce_GlobalReduceFnClient is a mock of GlobalReduce_GlobalReduceFnClient interface.
type MockGlobalReduce_GlobalReduceFnClient struct {
	ctrl     *gomock.Controller
	recorder *MockGlobalReduce_GlobalReduceFnClientMockRecorder
}

// MockGlobalReduce_GlobalReduceFnClientMockRecorder is the mock recorder for MockGlobalReduce_GlobalReduceFnClient.
type MockGlobalReduce_GlobalReduceFnClientMockRecorder struct {
	mock *MockGlobalReduce_GlobalReduceFnClient
}

// NewMockGlobalReduce_GlobalReduceFnClient creates a new mock instance.
func NewMockGlobalReduce_GlobalReduceFnClient(ctrl *gomock.Controller) *MockGlobalReduce_GlobalReduceFnClient {
	mock := &MockGlobalReduce_GlobalReduceFnClient{ctrl: ctrl}
	mock.recorder = &MockGlobalReduce_GlobalReduceFnClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGlobalReduce_GlobalReduceFnClient) EXPECT() *MockGlobalReduce_GlobalReduceFnClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockGlobalReduce_GlobalReduceFnClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockGlobalReduce_GlobalReduceFnClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockGlobalReduce_GlobalReduceFnClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockGlobalReduce_GlobalReduceFnClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockGlobalReduce_GlobalReduceFnClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockGlobalReduce_GlobalReduceFnClient)(nil).Context))
}

// Header mocks base method.
func (m *MockGlobalReduce_GlobalReduceFnClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockGlobalReduce_GlobalReduceFnClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockGlobalReduce_GlobalReduceFnClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockGlobalReduce_GlobalReduceFnClient) Recv() (*v1.GlobalReduceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*v1.GlobalReduceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockGlobalReduce_GlobalReduceFnClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockGlobalReduce_GlobalReduceFnClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m *MockGlobalReduce_GlobalReduceFnClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockGlobalReduce_GlobalReduceFnClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockGlobalReduce_GlobalReduceFnClient)(nil).RecvMsg), arg0)
}

// Send mocks base method.
func (m *MockGlobalReduce_GlobalReduceFnClient) Send(arg0 *v1.GlobalReduceRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockGlobalReduce_GlobalReduceFnClientMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockGlobalReduce_GlobalReduceFnClient)(nil).Send), arg0)
}

// SendMsg mocks base method.
func (m *MockGlobalReduce_GlobalReduceFnClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockGlobalReduce_GlobalReduceFnClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockGlobalReduce_GlobalReduceFnClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method.
func (m *MockGlobalReduce_GlobalReduceFnClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockGlobalReduce_GlobalReduceFnClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockGlobalReduce_GlobalReduceFnClient)(nil).Trailer))
}