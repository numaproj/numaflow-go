// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1 (interfaces: SinkClient,Sink_SinkFnClient)
//
// Generated by this command:
//
//	mockgen -destination sinkmock/sinkmock.go -package sinkmock github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1 SinkClient,Sink_SinkFnClient
//

// Package sinkmock is a generated GoMock package.
package sinkmock

import (
	context "context"
	reflect "reflect"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/sink/v1"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockSinkClient is a mock of SinkClient interface.
type MockSinkClient struct {
	ctrl     *gomock.Controller
	recorder *MockSinkClientMockRecorder
	isgomock struct{}
}

// MockSinkClientMockRecorder is the mock recorder for MockSinkClient.
type MockSinkClientMockRecorder struct {
	mock *MockSinkClient
}

// NewMockSinkClient creates a new mock instance.
func NewMockSinkClient(ctrl *gomock.Controller) *MockSinkClient {
	mock := &MockSinkClient{ctrl: ctrl}
	mock.recorder = &MockSinkClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSinkClient) EXPECT() *MockSinkClientMockRecorder {
	return m.recorder
}

// IsReady mocks base method.
func (m *MockSinkClient) IsReady(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1.ReadyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IsReady", varargs...)
	ret0, _ := ret[0].(*v1.ReadyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsReady indicates an expected call of IsReady.
func (mr *MockSinkClientMockRecorder) IsReady(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsReady", reflect.TypeOf((*MockSinkClient)(nil).IsReady), varargs...)
}

// SinkFn mocks base method.
func (m *MockSinkClient) SinkFn(ctx context.Context, opts ...grpc.CallOption) (v1.Sink_SinkFnClient, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SinkFn", varargs...)
	ret0, _ := ret[0].(v1.Sink_SinkFnClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SinkFn indicates an expected call of SinkFn.
func (mr *MockSinkClientMockRecorder) SinkFn(ctx any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SinkFn", reflect.TypeOf((*MockSinkClient)(nil).SinkFn), varargs...)
}

// MockSink_SinkFnClient is a mock of Sink_SinkFnClient interface.
type MockSink_SinkFnClient struct {
	ctrl     *gomock.Controller
	recorder *MockSink_SinkFnClientMockRecorder
	isgomock struct{}
}

// MockSink_SinkFnClientMockRecorder is the mock recorder for MockSink_SinkFnClient.
type MockSink_SinkFnClientMockRecorder struct {
	mock *MockSink_SinkFnClient
}

// NewMockSink_SinkFnClient creates a new mock instance.
func NewMockSink_SinkFnClient(ctrl *gomock.Controller) *MockSink_SinkFnClient {
	mock := &MockSink_SinkFnClient{ctrl: ctrl}
	mock.recorder = &MockSink_SinkFnClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSink_SinkFnClient) EXPECT() *MockSink_SinkFnClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockSink_SinkFnClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockSink_SinkFnClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockSink_SinkFnClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockSink_SinkFnClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockSink_SinkFnClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockSink_SinkFnClient)(nil).Context))
}

// Header mocks base method.
func (m *MockSink_SinkFnClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockSink_SinkFnClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockSink_SinkFnClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockSink_SinkFnClient) Recv() (*v1.SinkResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*v1.SinkResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockSink_SinkFnClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockSink_SinkFnClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m_2 *MockSink_SinkFnClient) RecvMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockSink_SinkFnClientMockRecorder) RecvMsg(m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockSink_SinkFnClient)(nil).RecvMsg), m)
}

// Send mocks base method.
func (m *MockSink_SinkFnClient) Send(arg0 *v1.SinkRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockSink_SinkFnClientMockRecorder) Send(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockSink_SinkFnClient)(nil).Send), arg0)
}

// SendMsg mocks base method.
func (m_2 *MockSink_SinkFnClient) SendMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockSink_SinkFnClientMockRecorder) SendMsg(m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockSink_SinkFnClient)(nil).SendMsg), m)
}

// Trailer mocks base method.
func (m *MockSink_SinkFnClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockSink_SinkFnClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockSink_SinkFnClient)(nil).Trailer))
}
