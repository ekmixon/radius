// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/project-radius/radius/pkg/linkrp/processors (interfaces: ResourceClient)

// Package processors is a generated GoMock package.
package processors

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockResourceClient is a mock of ResourceClient interface.
type MockResourceClient struct {
	ctrl     *gomock.Controller
	recorder *MockResourceClientMockRecorder
}

// MockResourceClientMockRecorder is the mock recorder for MockResourceClient.
type MockResourceClientMockRecorder struct {
	mock *MockResourceClient
}

// NewMockResourceClient creates a new mock instance.
func NewMockResourceClient(ctrl *gomock.Controller) *MockResourceClient {
	mock := &MockResourceClient{ctrl: ctrl}
	mock.recorder = &MockResourceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResourceClient) EXPECT() *MockResourceClientMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockResourceClient) Delete(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockResourceClientMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockResourceClient)(nil).Delete), arg0, arg1, arg2)
}
