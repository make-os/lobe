// Code generated by MockGen. DO NOT EDIT.
// Source: ticket.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/makeos/mosdef/types"
)

// MockTicketManager is a mock of TicketManager interface
type MockTicketManager struct {
	ctrl     *gomock.Controller
	recorder *MockTicketManagerMockRecorder
}

// MockTicketManagerMockRecorder is the mock recorder for MockTicketManager
type MockTicketManagerMockRecorder struct {
	mock *MockTicketManager
}

// NewMockTicketManager creates a new mock instance
func NewMockTicketManager(ctrl *gomock.Controller) *MockTicketManager {
	mock := &MockTicketManager{ctrl: ctrl}
	mock.recorder = &MockTicketManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTicketManager) EXPECT() *MockTicketManagerMockRecorder {
	return m.recorder
}

// Index mocks base method
func (m *MockTicketManager) Index(tx *types.Transaction, proposerPubKey string, blockHeight uint64, txIndex int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Index", tx, proposerPubKey, blockHeight, txIndex)
	ret0, _ := ret[0].(error)
	return ret0
}

// Index indicates an expected call of Index
func (mr *MockTicketManagerMockRecorder) Index(tx, proposerPubKey, blockHeight, txIndex interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Index", reflect.TypeOf((*MockTicketManager)(nil).Index), tx, proposerPubKey, blockHeight, txIndex)
}

// Get mocks base method
func (m *MockTicketManager) Get(proposerPubKey string, queryOpt types.QueryOptions) ([]*types.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", proposerPubKey, queryOpt)
	ret0, _ := ret[0].([]*types.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockTicketManagerMockRecorder) Get(proposerPubKey, queryOpt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTicketManager)(nil).Get), proposerPubKey, queryOpt)
}

// Stop mocks base method
func (m *MockTicketManager) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
func (mr *MockTicketManagerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockTicketManager)(nil).Stop))
}