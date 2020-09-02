// Code generated by MockGen. DO NOT EDIT.
// Source: types/core/remote.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	config "github.com/make-os/lobe/config"
	crypto "github.com/make-os/lobe/crypto"
	types "github.com/make-os/lobe/dht/types"
	types0 "github.com/make-os/lobe/modules/types"
	logger "github.com/make-os/lobe/pkgs/logger"
	fetcher "github.com/make-os/lobe/remote/fetcher"
	types1 "github.com/make-os/lobe/remote/push/types"
	types2 "github.com/make-os/lobe/remote/types"
	core "github.com/make-os/lobe/types/core"
	reflect "reflect"
)

// MockPoolGetter is a mock of PoolGetter interface
type MockPoolGetter struct {
	ctrl     *gomock.Controller
	recorder *MockPoolGetterMockRecorder
}

// MockPoolGetterMockRecorder is the mock recorder for MockPoolGetter
type MockPoolGetterMockRecorder struct {
	mock *MockPoolGetter
}

// NewMockPoolGetter creates a new mock instance
func NewMockPoolGetter(ctrl *gomock.Controller) *MockPoolGetter {
	mock := &MockPoolGetter{ctrl: ctrl}
	mock.recorder = &MockPoolGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPoolGetter) EXPECT() *MockPoolGetterMockRecorder {
	return m.recorder
}

// GetPushPool mocks base method
func (m *MockPoolGetter) GetPushPool() types1.PushPool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPushPool")
	ret0, _ := ret[0].(types1.PushPool)
	return ret0
}

// GetPushPool indicates an expected call of GetPushPool
func (mr *MockPoolGetterMockRecorder) GetPushPool() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPushPool", reflect.TypeOf((*MockPoolGetter)(nil).GetPushPool))
}

// GetMempool mocks base method
func (m *MockPoolGetter) GetMempool() core.Mempool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMempool")
	ret0, _ := ret[0].(core.Mempool)
	return ret0
}

// GetMempool indicates an expected call of GetMempool
func (mr *MockPoolGetterMockRecorder) GetMempool() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMempool", reflect.TypeOf((*MockPoolGetter)(nil).GetMempool))
}

// MockRemoteServer is a mock of RemoteServer interface
type MockRemoteServer struct {
	ctrl     *gomock.Controller
	recorder *MockRemoteServerMockRecorder
}

// MockRemoteServerMockRecorder is the mock recorder for MockRemoteServer
type MockRemoteServerMockRecorder struct {
	mock *MockRemoteServer
}

// NewMockRemoteServer creates a new mock instance
func NewMockRemoteServer(ctrl *gomock.Controller) *MockRemoteServer {
	mock := &MockRemoteServer{ctrl: ctrl}
	mock.recorder = &MockRemoteServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRemoteServer) EXPECT() *MockRemoteServerMockRecorder {
	return m.recorder
}

// GetPushPool mocks base method
func (m *MockRemoteServer) GetPushPool() types1.PushPool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPushPool")
	ret0, _ := ret[0].(types1.PushPool)
	return ret0
}

// GetPushPool indicates an expected call of GetPushPool
func (mr *MockRemoteServerMockRecorder) GetPushPool() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPushPool", reflect.TypeOf((*MockRemoteServer)(nil).GetPushPool))
}

// GetMempool mocks base method
func (m *MockRemoteServer) GetMempool() core.Mempool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMempool")
	ret0, _ := ret[0].(core.Mempool)
	return ret0
}

// GetMempool indicates an expected call of GetMempool
func (mr *MockRemoteServerMockRecorder) GetMempool() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMempool", reflect.TypeOf((*MockRemoteServer)(nil).GetMempool))
}

// Log mocks base method
func (m *MockRemoteServer) Log() logger.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Log")
	ret0, _ := ret[0].(logger.Logger)
	return ret0
}

// Log indicates an expected call of Log
func (mr *MockRemoteServerMockRecorder) Log() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockRemoteServer)(nil).Log))
}

// Cfg mocks base method
func (m *MockRemoteServer) Cfg() *config.AppConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cfg")
	ret0, _ := ret[0].(*config.AppConfig)
	return ret0
}

// Cfg indicates an expected call of Cfg
func (mr *MockRemoteServerMockRecorder) Cfg() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cfg", reflect.TypeOf((*MockRemoteServer)(nil).Cfg))
}

// GetRepoState mocks base method
func (m *MockRemoteServer) GetRepoState(target types2.LocalRepo, options ...types2.KVOption) (types2.BareRepoRefsState, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{target}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetRepoState", varargs...)
	ret0, _ := ret[0].(types2.BareRepoRefsState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepoState indicates an expected call of GetRepoState
func (mr *MockRemoteServerMockRecorder) GetRepoState(target interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{target}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepoState", reflect.TypeOf((*MockRemoteServer)(nil).GetRepoState), varargs...)
}

// GetPushKeyGetter mocks base method
func (m *MockRemoteServer) GetPushKeyGetter() core.PushKeyGetter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPushKeyGetter")
	ret0, _ := ret[0].(core.PushKeyGetter)
	return ret0
}

// GetPushKeyGetter indicates an expected call of GetPushKeyGetter
func (mr *MockRemoteServerMockRecorder) GetPushKeyGetter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPushKeyGetter", reflect.TypeOf((*MockRemoteServer)(nil).GetPushKeyGetter))
}

// GetLogic mocks base method
func (m *MockRemoteServer) GetLogic() core.Logic {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogic")
	ret0, _ := ret[0].(core.Logic)
	return ret0
}

// GetLogic indicates an expected call of GetLogic
func (mr *MockRemoteServerMockRecorder) GetLogic() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogic", reflect.TypeOf((*MockRemoteServer)(nil).GetLogic))
}

// GetRepo mocks base method
func (m *MockRemoteServer) GetRepo(name string) (types2.LocalRepo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepo", name)
	ret0, _ := ret[0].(types2.LocalRepo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepo indicates an expected call of GetRepo
func (mr *MockRemoteServerMockRecorder) GetRepo(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepo", reflect.TypeOf((*MockRemoteServer)(nil).GetRepo), name)
}

// GetPrivateValidatorKey mocks base method
func (m *MockRemoteServer) GetPrivateValidatorKey() *crypto.Key {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPrivateValidatorKey")
	ret0, _ := ret[0].(*crypto.Key)
	return ret0
}

// GetPrivateValidatorKey indicates an expected call of GetPrivateValidatorKey
func (mr *MockRemoteServerMockRecorder) GetPrivateValidatorKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPrivateValidatorKey", reflect.TypeOf((*MockRemoteServer)(nil).GetPrivateValidatorKey))
}

// Start mocks base method
func (m *MockRemoteServer) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockRemoteServerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockRemoteServer)(nil).Start))
}

// Wait mocks base method
func (m *MockRemoteServer) Wait() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Wait")
}

// Wait indicates an expected call of Wait
func (mr *MockRemoteServerMockRecorder) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockRemoteServer)(nil).Wait))
}

// InitRepository mocks base method
func (m *MockRemoteServer) InitRepository(name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitRepository", name)
	ret0, _ := ret[0].(error)
	return ret0
}

// InitRepository indicates an expected call of InitRepository
func (mr *MockRemoteServerMockRecorder) InitRepository(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitRepository", reflect.TypeOf((*MockRemoteServer)(nil).InitRepository), name)
}

// BroadcastMsg mocks base method
func (m *MockRemoteServer) BroadcastMsg(ch byte, msg []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BroadcastMsg", ch, msg)
}

// BroadcastMsg indicates an expected call of BroadcastMsg
func (mr *MockRemoteServerMockRecorder) BroadcastMsg(ch, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BroadcastMsg", reflect.TypeOf((*MockRemoteServer)(nil).BroadcastMsg), ch, msg)
}

// BroadcastNoteAndEndorsement mocks base method
func (m *MockRemoteServer) BroadcastNoteAndEndorsement(note types1.PushNote) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BroadcastNoteAndEndorsement", note)
	ret0, _ := ret[0].(error)
	return ret0
}

// BroadcastNoteAndEndorsement indicates an expected call of BroadcastNoteAndEndorsement
func (mr *MockRemoteServerMockRecorder) BroadcastNoteAndEndorsement(note interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BroadcastNoteAndEndorsement", reflect.TypeOf((*MockRemoteServer)(nil).BroadcastNoteAndEndorsement), note)
}

// RegisterAPIHandlers mocks base method
func (m *MockRemoteServer) RegisterAPIHandlers(agg types0.ModulesHub) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterAPIHandlers", agg)
}

// RegisterAPIHandlers indicates an expected call of RegisterAPIHandlers
func (mr *MockRemoteServerMockRecorder) RegisterAPIHandlers(agg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterAPIHandlers", reflect.TypeOf((*MockRemoteServer)(nil).RegisterAPIHandlers), agg)
}

// Announce mocks base method
func (m *MockRemoteServer) Announce(objType int, repo string, hash []byte, doneCB func(error)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Announce", objType, repo, hash, doneCB)
}

// Announce indicates an expected call of Announce
func (mr *MockRemoteServerMockRecorder) Announce(objType, repo, hash, doneCB interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Announce", reflect.TypeOf((*MockRemoteServer)(nil).Announce), objType, repo, hash, doneCB)
}

// GetFetcher mocks base method
func (m *MockRemoteServer) GetFetcher() fetcher.ObjectFetcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFetcher")
	ret0, _ := ret[0].(fetcher.ObjectFetcher)
	return ret0
}

// GetFetcher indicates an expected call of GetFetcher
func (mr *MockRemoteServerMockRecorder) GetFetcher() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFetcher", reflect.TypeOf((*MockRemoteServer)(nil).GetFetcher))
}

// GetDHT mocks base method
func (m *MockRemoteServer) GetDHT() types.DHT {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDHT")
	ret0, _ := ret[0].(types.DHT)
	return ret0
}

// GetDHT indicates an expected call of GetDHT
func (mr *MockRemoteServerMockRecorder) GetDHT() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDHT", reflect.TypeOf((*MockRemoteServer)(nil).GetDHT))
}

// Shutdown mocks base method
func (m *MockRemoteServer) Shutdown(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Shutdown", ctx)
}

// Shutdown indicates an expected call of Shutdown
func (mr *MockRemoteServerMockRecorder) Shutdown(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockRemoteServer)(nil).Shutdown), ctx)
}

// Stop mocks base method
func (m *MockRemoteServer) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
func (mr *MockRemoteServerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockRemoteServer)(nil).Stop))
}
