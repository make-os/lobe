// Code generated by MockGen. DO NOT EDIT.
// Source: dht/streamer/types/types.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	network "github.com/libp2p/go-libp2p-core/network"
	types "gitlab.com/makeos/mosdef/dht/streamer/types"
	io "gitlab.com/makeos/mosdef/util/io"
	object "gopkg.in/src-d/go-git.v4/plumbing/object"
	reflect "reflect"
)

// MockObjectStreamer is a mock of ObjectStreamer interface
type MockObjectStreamer struct {
	ctrl     *gomock.Controller
	recorder *MockObjectStreamerMockRecorder
}

// MockObjectStreamerMockRecorder is the mock recorder for MockObjectStreamer
type MockObjectStreamerMockRecorder struct {
	mock *MockObjectStreamer
}

// NewMockObjectStreamer creates a new mock instance
func NewMockObjectStreamer(ctrl *gomock.Controller) *MockObjectStreamer {
	mock := &MockObjectStreamer{ctrl: ctrl}
	mock.recorder = &MockObjectStreamerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockObjectStreamer) EXPECT() *MockObjectStreamerMockRecorder {
	return m.recorder
}

// Announce mocks base method
func (m *MockObjectStreamer) Announce(hash []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Announce", hash)
	ret0, _ := ret[0].(error)
	return ret0
}

// Announce indicates an expected call of Announce
func (mr *MockObjectStreamerMockRecorder) Announce(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Announce", reflect.TypeOf((*MockObjectStreamer)(nil).Announce), hash)
}

// GetCommit mocks base method
func (m *MockObjectStreamer) GetCommit(ctx context.Context, repo string, hash []byte) (io.ReadSeekerCloser, *object.Commit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommit", ctx, repo, hash)
	ret0, _ := ret[0].(io.ReadSeekerCloser)
	ret1, _ := ret[1].(*object.Commit)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCommit indicates an expected call of GetCommit
func (mr *MockObjectStreamerMockRecorder) GetCommit(ctx, repo, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommit", reflect.TypeOf((*MockObjectStreamer)(nil).GetCommit), ctx, repo, hash)
}

// GetCommitWithAncestors mocks base method
func (m *MockObjectStreamer) GetCommitWithAncestors(ctx context.Context, args types.GetAncestorArgs) ([]io.ReadSeekerCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommitWithAncestors", ctx, args)
	ret0, _ := ret[0].([]io.ReadSeekerCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommitWithAncestors indicates an expected call of GetCommitWithAncestors
func (mr *MockObjectStreamerMockRecorder) GetCommitWithAncestors(ctx, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommitWithAncestors", reflect.TypeOf((*MockObjectStreamer)(nil).GetCommitWithAncestors), ctx, args)
}

// GetTaggedCommitWithAncestors mocks base method
func (m *MockObjectStreamer) GetTaggedCommitWithAncestors(ctx context.Context, args types.GetAncestorArgs) ([]io.ReadSeekerCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaggedCommitWithAncestors", ctx, args)
	ret0, _ := ret[0].([]io.ReadSeekerCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaggedCommitWithAncestors indicates an expected call of GetTaggedCommitWithAncestors
func (mr *MockObjectStreamerMockRecorder) GetTaggedCommitWithAncestors(ctx, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaggedCommitWithAncestors", reflect.TypeOf((*MockObjectStreamer)(nil).GetTaggedCommitWithAncestors), ctx, args)
}

// GetTag mocks base method
func (m *MockObjectStreamer) GetTag(ctx context.Context, repo string, hash []byte) (io.ReadSeekerCloser, *object.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTag", ctx, repo, hash)
	ret0, _ := ret[0].(io.ReadSeekerCloser)
	ret1, _ := ret[1].(*object.Tag)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetTag indicates an expected call of GetTag
func (mr *MockObjectStreamerMockRecorder) GetTag(ctx, repo, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTag", reflect.TypeOf((*MockObjectStreamer)(nil).GetTag), ctx, repo, hash)
}

// OnRequest mocks base method
func (m *MockObjectStreamer) OnRequest(s network.Stream) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnRequest", s)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OnRequest indicates an expected call of OnRequest
func (mr *MockObjectStreamerMockRecorder) OnRequest(s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnRequest", reflect.TypeOf((*MockObjectStreamer)(nil).OnRequest), s)
}
