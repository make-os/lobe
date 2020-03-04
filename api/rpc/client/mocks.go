// Code generated by MockGen. DO NOT EDIT.
// Source: api/rpc/client/client.go

// Package client is a generated GoMock package.
package client

import (
	gomock "github.com/golang/mock/gomock"
	types "gitlab.com/makeos/mosdef/api/types"
	state "gitlab.com/makeos/mosdef/types/state"
	util "gitlab.com/makeos/mosdef/util"
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

// TxSendPayload mocks base method
func (m *MockClient) TxSendPayload(data map[string]interface{}) (*types.TxSendPayloadResponse, *util.StatusError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxSendPayload", data)
	ret0, _ := ret[0].(*types.TxSendPayloadResponse)
	ret1, _ := ret[1].(*util.StatusError)
	return ret0, ret1
}

// TxSendPayload indicates an expected call of TxSendPayload
func (mr *MockClientMockRecorder) TxSendPayload(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxSendPayload", reflect.TypeOf((*MockClient)(nil).TxSendPayload), data)
}

// AccountGet mocks base method
func (m *MockClient) AccountGet(address string, blockHeight ...uint64) (*state.Account, *util.StatusError) {
	m.ctrl.T.Helper()
	varargs := []interface{}{address}
	for _, a := range blockHeight {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AccountGet", varargs...)
	ret0, _ := ret[0].(*state.Account)
	ret1, _ := ret[1].(*util.StatusError)
	return ret0, ret1
}

// AccountGet indicates an expected call of AccountGet
func (mr *MockClientMockRecorder) AccountGet(address interface{}, blockHeight ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{address}, blockHeight...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountGet", reflect.TypeOf((*MockClient)(nil).AccountGet), varargs...)
}

// GPGGetAccountOfOwner mocks base method
func (m *MockClient) GPGGetAccountOfOwner(id string, blockHeight ...uint64) (*state.Account, *util.StatusError) {
	m.ctrl.T.Helper()
	varargs := []interface{}{id}
	for _, a := range blockHeight {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GPGGetAccountOfOwner", varargs...)
	ret0, _ := ret[0].(*state.Account)
	ret1, _ := ret[1].(*util.StatusError)
	return ret0, ret1
}

// GPGGetAccountOfOwner indicates an expected call of GPGGetAccountOfOwner
func (mr *MockClientMockRecorder) GPGGetAccountOfOwner(id interface{}, blockHeight ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{id}, blockHeight...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GPGGetAccountOfOwner", reflect.TypeOf((*MockClient)(nil).GPGGetAccountOfOwner), varargs...)
}

// GetOptions mocks base method
func (m *MockClient) GetOptions() *Options {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOptions")
	ret0, _ := ret[0].(*Options)
	return ret0
}

// GetOptions indicates an expected call of GetOptions
func (mr *MockClientMockRecorder) GetOptions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOptions", reflect.TypeOf((*MockClient)(nil).GetOptions))
}

// Call mocks base method
func (m *MockClient) Call(method string, params interface{}) (util.Map, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", method, params)
	ret0, _ := ret[0].(util.Map)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Call indicates an expected call of Call
func (mr *MockClientMockRecorder) Call(method, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockClient)(nil).Call), method, params)
}
