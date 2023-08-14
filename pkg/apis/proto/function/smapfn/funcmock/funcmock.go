// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/numaproj/numaflow-go/pkg/apis/proto/function/smapfn (interfaces: MapStreamClient)

// Package funcmock is a generated GoMock package.
package funcmock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	smapfn "github.com/numaproj/numaflow-go/pkg/apis/proto/function/smapfn"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockMapStreamClient is a mock of MapStreamClient interface.
type MockMapStreamClient struct {
	ctrl     *gomock.Controller
	recorder *MockMapStreamClientMockRecorder
}

// MockMapStreamClientMockRecorder is the mock recorder for MockMapStreamClient.
type MockMapStreamClientMockRecorder struct {
	mock *MockMapStreamClient
}

// NewMockMapStreamClient creates a new mock instance.
func NewMockMapStreamClient(ctrl *gomock.Controller) *MockMapStreamClient {
	mock := &MockMapStreamClient{ctrl: ctrl}
	mock.recorder = &MockMapStreamClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMapStreamClient) EXPECT() *MockMapStreamClientMockRecorder {
	return m.recorder
}

// IsReady mocks base method.
func (m *MockMapStreamClient) IsReady(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*smapfn.ReadyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IsReady", varargs...)
	ret0, _ := ret[0].(*smapfn.ReadyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsReady indicates an expected call of IsReady.
func (mr *MockMapStreamClientMockRecorder) IsReady(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsReady", reflect.TypeOf((*MockMapStreamClient)(nil).IsReady), varargs...)
}

// MapStreamFn mocks base method.
func (m *MockMapStreamClient) MapStreamFn(arg0 context.Context, arg1 *smapfn.MapStreamRequest, arg2 ...grpc.CallOption) (smapfn.MapStream_MapStreamFnClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "MapStreamFn", varargs...)
	ret0, _ := ret[0].(smapfn.MapStream_MapStreamFnClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MapStreamFn indicates an expected call of MapStreamFn.
func (mr *MockMapStreamClientMockRecorder) MapStreamFn(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MapStreamFn", reflect.TypeOf((*MockMapStreamClient)(nil).MapStreamFn), varargs...)
}
