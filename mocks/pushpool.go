// Code generated by MockGen. DO NOT EDIT.
// Source: remote/push/types/objects.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	msgpack "github.com/vmihailenco/msgpack"
	types "gitlab.com/makeos/mosdef/remote/push/types"
	types0 "gitlab.com/makeos/mosdef/remote/types"
	util "gitlab.com/makeos/mosdef/util"
	reflect "reflect"
)

// MockEndorsement is a mock of Endorsement interface
type MockEndorsement struct {
	ctrl     *gomock.Controller
	recorder *MockEndorsementMockRecorder
}

// MockEndorsementMockRecorder is the mock recorder for MockEndorsement
type MockEndorsementMockRecorder struct {
	mock *MockEndorsement
}

// NewMockEndorsement creates a new mock instance
func NewMockEndorsement(ctrl *gomock.Controller) *MockEndorsement {
	mock := &MockEndorsement{ctrl: ctrl}
	mock.recorder = &MockEndorsementMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEndorsement) EXPECT() *MockEndorsementMockRecorder {
	return m.recorder
}

// ID mocks base method
func (m *MockEndorsement) ID() util.Bytes32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(util.Bytes32)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockEndorsementMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockEndorsement)(nil).ID))
}

// Bytes mocks base method
func (m *MockEndorsement) Bytes() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bytes")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Bytes indicates an expected call of Bytes
func (mr *MockEndorsementMockRecorder) Bytes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bytes", reflect.TypeOf((*MockEndorsement)(nil).Bytes))
}

// BytesAndID mocks base method
func (m *MockEndorsement) BytesAndID() ([]byte, util.Bytes32) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BytesAndID")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(util.Bytes32)
	return ret0, ret1
}

// BytesAndID indicates an expected call of BytesAndID
func (mr *MockEndorsementMockRecorder) BytesAndID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BytesAndID", reflect.TypeOf((*MockEndorsement)(nil).BytesAndID))
}

// MockPushPool is a mock of PushPool interface
type MockPushPool struct {
	ctrl     *gomock.Controller
	recorder *MockPushPoolMockRecorder
}

// MockPushPoolMockRecorder is the mock recorder for MockPushPool
type MockPushPoolMockRecorder struct {
	mock *MockPushPool
}

// NewMockPushPool creates a new mock instance
func NewMockPushPool(ctrl *gomock.Controller) *MockPushPool {
	mock := &MockPushPool{ctrl: ctrl}
	mock.recorder = &MockPushPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPushPool) EXPECT() *MockPushPoolMockRecorder {
	return m.recorder
}

// Add mocks base method
func (m *MockPushPool) Add(tx types.PushNote, noValidation ...bool) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{tx}
	for _, a := range noValidation {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Add", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add
func (mr *MockPushPoolMockRecorder) Add(tx interface{}, noValidation ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{tx}, noValidation...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockPushPool)(nil).Add), varargs...)
}

// Full mocks base method
func (m *MockPushPool) Full() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Full")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Full indicates an expected call of Full
func (mr *MockPushPoolMockRecorder) Full() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Full", reflect.TypeOf((*MockPushPool)(nil).Full))
}

// RepoHasPushNote mocks base method
func (m *MockPushPool) RepoHasPushNote(repo string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RepoHasPushNote", repo)
	ret0, _ := ret[0].(bool)
	return ret0
}

// RepoHasPushNote indicates an expected call of RepoHasPushNote
func (mr *MockPushPoolMockRecorder) RepoHasPushNote(repo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RepoHasPushNote", reflect.TypeOf((*MockPushPool)(nil).RepoHasPushNote), repo)
}

// Get mocks base method
func (m *MockPushPool) Get(noteID string) *types.Note {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", noteID)
	ret0, _ := ret[0].(*types.Note)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockPushPoolMockRecorder) Get(noteID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPushPool)(nil).Get), noteID)
}

// Len mocks base method
func (m *MockPushPool) Len() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Len")
	ret0, _ := ret[0].(int)
	return ret0
}

// Len indicates an expected call of Len
func (mr *MockPushPoolMockRecorder) Len() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Len", reflect.TypeOf((*MockPushPool)(nil).Len))
}

// Remove mocks base method
func (m *MockPushPool) Remove(pushNote types.PushNote) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Remove", pushNote)
}

// Remove indicates an expected call of Remove
func (mr *MockPushPoolMockRecorder) Remove(pushNote interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockPushPool)(nil).Remove), pushNote)
}

// MockPushNote is a mock of PushNote interface
type MockPushNote struct {
	ctrl     *gomock.Controller
	recorder *MockPushNoteMockRecorder
}

// MockPushNoteMockRecorder is the mock recorder for MockPushNote
type MockPushNoteMockRecorder struct {
	mock *MockPushNote
}

// NewMockPushNote creates a new mock instance
func NewMockPushNote(ctrl *gomock.Controller) *MockPushNote {
	mock := &MockPushNote{ctrl: ctrl}
	mock.recorder = &MockPushNoteMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPushNote) EXPECT() *MockPushNoteMockRecorder {
	return m.recorder
}

// GetTargetRepo mocks base method
func (m *MockPushNote) GetTargetRepo() types0.LocalRepo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTargetRepo")
	ret0, _ := ret[0].(types0.LocalRepo)
	return ret0
}

// GetTargetRepo indicates an expected call of GetTargetRepo
func (mr *MockPushNoteMockRecorder) GetTargetRepo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTargetRepo", reflect.TypeOf((*MockPushNote)(nil).GetTargetRepo))
}

// SetTargetRepo mocks base method
func (m *MockPushNote) SetTargetRepo(repo types0.LocalRepo) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTargetRepo", repo)
}

// SetTargetRepo indicates an expected call of SetTargetRepo
func (mr *MockPushNoteMockRecorder) SetTargetRepo(repo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTargetRepo", reflect.TypeOf((*MockPushNote)(nil).SetTargetRepo), repo)
}

// GetPusherKeyID mocks base method
func (m *MockPushNote) GetPusherKeyID() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPusherKeyID")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// GetPusherKeyID indicates an expected call of GetPusherKeyID
func (mr *MockPushNoteMockRecorder) GetPusherKeyID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPusherKeyID", reflect.TypeOf((*MockPushNote)(nil).GetPusherKeyID))
}

// GetPusherAddress mocks base method
func (m *MockPushNote) GetPusherAddress() util.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPusherAddress")
	ret0, _ := ret[0].(util.Address)
	return ret0
}

// GetPusherAddress indicates an expected call of GetPusherAddress
func (mr *MockPushNoteMockRecorder) GetPusherAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPusherAddress", reflect.TypeOf((*MockPushNote)(nil).GetPusherAddress))
}

// GetPusherAccountNonce mocks base method
func (m *MockPushNote) GetPusherAccountNonce() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPusherAccountNonce")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetPusherAccountNonce indicates an expected call of GetPusherAccountNonce
func (mr *MockPushNoteMockRecorder) GetPusherAccountNonce() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPusherAccountNonce", reflect.TypeOf((*MockPushNote)(nil).GetPusherAccountNonce))
}

// GetPusherKeyIDString mocks base method
func (m *MockPushNote) GetPusherKeyIDString() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPusherKeyIDString")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPusherKeyIDString indicates an expected call of GetPusherKeyIDString
func (mr *MockPushNoteMockRecorder) GetPusherKeyIDString() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPusherKeyIDString", reflect.TypeOf((*MockPushNote)(nil).GetPusherKeyIDString))
}

// EncodeMsgpack mocks base method
func (m *MockPushNote) EncodeMsgpack(enc *msgpack.Encoder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncodeMsgpack", enc)
	ret0, _ := ret[0].(error)
	return ret0
}

// EncodeMsgpack indicates an expected call of EncodeMsgpack
func (mr *MockPushNoteMockRecorder) EncodeMsgpack(enc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncodeMsgpack", reflect.TypeOf((*MockPushNote)(nil).EncodeMsgpack), enc)
}

// DecodeMsgpack mocks base method
func (m *MockPushNote) DecodeMsgpack(dec *msgpack.Decoder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecodeMsgpack", dec)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecodeMsgpack indicates an expected call of DecodeMsgpack
func (mr *MockPushNoteMockRecorder) DecodeMsgpack(dec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecodeMsgpack", reflect.TypeOf((*MockPushNote)(nil).DecodeMsgpack), dec)
}

// Bytes mocks base method
func (m *MockPushNote) Bytes(recompute ...bool) []byte {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range recompute {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Bytes", varargs...)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Bytes indicates an expected call of Bytes
func (mr *MockPushNoteMockRecorder) Bytes(recompute ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bytes", reflect.TypeOf((*MockPushNote)(nil).Bytes), recompute...)
}

// BytesNoCache mocks base method
func (m *MockPushNote) BytesNoCache() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BytesNoCache")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// BytesNoCache indicates an expected call of BytesNoCache
func (mr *MockPushNoteMockRecorder) BytesNoCache() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BytesNoCache", reflect.TypeOf((*MockPushNote)(nil).BytesNoCache))
}

// BytesNoSig mocks base method
func (m *MockPushNote) BytesNoSig() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BytesNoSig")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// BytesNoSig indicates an expected call of BytesNoSig
func (mr *MockPushNoteMockRecorder) BytesNoSig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BytesNoSig", reflect.TypeOf((*MockPushNote)(nil).BytesNoSig))
}

// GetLocalSize mocks base method
func (m *MockPushNote) GetLocalSize() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocalSize")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetLocalSize indicates an expected call of GetLocalSize
func (mr *MockPushNoteMockRecorder) GetLocalSize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocalSize", reflect.TypeOf((*MockPushNote)(nil).GetLocalSize))
}

// GetEcoSize mocks base method
func (m *MockPushNote) GetEcoSize() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEcoSize")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetEcoSize indicates an expected call of GetEcoSize
func (mr *MockPushNoteMockRecorder) GetEcoSize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEcoSize", reflect.TypeOf((*MockPushNote)(nil).GetEcoSize))
}

// GetCreatorPubKey mocks base method
func (m *MockPushNote) GetCreatorPubKey() util.Bytes32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCreatorPubKey")
	ret0, _ := ret[0].(util.Bytes32)
	return ret0
}

// GetCreatorPubKey indicates an expected call of GetCreatorPubKey
func (mr *MockPushNoteMockRecorder) GetCreatorPubKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCreatorPubKey", reflect.TypeOf((*MockPushNote)(nil).GetCreatorPubKey))
}

// GetNodeSignature mocks base method
func (m *MockPushNote) GetNodeSignature() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeSignature")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// GetNodeSignature indicates an expected call of GetNodeSignature
func (mr *MockPushNoteMockRecorder) GetNodeSignature() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeSignature", reflect.TypeOf((*MockPushNote)(nil).GetNodeSignature))
}

// GetRepoName mocks base method
func (m *MockPushNote) GetRepoName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepoName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetRepoName indicates an expected call of GetRepoName
func (mr *MockPushNoteMockRecorder) GetRepoName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepoName", reflect.TypeOf((*MockPushNote)(nil).GetRepoName))
}

// GetNamespace mocks base method
func (m *MockPushNote) GetNamespace() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNamespace")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetNamespace indicates an expected call of GetNamespace
func (mr *MockPushNoteMockRecorder) GetNamespace() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNamespace", reflect.TypeOf((*MockPushNote)(nil).GetNamespace))
}

// GetTimestamp mocks base method
func (m *MockPushNote) GetTimestamp() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTimestamp")
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetTimestamp indicates an expected call of GetTimestamp
func (mr *MockPushNoteMockRecorder) GetTimestamp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTimestamp", reflect.TypeOf((*MockPushNote)(nil).GetTimestamp))
}

// GetPushedReferences mocks base method
func (m *MockPushNote) GetPushedReferences() types.PushedReferences {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPushedReferences")
	ret0, _ := ret[0].(types.PushedReferences)
	return ret0
}

// GetPushedReferences indicates an expected call of GetPushedReferences
func (mr *MockPushNoteMockRecorder) GetPushedReferences() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPushedReferences", reflect.TypeOf((*MockPushNote)(nil).GetPushedReferences))
}

// Len mocks base method
func (m *MockPushNote) Len() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Len")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// Len indicates an expected call of Len
func (mr *MockPushNoteMockRecorder) Len() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Len", reflect.TypeOf((*MockPushNote)(nil).Len))
}

// ID mocks base method
func (m *MockPushNote) ID(recompute ...bool) util.Bytes32 {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range recompute {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ID", varargs...)
	ret0, _ := ret[0].(util.Bytes32)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockPushNoteMockRecorder) ID(recompute ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockPushNote)(nil).ID), recompute...)
}

// BytesAndID mocks base method
func (m *MockPushNote) BytesAndID(recompute ...bool) ([]byte, util.Bytes32) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range recompute {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BytesAndID", varargs...)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(util.Bytes32)
	return ret0, ret1
}

// BytesAndID indicates an expected call of BytesAndID
func (mr *MockPushNoteMockRecorder) BytesAndID(recompute ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BytesAndID", reflect.TypeOf((*MockPushNote)(nil).BytesAndID), recompute...)
}

// TxSize mocks base method
func (m *MockPushNote) TxSize() uint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxSize")
	ret0, _ := ret[0].(uint)
	return ret0
}

// TxSize indicates an expected call of TxSize
func (mr *MockPushNoteMockRecorder) TxSize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxSize", reflect.TypeOf((*MockPushNote)(nil).TxSize))
}

// SizeForFeeCal mocks base method
func (m *MockPushNote) SizeForFeeCal() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SizeForFeeCal")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// SizeForFeeCal indicates an expected call of SizeForFeeCal
func (mr *MockPushNoteMockRecorder) SizeForFeeCal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SizeForFeeCal", reflect.TypeOf((*MockPushNote)(nil).SizeForFeeCal))
}

// GetSize mocks base method
func (m *MockPushNote) GetSize() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSize")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetSize indicates an expected call of GetSize
func (mr *MockPushNoteMockRecorder) GetSize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSize", reflect.TypeOf((*MockPushNote)(nil).GetSize))
}

// GetFee mocks base method
func (m *MockPushNote) GetFee() util.String {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFee")
	ret0, _ := ret[0].(util.String)
	return ret0
}

// GetFee indicates an expected call of GetFee
func (mr *MockPushNoteMockRecorder) GetFee() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFee", reflect.TypeOf((*MockPushNote)(nil).GetFee))
}

// GetValue mocks base method
func (m *MockPushNote) GetValue() util.String {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValue")
	ret0, _ := ret[0].(util.String)
	return ret0
}

// GetValue indicates an expected call of GetValue
func (mr *MockPushNoteMockRecorder) GetValue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValue", reflect.TypeOf((*MockPushNote)(nil).GetValue))
}

// IsFromPeer mocks base method
func (m *MockPushNote) IsFromPeer() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsFromPeer")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsFromPeer indicates an expected call of IsFromPeer
func (mr *MockPushNoteMockRecorder) IsFromPeer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsFromPeer", reflect.TypeOf((*MockPushNote)(nil).IsFromPeer))
}
