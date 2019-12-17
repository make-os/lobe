// Code generated by MockGen. DO NOT EDIT.
// Source: jsonrpc.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	types "github.com/makeos/mosdef/types"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Call mocks base method
func (m *MockClient) Call(method string, params interface{}) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", method, params)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Call indicates an expected call of Call
func (mr *MockClientMockRecorder) Call(method, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockClient)(nil).Call), method, params)
}

// New mocks base method
func (m *MockClient) New(opts *types.Options) types.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New", opts)
	ret0, _ := ret[0].(types.Client)
	return ret0
}

// New indicates an expected call of New
func (mr *MockClientMockRecorder) New(opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockClient)(nil).New), opts)
}

// GetOptions mocks base method
func (m *MockClient) GetOptions() *types.Options {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOptions")
	ret0, _ := ret[0].(*types.Options)
	return ret0
}

// GetOptions indicates an expected call of GetOptions
func (mr *MockClientMockRecorder) GetOptions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOptions", reflect.TypeOf((*MockClient)(nil).GetOptions))
}