// Code generated by MockGen. DO NOT EDIT.
// Source: api/rest/client/types.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	req "github.com/imroc/req"
	types "gitlab.com/makeos/mosdef/api/types"
	state "gitlab.com/makeos/mosdef/types/state"
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

// SendTxPayload mocks base method
func (m *MockClient) SendTxPayload(data map[string]interface{}) (*types.SendTxPayloadResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendTxPayload", data)
	ret0, _ := ret[0].(*types.SendTxPayloadResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendTxPayload indicates an expected call of SendTxPayload
func (mr *MockClientMockRecorder) SendTxPayload(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTxPayload", reflect.TypeOf((*MockClient)(nil).SendTxPayload), data)
}

// GetAccountNonce mocks base method
func (m *MockClient) GetAccountNonce(address string, blockHeight ...uint64) (*types.GetAccountNonceResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{address}
	for _, a := range blockHeight {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAccountNonce", varargs...)
	ret0, _ := ret[0].(*types.GetAccountNonceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountNonce indicates an expected call of GetAccountNonce
func (mr *MockClientMockRecorder) GetAccountNonce(address interface{}, blockHeight ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{address}, blockHeight...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountNonce", reflect.TypeOf((*MockClient)(nil).GetAccountNonce), varargs...)
}

// GetAccount mocks base method
func (m *MockClient) GetAccount(address string, blockHeight ...uint64) (*state.Account, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{address}
	for _, a := range blockHeight {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAccount", varargs...)
	ret0, _ := ret[0].(*state.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount
func (mr *MockClientMockRecorder) GetAccount(address interface{}, blockHeight ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{address}, blockHeight...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockClient)(nil).GetAccount), varargs...)
}

// GetCall mocks base method
func (m *MockClient) GetCall(endpoint string, params map[string]interface{}) (*req.Resp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCall", endpoint, params)
	ret0, _ := ret[0].(*req.Resp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCall indicates an expected call of GetCall
func (mr *MockClientMockRecorder) GetCall(endpoint, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCall", reflect.TypeOf((*MockClient)(nil).GetCall), endpoint, params)
}

// PostCall mocks base method
func (m *MockClient) PostCall(endpoint string, body map[string]interface{}) (*req.Resp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostCall", endpoint, body)
	ret0, _ := ret[0].(*req.Resp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostCall indicates an expected call of PostCall
func (mr *MockClientMockRecorder) PostCall(endpoint, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostCall", reflect.TypeOf((*MockClient)(nil).PostCall), endpoint, body)
}

// GetPushKeyOwnerNonce mocks base method
func (m *MockClient) GetPushKeyOwnerNonce(pushKeyID string, blockHeight ...uint64) (*types.GetAccountNonceResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{pushKeyID}
	for _, a := range blockHeight {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetPushKeyOwnerNonce", varargs...)
	ret0, _ := ret[0].(*types.GetAccountNonceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPushKeyOwnerNonce indicates an expected call of GetPushKeyOwnerNonce
func (mr *MockClientMockRecorder) GetPushKeyOwnerNonce(pushKeyID interface{}, blockHeight ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{pushKeyID}, blockHeight...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPushKeyOwnerNonce", reflect.TypeOf((*MockClient)(nil).GetPushKeyOwnerNonce), varargs...)
}

// GetPushKey mocks base method
func (m *MockClient) GetPushKey(pushKeyID string, blockHeight ...uint64) (*state.PushKey, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{pushKeyID}
	for _, a := range blockHeight {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetPushKey", varargs...)
	ret0, _ := ret[0].(*state.PushKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPushKey indicates an expected call of GetPushKey
func (mr *MockClientMockRecorder) GetPushKey(pushKeyID interface{}, blockHeight ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{pushKeyID}, blockHeight...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPushKey", reflect.TypeOf((*MockClient)(nil).GetPushKey), varargs...)
}
