// Code generated by MockGen. DO NOT EDIT.
// Source: remote/fetcher/objectfetcher.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	types "gitlab.com/makeos/mosdef/remote/push/types"
	io "io"
	reflect "reflect"
)

// MockObjectFetcher is a mock of ObjectFetcher interface
type MockObjectFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockObjectFetcherMockRecorder
}

// MockObjectFetcherMockRecorder is the mock recorder for MockObjectFetcher
type MockObjectFetcherMockRecorder struct {
	mock *MockObjectFetcher
}

// NewMockObjectFetcher creates a new mock instance
func NewMockObjectFetcher(ctrl *gomock.Controller) *MockObjectFetcher {
	mock := &MockObjectFetcher{ctrl: ctrl}
	mock.recorder = &MockObjectFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockObjectFetcher) EXPECT() *MockObjectFetcherMockRecorder {
	return m.recorder
}

// Fetch mocks base method
func (m *MockObjectFetcher) Fetch(note types.PushNote, cb func(error)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Fetch", note, cb)
}

// Fetch indicates an expected call of Fetch
func (mr *MockObjectFetcherMockRecorder) Fetch(note, cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockObjectFetcher)(nil).Fetch), note, cb)
}

// OnPackReceived mocks base method
func (m *MockObjectFetcher) OnPackReceived(cb func(string, io.ReadSeeker)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnPackReceived", cb)
}

// OnPackReceived indicates an expected call of OnPackReceived
func (mr *MockObjectFetcherMockRecorder) OnPackReceived(cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnPackReceived", reflect.TypeOf((*MockObjectFetcher)(nil).OnPackReceived), cb)
}

// Start mocks base method
func (m *MockObjectFetcher) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start
func (mr *MockObjectFetcherMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockObjectFetcher)(nil).Start))
}

// Stop mocks base method
func (m *MockObjectFetcher) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop
func (mr *MockObjectFetcherMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockObjectFetcher)(nil).Stop))
}

// MockObjectFetcherService is a mock of ObjectFetcherService interface
type MockObjectFetcherService struct {
	ctrl     *gomock.Controller
	recorder *MockObjectFetcherServiceMockRecorder
}

// MockObjectFetcherServiceMockRecorder is the mock recorder for MockObjectFetcherService
type MockObjectFetcherServiceMockRecorder struct {
	mock *MockObjectFetcherService
}

// NewMockObjectFetcherService creates a new mock instance
func NewMockObjectFetcherService(ctrl *gomock.Controller) *MockObjectFetcherService {
	mock := &MockObjectFetcherService{ctrl: ctrl}
	mock.recorder = &MockObjectFetcherServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockObjectFetcherService) EXPECT() *MockObjectFetcherServiceMockRecorder {
	return m.recorder
}

// Fetch mocks base method
func (m *MockObjectFetcherService) Fetch(note types.PushNote, cb func(error)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Fetch", note, cb)
}

// Fetch indicates an expected call of Fetch
func (mr *MockObjectFetcherServiceMockRecorder) Fetch(note, cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockObjectFetcherService)(nil).Fetch), note, cb)
}