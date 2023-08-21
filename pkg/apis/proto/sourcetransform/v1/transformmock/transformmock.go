// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/numaproj/numaflow-go/pkg/apis/proto/sourcetransform/v1 (interfaces: SourceTransformClient)

// Package transformermock is a generated GoMock package.
package transformermock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/sourcetransform/v1"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockSourceTransformClient is a mock of SourceTransformClient interface.
type MockSourceTransformClient struct {
	ctrl     *gomock.Controller
	recorder *MockSourceTransformClientMockRecorder
}

// MockSourceTransformClientMockRecorder is the mock recorder for MockSourceTransformClient.
type MockSourceTransformClientMockRecorder struct {
	mock *MockSourceTransformClient
}

// NewMockSourceTransformClient creates a new mock instance.
func NewMockSourceTransformClient(ctrl *gomock.Controller) *MockSourceTransformClient {
	mock := &MockSourceTransformClient{ctrl: ctrl}
	mock.recorder = &MockSourceTransformClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSourceTransformClient) EXPECT() *MockSourceTransformClientMockRecorder {
	return m.recorder
}

// IsReady mocks base method.
func (m *MockSourceTransformClient) IsReady(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*v1.ReadyResponse, error) {
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
func (mr *MockSourceTransformClientMockRecorder) IsReady(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsReady", reflect.TypeOf((*MockSourceTransformClient)(nil).IsReady), varargs...)
}

// SourceTransformFn mocks base method.
func (m *MockSourceTransformClient) SourceTransformFn(arg0 context.Context, arg1 *v1.SourceTransformRequest, arg2 ...grpc.CallOption) (*v1.SourceTransformResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SourceTransformFn", varargs...)
	ret0, _ := ret[0].(*v1.SourceTransformResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SourceTransformFn indicates an expected call of SourceTransformFn.
func (mr *MockSourceTransformClientMockRecorder) SourceTransformFn(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SourceTransformFn", reflect.TypeOf((*MockSourceTransformClient)(nil).SourceTransformFn), varargs...)
}
