// Code generated by MockGen. DO NOT EDIT.
// Source: api/remote/client/client.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	req "github.com/imroc/req"
	types "gitlab.com/makeos/mosdef/api/types"
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

// SendTxPayload mocks base method
func (m *MockClient) SendTxPayload(data map[string]interface{}) (*types.HashResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendTxPayload", data)
	ret0, _ := ret[0].(*types.HashResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendTxPayload indicates an expected call of SendTxPayload
func (mr *MockClientMockRecorder) SendTxPayload(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTxPayload", reflect.TypeOf((*MockClient)(nil).SendTxPayload), data)
}

// GetTransaction mocks base method
func (m *MockClient) GetTransaction(hash string) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransaction", hash)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransaction indicates an expected call of GetTransaction
func (mr *MockClientMockRecorder) GetTransaction(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*MockClient)(nil).GetTransaction), hash)
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
func (m *MockClient) GetAccount(address string, blockHeight ...uint64) (*types.GetAccountResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{address}
	for _, a := range blockHeight {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAccount", varargs...)
	ret0, _ := ret[0].(*types.GetAccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount
func (mr *MockClientMockRecorder) GetAccount(address interface{}, blockHeight ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{address}, blockHeight...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockClient)(nil).GetAccount), varargs...)
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
func (m *MockClient) GetPushKey(pushKeyID string, blockHeight ...uint64) (*types.GetPushKeyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{pushKeyID}
	for _, a := range blockHeight {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetPushKey", varargs...)
	ret0, _ := ret[0].(*types.GetPushKeyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPushKey indicates an expected call of GetPushKey
func (mr *MockClientMockRecorder) GetPushKey(pushKeyID interface{}, blockHeight ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{pushKeyID}, blockHeight...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPushKey", reflect.TypeOf((*MockClient)(nil).GetPushKey), varargs...)
}

// RegisterPushKey mocks base method
func (m *MockClient) RegisterPushKey(body *types.RegisterPushKeyBody) (*types.RegisterPushKeyResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterPushKey", body)
	ret0, _ := ret[0].(*types.RegisterPushKeyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterPushKey indicates an expected call of RegisterPushKey
func (mr *MockClientMockRecorder) RegisterPushKey(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterPushKey", reflect.TypeOf((*MockClient)(nil).RegisterPushKey), body)
}

// CreateRepo mocks base method
func (m *MockClient) CreateRepo(body *types.CreateRepoBody) (*types.CreateRepoResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRepo", body)
	ret0, _ := ret[0].(*types.CreateRepoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRepo indicates an expected call of CreateRepo
func (mr *MockClientMockRecorder) CreateRepo(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRepo", reflect.TypeOf((*MockClient)(nil).CreateRepo), body)
}

// GetRepo mocks base method
func (m *MockClient) GetRepo(name string, opts ...*types.GetRepoOpts) (*types.GetRepoResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{name}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetRepo", varargs...)
	ret0, _ := ret[0].(*types.GetRepoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepo indicates an expected call of GetRepo
func (mr *MockClientMockRecorder) GetRepo(name interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{name}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepo", reflect.TypeOf((*MockClient)(nil).GetRepo), varargs...)
}

// AddRepoContributors mocks base method
func (m *MockClient) AddRepoContributors(body *types.AddRepoContribsBody) (*types.HashResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRepoContributors", body)
	ret0, _ := ret[0].(*types.HashResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddRepoContributors indicates an expected call of AddRepoContributors
func (mr *MockClientMockRecorder) AddRepoContributors(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRepoContributors", reflect.TypeOf((*MockClient)(nil).AddRepoContributors), body)
}

// SendCoin mocks base method
func (m *MockClient) SendCoin(body *types.SendCoinBody) (*types.HashResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoin", body)
	ret0, _ := ret[0].(*types.HashResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendCoin indicates an expected call of SendCoin
func (mr *MockClientMockRecorder) SendCoin(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoin", reflect.TypeOf((*MockClient)(nil).SendCoin), body)
}
