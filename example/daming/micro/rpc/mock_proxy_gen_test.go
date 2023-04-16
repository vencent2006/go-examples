// Code generated by MockGen. DO NOT EDIT.
// Source: micro/rpc/types.go

// Package rpc is a generated GoMock package.
package rpc

import (
	context "context"
	"example/daming/micro/rpc/message"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Name mocks base method.
func (m *MockService) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockServiceMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockService)(nil).Name))
}

// MockProxy is a mock of Proxy interface.
type MockProxy struct {
	ctrl     *gomock.Controller
	recorder *MockProxyMockRecorder
}

// MockProxyMockRecorder is the mock recorder for MockProxy.
type MockProxyMockRecorder struct {
	mock *MockProxy
}

// NewMockProxy creates a new mock instance.
func NewMockProxy(ctrl *gomock.Controller) *MockProxy {
	mock := &MockProxy{ctrl: ctrl}
	mock.recorder = &MockProxyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProxy) EXPECT() *MockProxyMockRecorder {
	return m.recorder
}

// Invoke mocks base method.
func (m *MockProxy) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Invoke", ctx, req)
	ret0, _ := ret[0].(*message.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Invoke indicates an expected call of Invoke.
func (mr *MockProxyMockRecorder) Invoke(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Invoke", reflect.TypeOf((*MockProxy)(nil).Invoke), ctx, req)
}
