// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/numaproj/numaflow-go/pkg/apis/proto/function/reducefn (interfaces: ReduceClient)

// Package funcmock is a generated GoMock package.
package funcmock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	reducefn "github.com/numaproj/numaflow-go/pkg/apis/proto/function/reducefn"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockReduceClient is a mock of ReduceClient interface.
type MockReduceClient struct {
	ctrl     *gomock.Controller
	recorder *MockReduceClientMockRecorder
}

// MockReduceClientMockRecorder is the mock recorder for MockReduceClient.
type MockReduceClientMockRecorder struct {
	mock *MockReduceClient
}

// NewMockReduceClient creates a new mock instance.
func NewMockReduceClient(ctrl *gomock.Controller) *MockReduceClient {
	mock := &MockReduceClient{ctrl: ctrl}
	mock.recorder = &MockReduceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReduceClient) EXPECT() *MockReduceClientMockRecorder {
	return m.recorder
}

// IsReady mocks base method.
func (m *MockReduceClient) IsReady(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*reducefn.ReadyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IsReady", varargs...)
	ret0, _ := ret[0].(*reducefn.ReadyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsReady indicates an expected call of IsReady.
func (mr *MockReduceClientMockRecorder) IsReady(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsReady", reflect.TypeOf((*MockReduceClient)(nil).IsReady), varargs...)
}

// ReduceFn mocks base method.
func (m *MockReduceClient) ReduceFn(arg0 context.Context, arg1 ...grpc.CallOption) (reducefn.Reduce_ReduceFnClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReduceFn", varargs...)
	ret0, _ := ret[0].(reducefn.Reduce_ReduceFnClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReduceFn indicates an expected call of ReduceFn.
func (mr *MockReduceClientMockRecorder) ReduceFn(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReduceFn", reflect.TypeOf((*MockReduceClient)(nil).ReduceFn), varargs...)
}
