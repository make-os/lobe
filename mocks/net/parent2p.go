// Code generated by MockGen. DO NOT EDIT.
// Source: net/parent2p/parent2p.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	peer "github.com/libp2p/go-libp2p-core/peer"
	parent2p "github.com/make-os/kit/net/parent2p"
	reflect "reflect"
)

// MockParent2P is a mock of Parent2P interface
type MockParent2P struct {
	ctrl     *gomock.Controller
	recorder *MockParent2PMockRecorder
}

// MockParent2PMockRecorder is the mock recorder for MockParent2P
type MockParent2PMockRecorder struct {
	mock *MockParent2P
}

// NewMockParent2P creates a new mock instance
func NewMockParent2P(ctrl *gomock.Controller) *MockParent2P {
	mock := &MockParent2P{ctrl: ctrl}
	mock.recorder = &MockParent2PMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockParent2P) EXPECT() *MockParent2PMockRecorder {
	return m.recorder
}

// ConnectToParent mocks base method
func (m *MockParent2P) ConnectToParent(ctx context.Context, parentAddr string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnectToParent", ctx, parentAddr)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConnectToParent indicates an expected call of ConnectToParent
func (mr *MockParent2PMockRecorder) ConnectToParent(ctx, parentAddr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectToParent", reflect.TypeOf((*MockParent2P)(nil).ConnectToParent), ctx, parentAddr)
}

// SendHandshakeMsg mocks base method
func (m *MockParent2P) SendHandshakeMsg(ctx context.Context, trackList []string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHandshakeMsg", ctx, trackList)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendHandshakeMsg indicates an expected call of SendHandshakeMsg
func (mr *MockParent2PMockRecorder) SendHandshakeMsg(ctx, trackList interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHandshakeMsg", reflect.TypeOf((*MockParent2P)(nil).SendHandshakeMsg), ctx, trackList)
}

// Parent mocks base method
func (m *MockParent2P) Parent() peer.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parent")
	ret0, _ := ret[0].(peer.ID)
	return ret0
}

// Parent indicates an expected call of Parent
func (mr *MockParent2PMockRecorder) Parent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parent", reflect.TypeOf((*MockParent2P)(nil).Parent))
}

// Peers mocks base method
func (m *MockParent2P) Peers() parent2p.Peers {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Peers")
	ret0, _ := ret[0].(parent2p.Peers)
	return ret0
}

// Peers indicates an expected call of Peers
func (mr *MockParent2PMockRecorder) Peers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Peers", reflect.TypeOf((*MockParent2P)(nil).Peers))
}