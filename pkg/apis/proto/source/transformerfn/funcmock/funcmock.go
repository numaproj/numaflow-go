// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/numaproj/numaflow-go/pkg/apis/proto/source/transformerfn (interfaces: SourceTransformerClient)

// Package funcmock is a generated GoMock package.
package funcmock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	transformerfn "github.com/numaproj/numaflow-go/pkg/apis/proto/source/transformerfn"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockSourceTransformerClient is a mock of SourceTransformerClient interface.
type MockSourceTransformerClient struct {
	ctrl     *gomock.Controller
	recorder *MockSourceTransformerClientMockRecorder
}

// MockSourceTransformerClientMockRecorder is the mock recorder for MockSourceTransformerClient.
type MockSourceTransformerClientMockRecorder struct {
	mock *MockSourceTransformerClient
}

// NewMockSourceTransformerClient creates a new mock instance.
func NewMockSourceTransformerClient(ctrl *gomock.Controller) *MockSourceTransformerClient {
	mock := &MockSourceTransformerClient{ctrl: ctrl}
	mock.recorder = &MockSourceTransformerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSourceTransformerClient) EXPECT() *MockSourceTransformerClientMockRecorder {
	return m.recorder
}

// IsReady mocks base method.
func (m *MockSourceTransformerClient) IsReady(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*transformerfn.ReadyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IsReady", varargs...)
	ret0, _ := ret[0].(*transformerfn.ReadyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsReady indicates an expected call of IsReady.
func (mr *MockSourceTransformerClientMockRecorder) IsReady(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsReady", reflect.TypeOf((*MockSourceTransformerClient)(nil).IsReady), varargs...)
}

// SourceTransformer mocks base method.
func (m *MockSourceTransformerClient) SourceTransformer(arg0 context.Context, arg1 *transformerfn.SourceTransformerRequest, arg2 ...grpc.CallOption) (*transformerfn.SourceTransformerResponseList, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SourceTransformer", varargs...)
	ret0, _ := ret[0].(*transformerfn.SourceTransformerResponseList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SourceTransformer indicates an expected call of SourceTransformer.
func (mr *MockSourceTransformerClientMockRecorder) SourceTransformer(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SourceTransformer", reflect.TypeOf((*MockSourceTransformerClient)(nil).SourceTransformer), varargs...)
}