// Code generated by MockGen. DO NOT EDIT.
// Source: net/dht/announcer.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	dht "github.com/make-os/kit/net/dht"
	reflect "reflect"
)

// MockAnnouncer is a mock of Announcer interface
type MockAnnouncer struct {
	ctrl     *gomock.Controller
	recorder *MockAnnouncerMockRecorder
}

// MockAnnouncerMockRecorder is the mock recorder for MockAnnouncer
type MockAnnouncerMockRecorder struct {
	mock *MockAnnouncer
}

// NewMockAnnouncer creates a new mock instance
func NewMockAnnouncer(ctrl *gomock.Controller) *MockAnnouncer {
	mock := &MockAnnouncer{ctrl: ctrl}
	mock.recorder = &MockAnnouncerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAnnouncer) EXPECT() *MockAnnouncerMockRecorder {
	return m.recorder
}

// Announce mocks base method
func (m *MockAnnouncer) Announce(objType int, repo string, key []byte, doneCB func(error)) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Announce", objType, repo, key, doneCB)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Announce indicates an expected call of Announce
func (mr *MockAnnouncerMockRecorder) Announce(objType, repo, key, doneCB interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Announce", reflect.TypeOf((*MockAnnouncer)(nil).Announce), objType, repo, key, doneCB)
}

// Start mocks base method
func (m *MockAnnouncer) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start
func (mr *MockAnnouncerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockAnnouncer)(nil).Start))
}

// IsRunning mocks base method
func (m *MockAnnouncer) IsRunning() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRunning")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsRunning indicates an expected call of IsRunning
func (mr *MockAnnouncerMockRecorder) IsRunning() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRunning", reflect.TypeOf((*MockAnnouncer)(nil).IsRunning))
}

// HasTask mocks base method
func (m *MockAnnouncer) HasTask() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasTask")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasTask indicates an expected call of HasTask
func (mr *MockAnnouncerMockRecorder) HasTask() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasTask", reflect.TypeOf((*MockAnnouncer)(nil).HasTask))
}

// NewSession mocks base method
func (m *MockAnnouncer) NewSession() dht.Session {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSession")
	ret0, _ := ret[0].(dht.Session)
	return ret0
}

// NewSession indicates an expected call of NewSession
func (mr *MockAnnouncerMockRecorder) NewSession() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSession", reflect.TypeOf((*MockAnnouncer)(nil).NewSession))
}

// Stop mocks base method
func (m *MockAnnouncer) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop
func (mr *MockAnnouncerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockAnnouncer)(nil).Stop))
}

// RegisterChecker mocks base method
func (m *MockAnnouncer) RegisterChecker(objType int, checker dht.CheckFunc) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterChecker", objType, checker)
}

// RegisterChecker indicates an expected call of RegisterChecker
func (mr *MockAnnouncerMockRecorder) RegisterChecker(objType, checker interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterChecker", reflect.TypeOf((*MockAnnouncer)(nil).RegisterChecker), objType, checker)
}

// MockAnnouncerService is a mock of AnnouncerService interface
type MockAnnouncerService struct {
	ctrl     *gomock.Controller
	recorder *MockAnnouncerServiceMockRecorder
}

// MockAnnouncerServiceMockRecorder is the mock recorder for MockAnnouncerService
type MockAnnouncerServiceMockRecorder struct {
	mock *MockAnnouncerService
}

// NewMockAnnouncerService creates a new mock instance
func NewMockAnnouncerService(ctrl *gomock.Controller) *MockAnnouncerService {
	mock := &MockAnnouncerService{ctrl: ctrl}
	mock.recorder = &MockAnnouncerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAnnouncerService) EXPECT() *MockAnnouncerServiceMockRecorder {
	return m.recorder
}

// Announce mocks base method
func (m *MockAnnouncerService) Announce(objType int, repo string, key []byte, doneCB func(error)) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Announce", objType, repo, key, doneCB)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Announce indicates an expected call of Announce
func (mr *MockAnnouncerServiceMockRecorder) Announce(objType, repo, key, doneCB interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Announce", reflect.TypeOf((*MockAnnouncerService)(nil).Announce), objType, repo, key, doneCB)
}

// MockSession is a mock of Session interface
type MockSession struct {
	ctrl     *gomock.Controller
	recorder *MockSessionMockRecorder
}

// MockSessionMockRecorder is the mock recorder for MockSession
type MockSessionMockRecorder struct {
	mock *MockSession
}

// NewMockSession creates a new mock instance
func NewMockSession(ctrl *gomock.Controller) *MockSession {
	mock := &MockSession{ctrl: ctrl}
	mock.recorder = &MockSessionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSession) EXPECT() *MockSessionMockRecorder {
	return m.recorder
}

// Announce mocks base method
func (m *MockSession) Announce(objType int, repo string, key []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Announce", objType, repo, key)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Announce indicates an expected call of Announce
func (mr *MockSessionMockRecorder) Announce(objType, repo, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Announce", reflect.TypeOf((*MockSession)(nil).Announce), objType, repo, key)
}

// OnDone mocks base method
func (m *MockSession) OnDone(cb func(int)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnDone", cb)
}

// OnDone indicates an expected call of OnDone
func (mr *MockSessionMockRecorder) OnDone(cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnDone", reflect.TypeOf((*MockSession)(nil).OnDone), cb)
}
