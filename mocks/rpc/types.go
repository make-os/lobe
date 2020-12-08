// Code generated by MockGen. DO NOT EDIT.
// Source: rpc/types/types.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	rpc "github.com/make-os/kit/rpc"
	types "github.com/make-os/kit/rpc/types"
	api "github.com/make-os/kit/types/api"
	util "github.com/make-os/kit/util"
	reflect "reflect"
)

// MockPushKey is a mock of PushKey interface
type MockPushKey struct {
	ctrl     *gomock.Controller
	recorder *MockPushKeyMockRecorder
}

// MockPushKeyMockRecorder is the mock recorder for MockPushKey
type MockPushKeyMockRecorder struct {
	mock *MockPushKey
}

// NewMockPushKey creates a new mock instance
func NewMockPushKey(ctrl *gomock.Controller) *MockPushKey {
	mock := &MockPushKey{ctrl: ctrl}
	mock.recorder = &MockPushKeyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPushKey) EXPECT() *MockPushKeyMockRecorder {
	return m.recorder
}

// GetOwner mocks base method
func (m *MockPushKey) GetOwner(addr string, blockHeight ...uint64) (*api.ResultAccount, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{addr}
	for _, a := range blockHeight {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetOwner", varargs...)
	ret0, _ := ret[0].(*api.ResultAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOwner indicates an expected call of GetOwner
func (mr *MockPushKeyMockRecorder) GetOwner(addr interface{}, blockHeight ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{addr}, blockHeight...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOwner", reflect.TypeOf((*MockPushKey)(nil).GetOwner), varargs...)
}

// Register mocks base method
func (m *MockPushKey) Register(body *api.BodyRegisterPushKey) (*api.ResultRegisterPushKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", body)
	ret0, _ := ret[0].(*api.ResultRegisterPushKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register
func (mr *MockPushKeyMockRecorder) Register(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockPushKey)(nil).Register), body)
}

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

// GetOptions mocks base method
func (m *MockClient) GetOptions() *types.Options {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOptions")
	ret0, _ := ret[0].(*types.Options)
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

// Node mocks base method
func (m *MockClient) Node() types.Node {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Node")
	ret0, _ := ret[0].(types.Node)
	return ret0
}

// Node indicates an expected call of Node
func (mr *MockClientMockRecorder) Node() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Node", reflect.TypeOf((*MockClient)(nil).Node))
}

// PushKey mocks base method
func (m *MockClient) PushKey() types.PushKey {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushKey")
	ret0, _ := ret[0].(types.PushKey)
	return ret0
}

// PushKey indicates an expected call of PushKey
func (mr *MockClientMockRecorder) PushKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushKey", reflect.TypeOf((*MockClient)(nil).PushKey))
}

// Pool mocks base method
func (m *MockClient) Pool() types.Pool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pool")
	ret0, _ := ret[0].(types.Pool)
	return ret0
}

// Pool indicates an expected call of Pool
func (mr *MockClientMockRecorder) Pool() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pool", reflect.TypeOf((*MockClient)(nil).Pool))
}

// Repo mocks base method
func (m *MockClient) Repo() types.Repo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Repo")
	ret0, _ := ret[0].(types.Repo)
	return ret0
}

// Repo indicates an expected call of Repo
func (mr *MockClientMockRecorder) Repo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Repo", reflect.TypeOf((*MockClient)(nil).Repo))
}

// RPC mocks base method
func (m *MockClient) RPC() types.RPC {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RPC")
	ret0, _ := ret[0].(types.RPC)
	return ret0
}

// RPC indicates an expected call of RPC
func (mr *MockClientMockRecorder) RPC() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RPC", reflect.TypeOf((*MockClient)(nil).RPC))
}

// Tx mocks base method
func (m *MockClient) Tx() types.Tx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tx")
	ret0, _ := ret[0].(types.Tx)
	return ret0
}

// Tx indicates an expected call of Tx
func (mr *MockClientMockRecorder) Tx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tx", reflect.TypeOf((*MockClient)(nil).Tx))
}

// User mocks base method
func (m *MockClient) User() types.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "User")
	ret0, _ := ret[0].(types.User)
	return ret0
}

// User indicates an expected call of User
func (mr *MockClientMockRecorder) User() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "User", reflect.TypeOf((*MockClient)(nil).User))
}

// DHT mocks base method
func (m *MockClient) DHT() types.DHT {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DHT")
	ret0, _ := ret[0].(types.DHT)
	return ret0
}

// DHT indicates an expected call of DHT
func (mr *MockClientMockRecorder) DHT() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DHT", reflect.TypeOf((*MockClient)(nil).DHT))
}

// Ticket mocks base method
func (m *MockClient) Ticket() types.Ticket {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ticket")
	ret0, _ := ret[0].(types.Ticket)
	return ret0
}

// Ticket indicates an expected call of Ticket
func (mr *MockClientMockRecorder) Ticket() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ticket", reflect.TypeOf((*MockClient)(nil).Ticket))
}

// MockNode is a mock of Node interface
type MockNode struct {
	ctrl     *gomock.Controller
	recorder *MockNodeMockRecorder
}

// MockNodeMockRecorder is the mock recorder for MockNode
type MockNodeMockRecorder struct {
	mock *MockNode
}

// NewMockNode creates a new mock instance
func NewMockNode(ctrl *gomock.Controller) *MockNode {
	mock := &MockNode{ctrl: ctrl}
	mock.recorder = &MockNodeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNode) EXPECT() *MockNodeMockRecorder {
	return m.recorder
}

// GetBlock mocks base method
func (m *MockNode) GetBlock(height uint64) (*api.ResultBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlock", height)
	ret0, _ := ret[0].(*api.ResultBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlock indicates an expected call of GetBlock
func (mr *MockNodeMockRecorder) GetBlock(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*MockNode)(nil).GetBlock), height)
}

// GetHeight mocks base method
func (m *MockNode) GetHeight() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeight")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHeight indicates an expected call of GetHeight
func (mr *MockNodeMockRecorder) GetHeight() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeight", reflect.TypeOf((*MockNode)(nil).GetHeight))
}

// GetBlockInfo mocks base method
func (m *MockNode) GetBlockInfo(height uint64) (*api.ResultBlockInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockInfo", height)
	ret0, _ := ret[0].(*api.ResultBlockInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockInfo indicates an expected call of GetBlockInfo
func (mr *MockNodeMockRecorder) GetBlockInfo(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockInfo", reflect.TypeOf((*MockNode)(nil).GetBlockInfo), height)
}

// GetValidators mocks base method
func (m *MockNode) GetValidators(height uint64) ([]*api.ResultValidator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidators", height)
	ret0, _ := ret[0].([]*api.ResultValidator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidators indicates an expected call of GetValidators
func (mr *MockNodeMockRecorder) GetValidators(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidators", reflect.TypeOf((*MockNode)(nil).GetValidators), height)
}

// IsSyncing mocks base method
func (m *MockNode) IsSyncing() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSyncing")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsSyncing indicates an expected call of IsSyncing
func (mr *MockNodeMockRecorder) IsSyncing() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSyncing", reflect.TypeOf((*MockNode)(nil).IsSyncing))
}

// MockDHT is a mock of DHT interface
type MockDHT struct {
	ctrl     *gomock.Controller
	recorder *MockDHTMockRecorder
}

// MockDHTMockRecorder is the mock recorder for MockDHT
type MockDHTMockRecorder struct {
	mock *MockDHT
}

// NewMockDHT creates a new mock instance
func NewMockDHT(ctrl *gomock.Controller) *MockDHT {
	mock := &MockDHT{ctrl: ctrl}
	mock.recorder = &MockDHTMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDHT) EXPECT() *MockDHTMockRecorder {
	return m.recorder
}

// GetPeers mocks base method
func (m *MockDHT) GetPeers() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeers")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPeers indicates an expected call of GetPeers
func (mr *MockDHTMockRecorder) GetPeers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeers", reflect.TypeOf((*MockDHT)(nil).GetPeers))
}

// GetProviders mocks base method
func (m *MockDHT) GetProviders(key string) ([]*api.ResultDHTProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProviders", key)
	ret0, _ := ret[0].([]*api.ResultDHTProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProviders indicates an expected call of GetProviders
func (mr *MockDHTMockRecorder) GetProviders(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProviders", reflect.TypeOf((*MockDHT)(nil).GetProviders), key)
}

// Announce mocks base method
func (m *MockDHT) Announce(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Announce", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Announce indicates an expected call of Announce
func (mr *MockDHTMockRecorder) Announce(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Announce", reflect.TypeOf((*MockDHT)(nil).Announce), key)
}

// GetRepoObjectProviders mocks base method
func (m *MockDHT) GetRepoObjectProviders(hash string) ([]*api.ResultDHTProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepoObjectProviders", hash)
	ret0, _ := ret[0].([]*api.ResultDHTProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepoObjectProviders indicates an expected call of GetRepoObjectProviders
func (mr *MockDHTMockRecorder) GetRepoObjectProviders(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepoObjectProviders", reflect.TypeOf((*MockDHT)(nil).GetRepoObjectProviders), hash)
}

// Store mocks base method
func (m *MockDHT) Store(key, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store
func (mr *MockDHTMockRecorder) Store(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockDHT)(nil).Store), key, value)
}

// Lookup mocks base method
func (m *MockDHT) Lookup(key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Lookup", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Lookup indicates an expected call of Lookup
func (mr *MockDHTMockRecorder) Lookup(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lookup", reflect.TypeOf((*MockDHT)(nil).Lookup), key)
}

// MockPool is a mock of Pool interface
type MockPool struct {
	ctrl     *gomock.Controller
	recorder *MockPoolMockRecorder
}

// MockPoolMockRecorder is the mock recorder for MockPool
type MockPoolMockRecorder struct {
	mock *MockPool
}

// NewMockPool creates a new mock instance
func NewMockPool(ctrl *gomock.Controller) *MockPool {
	mock := &MockPool{ctrl: ctrl}
	mock.recorder = &MockPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPool) EXPECT() *MockPoolMockRecorder {
	return m.recorder
}

// GetSize mocks base method
func (m *MockPool) GetSize() (*api.ResultPoolSize, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSize")
	ret0, _ := ret[0].(*api.ResultPoolSize)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSize indicates an expected call of GetSize
func (mr *MockPoolMockRecorder) GetSize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSize", reflect.TypeOf((*MockPool)(nil).GetSize))
}

// GetPushPoolSize mocks base method
func (m *MockPool) GetPushPoolSize() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPushPoolSize")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPushPoolSize indicates an expected call of GetPushPoolSize
func (mr *MockPoolMockRecorder) GetPushPoolSize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPushPoolSize", reflect.TypeOf((*MockPool)(nil).GetPushPoolSize))
}

// MockRepo is a mock of Repo interface
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockRepo) Create(body *api.BodyCreateRepo) (*api.ResultCreateRepo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", body)
	ret0, _ := ret[0].(*api.ResultCreateRepo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockRepoMockRecorder) Create(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepo)(nil).Create), body)
}

// Get mocks base method
func (m *MockRepo) Get(name string, opts ...*api.GetRepoOpts) (*api.ResultRepository, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{name}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*api.ResultRepository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockRepoMockRecorder) Get(name interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{name}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepo)(nil).Get), varargs...)
}

// AddContributors mocks base method
func (m *MockRepo) AddContributors(body *api.BodyAddRepoContribs) (*api.ResultHash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddContributors", body)
	ret0, _ := ret[0].(*api.ResultHash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddContributors indicates an expected call of AddContributors
func (mr *MockRepoMockRecorder) AddContributors(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddContributors", reflect.TypeOf((*MockRepo)(nil).AddContributors), body)
}

// VoteProposal mocks base method
func (m *MockRepo) VoteProposal(body *api.BodyRepoVote) (*api.ResultHash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VoteProposal", body)
	ret0, _ := ret[0].(*api.ResultHash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VoteProposal indicates an expected call of VoteProposal
func (mr *MockRepoMockRecorder) VoteProposal(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VoteProposal", reflect.TypeOf((*MockRepo)(nil).VoteProposal), body)
}

// MockRPC is a mock of RPC interface
type MockRPC struct {
	ctrl     *gomock.Controller
	recorder *MockRPCMockRecorder
}

// MockRPCMockRecorder is the mock recorder for MockRPC
type MockRPCMockRecorder struct {
	mock *MockRPC
}

// NewMockRPC creates a new mock instance
func NewMockRPC(ctrl *gomock.Controller) *MockRPC {
	mock := &MockRPC{ctrl: ctrl}
	mock.recorder = &MockRPCMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRPC) EXPECT() *MockRPCMockRecorder {
	return m.recorder
}

// GetMethods mocks base method
func (m *MockRPC) GetMethods() ([]rpc.MethodInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMethods")
	ret0, _ := ret[0].([]rpc.MethodInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMethods indicates an expected call of GetMethods
func (mr *MockRPCMockRecorder) GetMethods() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMethods", reflect.TypeOf((*MockRPC)(nil).GetMethods))
}

// MockTx is a mock of Tx interface
type MockTx struct {
	ctrl     *gomock.Controller
	recorder *MockTxMockRecorder
}

// MockTxMockRecorder is the mock recorder for MockTx
type MockTxMockRecorder struct {
	mock *MockTx
}

// NewMockTx creates a new mock instance
func NewMockTx(ctrl *gomock.Controller) *MockTx {
	mock := &MockTx{ctrl: ctrl}
	mock.recorder = &MockTxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTx) EXPECT() *MockTxMockRecorder {
	return m.recorder
}

// Send mocks base method
func (m *MockTx) Send(data map[string]interface{}) (*api.ResultHash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", data)
	ret0, _ := ret[0].(*api.ResultHash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Send indicates an expected call of Send
func (mr *MockTxMockRecorder) Send(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockTx)(nil).Send), data)
}

// Get mocks base method
func (m *MockTx) Get(hash string) (*api.ResultTx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", hash)
	ret0, _ := ret[0].(*api.ResultTx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockTxMockRecorder) Get(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTx)(nil).Get), hash)
}

// MockUser is a mock of User interface
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockUser) Get(address string, blockHeight ...uint64) (*api.ResultAccount, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{address}
	for _, a := range blockHeight {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*api.ResultAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockUserMockRecorder) Get(address interface{}, blockHeight ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{address}, blockHeight...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUser)(nil).Get), varargs...)
}

// Send mocks base method
func (m *MockUser) Send(body *api.BodySendCoin) (*api.ResultHash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", body)
	ret0, _ := ret[0].(*api.ResultHash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Send indicates an expected call of Send
func (mr *MockUserMockRecorder) Send(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockUser)(nil).Send), body)
}

// GetNonce mocks base method
func (m *MockUser) GetNonce(address string, blockHeight ...uint64) (uint64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{address}
	for _, a := range blockHeight {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetNonce", varargs...)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNonce indicates an expected call of GetNonce
func (mr *MockUserMockRecorder) GetNonce(address interface{}, blockHeight ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{address}, blockHeight...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNonce", reflect.TypeOf((*MockUser)(nil).GetNonce), varargs...)
}

// GetKeys mocks base method
func (m *MockUser) GetKeys() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKeys")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKeys indicates an expected call of GetKeys
func (mr *MockUserMockRecorder) GetKeys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeys", reflect.TypeOf((*MockUser)(nil).GetKeys))
}

// GetBalance mocks base method
func (m *MockUser) GetBalance(address string, blockHeight ...uint64) (float64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{address}
	for _, a := range blockHeight {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetBalance", varargs...)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance
func (mr *MockUserMockRecorder) GetBalance(address interface{}, blockHeight ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{address}, blockHeight...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockUser)(nil).GetBalance), varargs...)
}

// GetStakedBalance mocks base method
func (m *MockUser) GetStakedBalance(address string, blockHeight ...uint64) (float64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{address}
	for _, a := range blockHeight {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetStakedBalance", varargs...)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStakedBalance indicates an expected call of GetStakedBalance
func (mr *MockUserMockRecorder) GetStakedBalance(address interface{}, blockHeight ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{address}, blockHeight...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStakedBalance", reflect.TypeOf((*MockUser)(nil).GetStakedBalance), varargs...)
}

// GetValidator mocks base method
func (m *MockUser) GetValidator(includePrivKey bool) (*api.ResultValidatorInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidator", includePrivKey)
	ret0, _ := ret[0].(*api.ResultValidatorInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidator indicates an expected call of GetValidator
func (mr *MockUserMockRecorder) GetValidator(includePrivKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidator", reflect.TypeOf((*MockUser)(nil).GetValidator), includePrivKey)
}

// GetPrivateKey mocks base method
func (m *MockUser) GetPrivateKey(address, passphrase string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPrivateKey", address, passphrase)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPrivateKey indicates an expected call of GetPrivateKey
func (mr *MockUserMockRecorder) GetPrivateKey(address, passphrase interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPrivateKey", reflect.TypeOf((*MockUser)(nil).GetPrivateKey), address, passphrase)
}

// GetPublicKey mocks base method
func (m *MockUser) GetPublicKey(address, passphrase string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublicKey", address, passphrase)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPublicKey indicates an expected call of GetPublicKey
func (mr *MockUserMockRecorder) GetPublicKey(address, passphrase interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublicKey", reflect.TypeOf((*MockUser)(nil).GetPublicKey), address, passphrase)
}

// SetCommission mocks base method
func (m *MockUser) SetCommission(body *api.BodySetCommission) (*api.ResultHash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCommission", body)
	ret0, _ := ret[0].(*api.ResultHash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetCommission indicates an expected call of SetCommission
func (mr *MockUserMockRecorder) SetCommission(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCommission", reflect.TypeOf((*MockUser)(nil).SetCommission), body)
}

// MockTicket is a mock of Ticket interface
type MockTicket struct {
	ctrl     *gomock.Controller
	recorder *MockTicketMockRecorder
}

// MockTicketMockRecorder is the mock recorder for MockTicket
type MockTicketMockRecorder struct {
	mock *MockTicket
}

// NewMockTicket creates a new mock instance
func NewMockTicket(ctrl *gomock.Controller) *MockTicket {
	mock := &MockTicket{ctrl: ctrl}
	mock.recorder = &MockTicketMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTicket) EXPECT() *MockTicketMockRecorder {
	return m.recorder
}

// Buy mocks base method
func (m *MockTicket) Buy(body *api.BodyBuyTicket) (*api.ResultHash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Buy", body)
	ret0, _ := ret[0].(*api.ResultHash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Buy indicates an expected call of Buy
func (mr *MockTicketMockRecorder) Buy(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Buy", reflect.TypeOf((*MockTicket)(nil).Buy), body)
}

// BuyHost mocks base method
func (m *MockTicket) BuyHost(body *api.BodyBuyTicket) (*api.ResultHash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuyHost", body)
	ret0, _ := ret[0].(*api.ResultHash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuyHost indicates an expected call of BuyHost
func (mr *MockTicketMockRecorder) BuyHost(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuyHost", reflect.TypeOf((*MockTicket)(nil).BuyHost), body)
}

// List mocks base method
func (m *MockTicket) List(body *api.BodyTicketQuery) ([]*api.ResultTicket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", body)
	ret0, _ := ret[0].([]*api.ResultTicket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockTicketMockRecorder) List(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTicket)(nil).List), body)
}

// ListHost mocks base method
func (m *MockTicket) ListHost(body *api.BodyTicketQuery) ([]*api.ResultTicket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListHost", body)
	ret0, _ := ret[0].([]*api.ResultTicket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListHost indicates an expected call of ListHost
func (mr *MockTicketMockRecorder) ListHost(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListHost", reflect.TypeOf((*MockTicket)(nil).ListHost), body)
}
