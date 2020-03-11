// Code generated by MockGen. DO NOT EDIT.
// Source: pkgs/tree/types.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTree is a mock of Tree interface
type MockTree struct {
	ctrl     *gomock.Controller
	recorder *MockTreeMockRecorder
}

// MockTreeMockRecorder is the mock recorder for MockTree
type MockTreeMockRecorder struct {
	mock *MockTree
}

// NewMockTree creates a new mock instance
func NewMockTree(ctrl *gomock.Controller) *MockTree {
	mock := &MockTree{ctrl: ctrl}
	mock.recorder = &MockTreeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTree) EXPECT() *MockTreeMockRecorder {
	return m.recorder
}

// Version mocks base method
func (m *MockTree) Version() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Version")
	ret0, _ := ret[0].(int64)
	return ret0
}

// Version indicates an expected call of Version
func (mr *MockTreeMockRecorder) Version() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Version", reflect.TypeOf((*MockTree)(nil).Version))
}

// GetVersioned mocks base method
func (m *MockTree) GetVersioned(key []byte, version int64) (int64, []byte) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersioned", key, version)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].([]byte)
	return ret0, ret1
}

// GetVersioned indicates an expected call of GetVersioned
func (mr *MockTreeMockRecorder) GetVersioned(key, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersioned", reflect.TypeOf((*MockTree)(nil).GetVersioned), key, version)
}

// Get mocks base method
func (m *MockTree) Get(key []byte) (int64, []byte) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].([]byte)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockTreeMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTree)(nil).Get), key)
}

// Set mocks base method
func (m *MockTree) Set(key, value []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", key, value)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Set indicates an expected call of Set
func (mr *MockTreeMockRecorder) Set(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockTree)(nil).Set), key, value)
}

// Remove mocks base method
func (m *MockTree) Remove(key []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", key)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockTreeMockRecorder) Remove(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockTree)(nil).Remove), key)
}

// SaveVersion mocks base method
func (m *MockTree) SaveVersion() ([]byte, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveVersion")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SaveVersion indicates an expected call of SaveVersion
func (mr *MockTreeMockRecorder) SaveVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveVersion", reflect.TypeOf((*MockTree)(nil).SaveVersion))
}

// Load mocks base method
func (m *MockTree) Load() (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load indicates an expected call of Load
func (mr *MockTreeMockRecorder) Load() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockTree)(nil).Load))
}

// WorkingHash mocks base method
func (m *MockTree) WorkingHash() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WorkingHash")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// WorkingHash indicates an expected call of WorkingHash
func (mr *MockTreeMockRecorder) WorkingHash() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WorkingHash", reflect.TypeOf((*MockTree)(nil).WorkingHash))
}

// Hash mocks base method
func (m *MockTree) Hash() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hash")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Hash indicates an expected call of Hash
func (mr *MockTreeMockRecorder) Hash() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hash", reflect.TypeOf((*MockTree)(nil).Hash))
}
