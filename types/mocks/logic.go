// Code generated by MockGen. DO NOT EDIT.
// Source: logic.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	crypto "github.com/makeos/mosdef/crypto"
	rand "github.com/makeos/mosdef/crypto/rand"
	storage "github.com/makeos/mosdef/storage"
	types "github.com/makeos/mosdef/types"
	util "github.com/makeos/mosdef/util"
	types0 "github.com/tendermint/tendermint/abci/types"
	reflect "reflect"
)

// MockSystemKeeper is a mock of SystemKeeper interface
type MockSystemKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockSystemKeeperMockRecorder
}

// MockSystemKeeperMockRecorder is the mock recorder for MockSystemKeeper
type MockSystemKeeperMockRecorder struct {
	mock *MockSystemKeeper
}

// NewMockSystemKeeper creates a new mock instance
func NewMockSystemKeeper(ctrl *gomock.Controller) *MockSystemKeeper {
	mock := &MockSystemKeeper{ctrl: ctrl}
	mock.recorder = &MockSystemKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSystemKeeper) EXPECT() *MockSystemKeeperMockRecorder {
	return m.recorder
}

// SaveBlockInfo mocks base method
func (m *MockSystemKeeper) SaveBlockInfo(info *types.BlockInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBlockInfo", info)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveBlockInfo indicates an expected call of SaveBlockInfo
func (mr *MockSystemKeeperMockRecorder) SaveBlockInfo(info interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBlockInfo", reflect.TypeOf((*MockSystemKeeper)(nil).SaveBlockInfo), info)
}

// GetLastBlockInfo mocks base method
func (m *MockSystemKeeper) GetLastBlockInfo() (*types.BlockInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastBlockInfo")
	ret0, _ := ret[0].(*types.BlockInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastBlockInfo indicates an expected call of GetLastBlockInfo
func (mr *MockSystemKeeperMockRecorder) GetLastBlockInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastBlockInfo", reflect.TypeOf((*MockSystemKeeper)(nil).GetLastBlockInfo))
}

// GetBlockInfo mocks base method
func (m *MockSystemKeeper) GetBlockInfo(height int64) (*types.BlockInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockInfo", height)
	ret0, _ := ret[0].(*types.BlockInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockInfo indicates an expected call of GetBlockInfo
func (mr *MockSystemKeeperMockRecorder) GetBlockInfo(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockInfo", reflect.TypeOf((*MockSystemKeeper)(nil).GetBlockInfo), height)
}

// MarkAsMatured mocks base method
func (m *MockSystemKeeper) MarkAsMatured(maturityHeight uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkAsMatured", maturityHeight)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkAsMatured indicates an expected call of MarkAsMatured
func (mr *MockSystemKeeperMockRecorder) MarkAsMatured(maturityHeight interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkAsMatured", reflect.TypeOf((*MockSystemKeeper)(nil).MarkAsMatured), maturityHeight)
}

// GetNetMaturityHeight mocks base method
func (m *MockSystemKeeper) GetNetMaturityHeight() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetMaturityHeight")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNetMaturityHeight indicates an expected call of GetNetMaturityHeight
func (mr *MockSystemKeeperMockRecorder) GetNetMaturityHeight() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetMaturityHeight", reflect.TypeOf((*MockSystemKeeper)(nil).GetNetMaturityHeight))
}

// IsMarkedAsMature mocks base method
func (m *MockSystemKeeper) IsMarkedAsMature() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsMarkedAsMature")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsMarkedAsMature indicates an expected call of IsMarkedAsMature
func (mr *MockSystemKeeperMockRecorder) IsMarkedAsMature() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsMarkedAsMature", reflect.TypeOf((*MockSystemKeeper)(nil).IsMarkedAsMature))
}

// SetHighestDrandRound mocks base method
func (m *MockSystemKeeper) SetHighestDrandRound(r uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHighestDrandRound", r)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHighestDrandRound indicates an expected call of SetHighestDrandRound
func (mr *MockSystemKeeperMockRecorder) SetHighestDrandRound(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHighestDrandRound", reflect.TypeOf((*MockSystemKeeper)(nil).SetHighestDrandRound), r)
}

// GetHighestDrandRound mocks base method
func (m *MockSystemKeeper) GetHighestDrandRound() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHighestDrandRound")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHighestDrandRound indicates an expected call of GetHighestDrandRound
func (mr *MockSystemKeeperMockRecorder) GetHighestDrandRound() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHighestDrandRound", reflect.TypeOf((*MockSystemKeeper)(nil).GetHighestDrandRound))
}

// GetSecrets mocks base method
func (m *MockSystemKeeper) GetSecrets(from, limit, skip int64) ([][]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecrets", from, limit, skip)
	ret0, _ := ret[0].([][]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecrets indicates an expected call of GetSecrets
func (mr *MockSystemKeeperMockRecorder) GetSecrets(from, limit, skip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecrets", reflect.TypeOf((*MockSystemKeeper)(nil).GetSecrets), from, limit, skip)
}

// MockTxKeeper is a mock of TxKeeper interface
type MockTxKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockTxKeeperMockRecorder
}

// MockTxKeeperMockRecorder is the mock recorder for MockTxKeeper
type MockTxKeeperMockRecorder struct {
	mock *MockTxKeeper
}

// NewMockTxKeeper creates a new mock instance
func NewMockTxKeeper(ctrl *gomock.Controller) *MockTxKeeper {
	mock := &MockTxKeeper{ctrl: ctrl}
	mock.recorder = &MockTxKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTxKeeper) EXPECT() *MockTxKeeperMockRecorder {
	return m.recorder
}

// Index mocks base method
func (m *MockTxKeeper) Index(tx types.Tx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Index", tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Index indicates an expected call of Index
func (mr *MockTxKeeperMockRecorder) Index(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Index", reflect.TypeOf((*MockTxKeeper)(nil).Index), tx)
}

// GetTx mocks base method
func (m *MockTxKeeper) GetTx(hash string) (types.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTx", hash)
	ret0, _ := ret[0].(types.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTx indicates an expected call of GetTx
func (mr *MockTxKeeperMockRecorder) GetTx(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTx", reflect.TypeOf((*MockTxKeeper)(nil).GetTx), hash)
}

// MockAccountKeeper is a mock of AccountKeeper interface
type MockAccountKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockAccountKeeperMockRecorder
}

// MockAccountKeeperMockRecorder is the mock recorder for MockAccountKeeper
type MockAccountKeeperMockRecorder struct {
	mock *MockAccountKeeper
}

// NewMockAccountKeeper creates a new mock instance
func NewMockAccountKeeper(ctrl *gomock.Controller) *MockAccountKeeper {
	mock := &MockAccountKeeper{ctrl: ctrl}
	mock.recorder = &MockAccountKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAccountKeeper) EXPECT() *MockAccountKeeperMockRecorder {
	return m.recorder
}

// GetAccount mocks base method
func (m *MockAccountKeeper) GetAccount(address util.String, blockNum ...int64) *types.Account {
	m.ctrl.T.Helper()
	varargs := []interface{}{address}
	for _, a := range blockNum {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAccount", varargs...)
	ret0, _ := ret[0].(*types.Account)
	return ret0
}

// GetAccount indicates an expected call of GetAccount
func (mr *MockAccountKeeperMockRecorder) GetAccount(address interface{}, blockNum ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{address}, blockNum...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockAccountKeeper)(nil).GetAccount), varargs...)
}

// Update mocks base method
func (m *MockAccountKeeper) Update(address util.String, upd *types.Account) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Update", address, upd)
}

// Update indicates an expected call of Update
func (mr *MockAccountKeeperMockRecorder) Update(address, upd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAccountKeeper)(nil).Update), address, upd)
}

// MockLogic is a mock of Logic interface
type MockLogic struct {
	ctrl     *gomock.Controller
	recorder *MockLogicMockRecorder
}

// MockLogicMockRecorder is the mock recorder for MockLogic
type MockLogicMockRecorder struct {
	mock *MockLogic
}

// NewMockLogic creates a new mock instance
func NewMockLogic(ctrl *gomock.Controller) *MockLogic {
	mock := &MockLogic{ctrl: ctrl}
	mock.recorder = &MockLogicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogic) EXPECT() *MockLogicMockRecorder {
	return m.recorder
}

// SysKeeper mocks base method
func (m *MockLogic) SysKeeper() types.SystemKeeper {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SysKeeper")
	ret0, _ := ret[0].(types.SystemKeeper)
	return ret0
}

// SysKeeper indicates an expected call of SysKeeper
func (mr *MockLogicMockRecorder) SysKeeper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SysKeeper", reflect.TypeOf((*MockLogic)(nil).SysKeeper))
}

// AccountKeeper mocks base method
func (m *MockLogic) AccountKeeper() types.AccountKeeper {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountKeeper")
	ret0, _ := ret[0].(types.AccountKeeper)
	return ret0
}

// AccountKeeper indicates an expected call of AccountKeeper
func (mr *MockLogicMockRecorder) AccountKeeper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountKeeper", reflect.TypeOf((*MockLogic)(nil).AccountKeeper))
}

// ValidatorKeeper mocks base method
func (m *MockLogic) ValidatorKeeper() types.ValidatorKeeper {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidatorKeeper")
	ret0, _ := ret[0].(types.ValidatorKeeper)
	return ret0
}

// ValidatorKeeper indicates an expected call of ValidatorKeeper
func (mr *MockLogicMockRecorder) ValidatorKeeper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidatorKeeper", reflect.TypeOf((*MockLogic)(nil).ValidatorKeeper))
}

// TxKeeper mocks base method
func (m *MockLogic) TxKeeper() types.TxKeeper {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxKeeper")
	ret0, _ := ret[0].(types.TxKeeper)
	return ret0
}

// TxKeeper indicates an expected call of TxKeeper
func (mr *MockLogicMockRecorder) TxKeeper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxKeeper", reflect.TypeOf((*MockLogic)(nil).TxKeeper))
}

// GetTicketManager mocks base method
func (m *MockLogic) GetTicketManager() types.TicketManager {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTicketManager")
	ret0, _ := ret[0].(types.TicketManager)
	return ret0
}

// GetTicketManager indicates an expected call of GetTicketManager
func (mr *MockLogicMockRecorder) GetTicketManager() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTicketManager", reflect.TypeOf((*MockLogic)(nil).GetTicketManager))
}

// Tx mocks base method
func (m *MockLogic) Tx() types.TxLogic {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tx")
	ret0, _ := ret[0].(types.TxLogic)
	return ret0
}

// Tx indicates an expected call of Tx
func (mr *MockLogicMockRecorder) Tx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tx", reflect.TypeOf((*MockLogic)(nil).Tx))
}

// Sys mocks base method
func (m *MockLogic) Sys() types.SysLogic {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sys")
	ret0, _ := ret[0].(types.SysLogic)
	return ret0
}

// Sys indicates an expected call of Sys
func (mr *MockLogicMockRecorder) Sys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sys", reflect.TypeOf((*MockLogic)(nil).Sys))
}

// Validator mocks base method
func (m *MockLogic) Validator() types.ValidatorLogic {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validator")
	ret0, _ := ret[0].(types.ValidatorLogic)
	return ret0
}

// Validator indicates an expected call of Validator
func (mr *MockLogicMockRecorder) Validator() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validator", reflect.TypeOf((*MockLogic)(nil).Validator))
}

// DB mocks base method
func (m *MockLogic) DB() storage.Engine {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DB")
	ret0, _ := ret[0].(storage.Engine)
	return ret0
}

// DB indicates an expected call of DB
func (mr *MockLogicMockRecorder) DB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DB", reflect.TypeOf((*MockLogic)(nil).DB))
}

// StateTree mocks base method
func (m *MockLogic) StateTree() types.Tree {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StateTree")
	ret0, _ := ret[0].(types.Tree)
	return ret0
}

// StateTree indicates an expected call of StateTree
func (mr *MockLogicMockRecorder) StateTree() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StateTree", reflect.TypeOf((*MockLogic)(nil).StateTree))
}

// WriteGenesisState mocks base method
func (m *MockLogic) WriteGenesisState() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteGenesisState")
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteGenesisState indicates an expected call of WriteGenesisState
func (mr *MockLogicMockRecorder) WriteGenesisState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteGenesisState", reflect.TypeOf((*MockLogic)(nil).WriteGenesisState))
}

// SetTicketManager mocks base method
func (m *MockLogic) SetTicketManager(tm types.TicketManager) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTicketManager", tm)
}

// SetTicketManager indicates an expected call of SetTicketManager
func (mr *MockLogicMockRecorder) SetTicketManager(tm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTicketManager", reflect.TypeOf((*MockLogic)(nil).SetTicketManager), tm)
}

// GetDRand mocks base method
func (m *MockLogic) GetDRand() rand.DRander {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDRand")
	ret0, _ := ret[0].(rand.DRander)
	return ret0
}

// GetDRand indicates an expected call of GetDRand
func (mr *MockLogicMockRecorder) GetDRand() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDRand", reflect.TypeOf((*MockLogic)(nil).GetDRand))
}

// MockKeepers is a mock of Keepers interface
type MockKeepers struct {
	ctrl     *gomock.Controller
	recorder *MockKeepersMockRecorder
}

// MockKeepersMockRecorder is the mock recorder for MockKeepers
type MockKeepersMockRecorder struct {
	mock *MockKeepers
}

// NewMockKeepers creates a new mock instance
func NewMockKeepers(ctrl *gomock.Controller) *MockKeepers {
	mock := &MockKeepers{ctrl: ctrl}
	mock.recorder = &MockKeepersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKeepers) EXPECT() *MockKeepersMockRecorder {
	return m.recorder
}

// SysKeeper mocks base method
func (m *MockKeepers) SysKeeper() types.SystemKeeper {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SysKeeper")
	ret0, _ := ret[0].(types.SystemKeeper)
	return ret0
}

// SysKeeper indicates an expected call of SysKeeper
func (mr *MockKeepersMockRecorder) SysKeeper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SysKeeper", reflect.TypeOf((*MockKeepers)(nil).SysKeeper))
}

// AccountKeeper mocks base method
func (m *MockKeepers) AccountKeeper() types.AccountKeeper {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountKeeper")
	ret0, _ := ret[0].(types.AccountKeeper)
	return ret0
}

// AccountKeeper indicates an expected call of AccountKeeper
func (mr *MockKeepersMockRecorder) AccountKeeper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountKeeper", reflect.TypeOf((*MockKeepers)(nil).AccountKeeper))
}

// ValidatorKeeper mocks base method
func (m *MockKeepers) ValidatorKeeper() types.ValidatorKeeper {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidatorKeeper")
	ret0, _ := ret[0].(types.ValidatorKeeper)
	return ret0
}

// ValidatorKeeper indicates an expected call of ValidatorKeeper
func (mr *MockKeepersMockRecorder) ValidatorKeeper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidatorKeeper", reflect.TypeOf((*MockKeepers)(nil).ValidatorKeeper))
}

// TxKeeper mocks base method
func (m *MockKeepers) TxKeeper() types.TxKeeper {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxKeeper")
	ret0, _ := ret[0].(types.TxKeeper)
	return ret0
}

// TxKeeper indicates an expected call of TxKeeper
func (mr *MockKeepersMockRecorder) TxKeeper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxKeeper", reflect.TypeOf((*MockKeepers)(nil).TxKeeper))
}

// GetTicketManager mocks base method
func (m *MockKeepers) GetTicketManager() types.TicketManager {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTicketManager")
	ret0, _ := ret[0].(types.TicketManager)
	return ret0
}

// GetTicketManager indicates an expected call of GetTicketManager
func (mr *MockKeepersMockRecorder) GetTicketManager() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTicketManager", reflect.TypeOf((*MockKeepers)(nil).GetTicketManager))
}

// MockLogicCommon is a mock of LogicCommon interface
type MockLogicCommon struct {
	ctrl     *gomock.Controller
	recorder *MockLogicCommonMockRecorder
}

// MockLogicCommonMockRecorder is the mock recorder for MockLogicCommon
type MockLogicCommonMockRecorder struct {
	mock *MockLogicCommon
}

// NewMockLogicCommon creates a new mock instance
func NewMockLogicCommon(ctrl *gomock.Controller) *MockLogicCommon {
	mock := &MockLogicCommon{ctrl: ctrl}
	mock.recorder = &MockLogicCommonMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogicCommon) EXPECT() *MockLogicCommonMockRecorder {
	return m.recorder
}

// MockValidatorKeeper is a mock of ValidatorKeeper interface
type MockValidatorKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorKeeperMockRecorder
}

// MockValidatorKeeperMockRecorder is the mock recorder for MockValidatorKeeper
type MockValidatorKeeperMockRecorder struct {
	mock *MockValidatorKeeper
}

// NewMockValidatorKeeper creates a new mock instance
func NewMockValidatorKeeper(ctrl *gomock.Controller) *MockValidatorKeeper {
	mock := &MockValidatorKeeper{ctrl: ctrl}
	mock.recorder = &MockValidatorKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockValidatorKeeper) EXPECT() *MockValidatorKeeperMockRecorder {
	return m.recorder
}

// GetByHeight mocks base method
func (m *MockValidatorKeeper) GetByHeight(height int64) (types.BlockValidators, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByHeight", height)
	ret0, _ := ret[0].(types.BlockValidators)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByHeight indicates an expected call of GetByHeight
func (mr *MockValidatorKeeperMockRecorder) GetByHeight(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByHeight", reflect.TypeOf((*MockValidatorKeeper)(nil).GetByHeight), height)
}

// Index mocks base method
func (m *MockValidatorKeeper) Index(height int64, validators []*types.Validator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Index", height, validators)
	ret0, _ := ret[0].(error)
	return ret0
}

// Index indicates an expected call of Index
func (mr *MockValidatorKeeperMockRecorder) Index(height, validators interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Index", reflect.TypeOf((*MockValidatorKeeper)(nil).Index), height, validators)
}

// MockValidatorLogic is a mock of ValidatorLogic interface
type MockValidatorLogic struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorLogicMockRecorder
}

// MockValidatorLogicMockRecorder is the mock recorder for MockValidatorLogic
type MockValidatorLogicMockRecorder struct {
	mock *MockValidatorLogic
}

// NewMockValidatorLogic creates a new mock instance
func NewMockValidatorLogic(ctrl *gomock.Controller) *MockValidatorLogic {
	mock := &MockValidatorLogic{ctrl: ctrl}
	mock.recorder = &MockValidatorLogicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockValidatorLogic) EXPECT() *MockValidatorLogicMockRecorder {
	return m.recorder
}

// Index mocks base method
func (m *MockValidatorLogic) Index(height int64, valUpdates []types0.ValidatorUpdate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Index", height, valUpdates)
	ret0, _ := ret[0].(error)
	return ret0
}

// Index indicates an expected call of Index
func (mr *MockValidatorLogicMockRecorder) Index(height, valUpdates interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Index", reflect.TypeOf((*MockValidatorLogic)(nil).Index), height, valUpdates)
}

// MockTxLogic is a mock of TxLogic interface
type MockTxLogic struct {
	ctrl     *gomock.Controller
	recorder *MockTxLogicMockRecorder
}

// MockTxLogicMockRecorder is the mock recorder for MockTxLogic
type MockTxLogicMockRecorder struct {
	mock *MockTxLogic
}

// NewMockTxLogic creates a new mock instance
func NewMockTxLogic(ctrl *gomock.Controller) *MockTxLogic {
	mock := &MockTxLogic{ctrl: ctrl}
	mock.recorder = &MockTxLogicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTxLogic) EXPECT() *MockTxLogicMockRecorder {
	return m.recorder
}

// PrepareExec mocks base method
func (m *MockTxLogic) PrepareExec(req types0.RequestDeliverTx) types0.ResponseDeliverTx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareExec", req)
	ret0, _ := ret[0].(types0.ResponseDeliverTx)
	return ret0
}

// PrepareExec indicates an expected call of PrepareExec
func (mr *MockTxLogicMockRecorder) PrepareExec(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareExec", reflect.TypeOf((*MockTxLogic)(nil).PrepareExec), req)
}

// Exec mocks base method
func (m *MockTxLogic) Exec(tx *types.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec", tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Exec indicates an expected call of Exec
func (mr *MockTxLogicMockRecorder) Exec(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockTxLogic)(nil).Exec), tx)
}

// CanTransferCoin mocks base method
func (m *MockTxLogic) CanTransferCoin(txType int, senderPubKey *crypto.PubKey, recipientAddr, value, fee util.String, nonce uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanTransferCoin", txType, senderPubKey, recipientAddr, value, fee, nonce)
	ret0, _ := ret[0].(error)
	return ret0
}

// CanTransferCoin indicates an expected call of CanTransferCoin
func (mr *MockTxLogicMockRecorder) CanTransferCoin(txType, senderPubKey, recipientAddr, value, fee, nonce interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanTransferCoin", reflect.TypeOf((*MockTxLogic)(nil).CanTransferCoin), txType, senderPubKey, recipientAddr, value, fee, nonce)
}

// MockSysLogic is a mock of SysLogic interface
type MockSysLogic struct {
	ctrl     *gomock.Controller
	recorder *MockSysLogicMockRecorder
}

// MockSysLogicMockRecorder is the mock recorder for MockSysLogic
type MockSysLogicMockRecorder struct {
	mock *MockSysLogic
}

// NewMockSysLogic creates a new mock instance
func NewMockSysLogic(ctrl *gomock.Controller) *MockSysLogic {
	mock := &MockSysLogic{ctrl: ctrl}
	mock.recorder = &MockSysLogicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSysLogic) EXPECT() *MockSysLogicMockRecorder {
	return m.recorder
}

// GetCurValidatorTicketPrice mocks base method
func (m *MockSysLogic) GetCurValidatorTicketPrice() float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurValidatorTicketPrice")
	ret0, _ := ret[0].(float64)
	return ret0
}

// GetCurValidatorTicketPrice indicates an expected call of GetCurValidatorTicketPrice
func (mr *MockSysLogicMockRecorder) GetCurValidatorTicketPrice() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurValidatorTicketPrice", reflect.TypeOf((*MockSysLogic)(nil).GetCurValidatorTicketPrice))
}

// CheckSetNetMaturity mocks base method
func (m *MockSysLogic) CheckSetNetMaturity() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckSetNetMaturity")
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckSetNetMaturity indicates an expected call of CheckSetNetMaturity
func (mr *MockSysLogicMockRecorder) CheckSetNetMaturity() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSetNetMaturity", reflect.TypeOf((*MockSysLogic)(nil).CheckSetNetMaturity))
}

// GetEpoch mocks base method
func (m *MockSysLogic) GetEpoch(curBlockHeight uint64) (int, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEpoch", curBlockHeight)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// GetEpoch indicates an expected call of GetEpoch
func (mr *MockSysLogicMockRecorder) GetEpoch(curBlockHeight interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEpoch", reflect.TypeOf((*MockSysLogic)(nil).GetEpoch), curBlockHeight)
}

// GetCurretEpochSecretTx mocks base method
func (m *MockSysLogic) GetCurretEpochSecretTx() (types.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurretEpochSecretTx")
	ret0, _ := ret[0].(types.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurretEpochSecretTx indicates an expected call of GetCurretEpochSecretTx
func (mr *MockSysLogicMockRecorder) GetCurretEpochSecretTx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurretEpochSecretTx", reflect.TypeOf((*MockSysLogic)(nil).GetCurretEpochSecretTx))
}

// MakeSecret mocks base method
func (m *MockSysLogic) MakeSecret(height int64) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeSecret", height)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MakeSecret indicates an expected call of MakeSecret
func (mr *MockSysLogicMockRecorder) MakeSecret(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeSecret", reflect.TypeOf((*MockSysLogic)(nil).MakeSecret), height)
}
