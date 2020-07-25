// Code generated by MockGen. DO NOT EDIT.
// Source: ticket/types/types.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	types "gitlab.com/makeos/lobe/ticket/types"
	types0 "gitlab.com/makeos/lobe/types"
	util "gitlab.com/makeos/lobe/util"
	reflect "reflect"
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
func (m *MockTicketManager) Index(tx types0.BaseTx, blockHeight uint64, txIndex int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Index", tx, blockHeight, txIndex)
	ret0, _ := ret[0].(error)
	return ret0
}

// Index indicates an expected call of Index
func (mr *MockTicketManagerMockRecorder) Index(tx, blockHeight, txIndex interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Index", reflect.TypeOf((*MockTicketManager)(nil).Index), tx, blockHeight, txIndex)
}

// Remove mocks base method
func (m *MockTicketManager) Remove(hash util.HexBytes) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", hash)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockTicketManagerMockRecorder) Remove(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockTicketManager)(nil).Remove), hash)
}

// GetByProposer mocks base method
func (m *MockTicketManager) GetByProposer(ticketType types0.TxCode, proposerPubKey util.Bytes32, queryOpt ...interface{}) ([]*types.Ticket, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ticketType, proposerPubKey}
	for _, a := range queryOpt {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByProposer", varargs...)
	ret0, _ := ret[0].([]*types.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByProposer indicates an expected call of GetByProposer
func (mr *MockTicketManagerMockRecorder) GetByProposer(ticketType, proposerPubKey interface{}, queryOpt ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ticketType, proposerPubKey}, queryOpt...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByProposer", reflect.TypeOf((*MockTicketManager)(nil).GetByProposer), varargs...)
}

// CountActiveValidatorTickets mocks base method
func (m *MockTicketManager) CountActiveValidatorTickets() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountActiveValidatorTickets")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountActiveValidatorTickets indicates an expected call of CountActiveValidatorTickets
func (mr *MockTicketManagerMockRecorder) CountActiveValidatorTickets() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountActiveValidatorTickets", reflect.TypeOf((*MockTicketManager)(nil).CountActiveValidatorTickets))
}

// GetNonDelegatedTickets mocks base method
func (m *MockTicketManager) GetNonDelegatedTickets(pubKey util.Bytes32, ticketType types0.TxCode) ([]*types.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNonDelegatedTickets", pubKey, ticketType)
	ret0, _ := ret[0].([]*types.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNonDelegatedTickets indicates an expected call of GetNonDelegatedTickets
func (mr *MockTicketManagerMockRecorder) GetNonDelegatedTickets(pubKey, ticketType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNonDelegatedTickets", reflect.TypeOf((*MockTicketManager)(nil).GetNonDelegatedTickets), pubKey, ticketType)
}

// Query mocks base method
func (m *MockTicketManager) Query(qf func(*types.Ticket) bool, queryOpt ...interface{}) []*types.Ticket {
	m.ctrl.T.Helper()
	varargs := []interface{}{qf}
	for _, a := range queryOpt {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].([]*types.Ticket)
	return ret0
}

// Query indicates an expected call of Query
func (mr *MockTicketManagerMockRecorder) Query(qf interface{}, queryOpt ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{qf}, queryOpt...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockTicketManager)(nil).Query), varargs...)
}

// QueryOne mocks base method
func (m *MockTicketManager) QueryOne(qf func(*types.Ticket) bool) *types.Ticket {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryOne", qf)
	ret0, _ := ret[0].(*types.Ticket)
	return ret0
}

// QueryOne indicates an expected call of QueryOne
func (mr *MockTicketManagerMockRecorder) QueryOne(qf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryOne", reflect.TypeOf((*MockTicketManager)(nil).QueryOne), qf)
}

// GetByHash mocks base method
func (m *MockTicketManager) GetByHash(hash util.HexBytes) *types.Ticket {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByHash", hash)
	ret0, _ := ret[0].(*types.Ticket)
	return ret0
}

// GetByHash indicates an expected call of GetByHash
func (mr *MockTicketManagerMockRecorder) GetByHash(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByHash", reflect.TypeOf((*MockTicketManager)(nil).GetByHash), hash)
}

// UpdateDecayBy mocks base method
func (m *MockTicketManager) UpdateDecayBy(hash util.HexBytes, newDecayHeight uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDecayBy", hash, newDecayHeight)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDecayBy indicates an expected call of UpdateDecayBy
func (mr *MockTicketManagerMockRecorder) UpdateDecayBy(hash, newDecayHeight interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDecayBy", reflect.TypeOf((*MockTicketManager)(nil).UpdateDecayBy), hash, newDecayHeight)
}

// GetTopHosts mocks base method
func (m *MockTicketManager) GetTopHosts(limit int) (types.SelectedTickets, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopHosts", limit)
	ret0, _ := ret[0].(types.SelectedTickets)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopHosts indicates an expected call of GetTopHosts
func (mr *MockTicketManagerMockRecorder) GetTopHosts(limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopHosts", reflect.TypeOf((*MockTicketManager)(nil).GetTopHosts), limit)
}

// GetTopValidators mocks base method
func (m *MockTicketManager) GetTopValidators(limit int) (types.SelectedTickets, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopValidators", limit)
	ret0, _ := ret[0].(types.SelectedTickets)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopValidators indicates an expected call of GetTopValidators
func (mr *MockTicketManagerMockRecorder) GetTopValidators(limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopValidators", reflect.TypeOf((*MockTicketManager)(nil).GetTopValidators), limit)
}

// ValueOfNonDelegatedTickets mocks base method
func (m *MockTicketManager) ValueOfNonDelegatedTickets(pubKey util.Bytes32, maturityHeight uint64) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValueOfNonDelegatedTickets", pubKey, maturityHeight)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValueOfNonDelegatedTickets indicates an expected call of ValueOfNonDelegatedTickets
func (mr *MockTicketManagerMockRecorder) ValueOfNonDelegatedTickets(pubKey, maturityHeight interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValueOfNonDelegatedTickets", reflect.TypeOf((*MockTicketManager)(nil).ValueOfNonDelegatedTickets), pubKey, maturityHeight)
}

// ValueOfDelegatedTickets mocks base method
func (m *MockTicketManager) ValueOfDelegatedTickets(pubKey util.Bytes32, maturityHeight uint64) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValueOfDelegatedTickets", pubKey, maturityHeight)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValueOfDelegatedTickets indicates an expected call of ValueOfDelegatedTickets
func (mr *MockTicketManagerMockRecorder) ValueOfDelegatedTickets(pubKey, maturityHeight interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValueOfDelegatedTickets", reflect.TypeOf((*MockTicketManager)(nil).ValueOfDelegatedTickets), pubKey, maturityHeight)
}

// ValueOfTickets mocks base method
func (m *MockTicketManager) ValueOfTickets(pubKey util.Bytes32, maturityHeight uint64) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValueOfTickets", pubKey, maturityHeight)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValueOfTickets indicates an expected call of ValueOfTickets
func (mr *MockTicketManagerMockRecorder) ValueOfTickets(pubKey, maturityHeight interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValueOfTickets", reflect.TypeOf((*MockTicketManager)(nil).ValueOfTickets), pubKey, maturityHeight)
}

// ValueOfAllTickets mocks base method
func (m *MockTicketManager) ValueOfAllTickets(maturityHeight uint64) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValueOfAllTickets", maturityHeight)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValueOfAllTickets indicates an expected call of ValueOfAllTickets
func (mr *MockTicketManagerMockRecorder) ValueOfAllTickets(maturityHeight interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValueOfAllTickets", reflect.TypeOf((*MockTicketManager)(nil).ValueOfAllTickets), maturityHeight)
}

// GetNonDecayedTickets mocks base method
func (m *MockTicketManager) GetNonDecayedTickets(pubKey util.Bytes32, maturityHeight uint64) ([]*types.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNonDecayedTickets", pubKey, maturityHeight)
	ret0, _ := ret[0].([]*types.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNonDecayedTickets indicates an expected call of GetNonDecayedTickets
func (mr *MockTicketManagerMockRecorder) GetNonDecayedTickets(pubKey, maturityHeight interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNonDecayedTickets", reflect.TypeOf((*MockTicketManager)(nil).GetNonDecayedTickets), pubKey, maturityHeight)
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
