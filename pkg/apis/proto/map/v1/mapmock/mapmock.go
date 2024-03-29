// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1 (interfaces: MapClient)

// Package mapmock is a generated GoMock package.
package mapmock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/map/v1"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockMapClient is a mock of MapClient interface.
type MockMapClient struct {
	ctrl     *gomock.Controller
	recorder *MockMapClientMockRecorder
}

// MockMapClientMockRecorder is the mock recorder for MockMapClient.
type MockMapClientMockRecorder struct {
	mock *MockMapClient
}

// NewMockMapClient creates a new mock instance.
func NewMockMapClient(ctrl *gomock.Controller) *MockMapClient {
	mock := &MockMapClient{ctrl: ctrl}
	mock.recorder = &MockMapClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMapClient) EXPECT() *MockMapClientMockRecorder {
	return m.recorder
}

// IsReady mocks base method.
func (m *MockMapClient) IsReady(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*v1.ReadyResponse, error) {
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
func (mr *MockMapClientMockRecorder) IsReady(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsReady", reflect.TypeOf((*MockMapClient)(nil).IsReady), varargs...)
}

// MapFn mocks base method.
func (m *MockMapClient) MapFn(arg0 context.Context, arg1 *v1.MapRequest, arg2 ...grpc.CallOption) (*v1.MapResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "MapFn", varargs...)
	ret0, _ := ret[0].(*v1.MapResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MapFn indicates an expected call of MapFn.
func (mr *MockMapClientMockRecorder) MapFn(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MapFn", reflect.TypeOf((*MockMapClient)(nil).MapFn), varargs...)
}
