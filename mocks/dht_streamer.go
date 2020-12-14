// Code generated by MockGen. DO NOT EDIT.
// Source: dht/streamer.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	network "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/make-os/kit/dht"
	io "github.com/make-os/kit/util/io"
	object "gopkg.in/src-d/go-git.v4/plumbing/object"
	reflect "reflect"
)

// MockStreamer is a mock of Streamer interface
type MockStreamer struct {
	ctrl     *gomock.Controller
	recorder *MockStreamerMockRecorder
}

// MockStreamerMockRecorder is the mock recorder for MockStreamer
type MockStreamerMockRecorder struct {
	mock *MockStreamer
}

// NewMockStreamer creates a new mock instance
func NewMockStreamer(ctrl *gomock.Controller) *MockStreamer {
	mock := &MockStreamer{ctrl: ctrl}
	mock.recorder = &MockStreamerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStreamer) EXPECT() *MockStreamerMockRecorder {
	return m.recorder
}

// GetCommit mocks base method
func (m *MockStreamer) GetCommit(ctx context.Context, repo string, hash []byte) (io.ReadSeekerCloser, *object.Commit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommit", ctx, repo, hash)
	ret0, _ := ret[0].(io.ReadSeekerCloser)
	ret1, _ := ret[1].(*object.Commit)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCommit indicates an expected call of GetCommit
func (mr *MockStreamerMockRecorder) GetCommit(ctx, repo, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommit", reflect.TypeOf((*MockStreamer)(nil).GetCommit), ctx, repo, hash)
}

// GetCommitWithAncestors mocks base method
func (m *MockStreamer) GetCommitWithAncestors(ctx context.Context, args dht.GetAncestorArgs) ([]io.ReadSeekerCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommitWithAncestors", ctx, args)
	ret0, _ := ret[0].([]io.ReadSeekerCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommitWithAncestors indicates an expected call of GetCommitWithAncestors
func (mr *MockStreamerMockRecorder) GetCommitWithAncestors(ctx, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommitWithAncestors", reflect.TypeOf((*MockStreamer)(nil).GetCommitWithAncestors), ctx, args)
}

// GetTaggedCommitWithAncestors mocks base method
func (m *MockStreamer) GetTaggedCommitWithAncestors(ctx context.Context, args dht.GetAncestorArgs) ([]io.ReadSeekerCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaggedCommitWithAncestors", ctx, args)
	ret0, _ := ret[0].([]io.ReadSeekerCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaggedCommitWithAncestors indicates an expected call of GetTaggedCommitWithAncestors
func (mr *MockStreamerMockRecorder) GetTaggedCommitWithAncestors(ctx, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaggedCommitWithAncestors", reflect.TypeOf((*MockStreamer)(nil).GetTaggedCommitWithAncestors), ctx, args)
}

// GetTag mocks base method
func (m *MockStreamer) GetTag(ctx context.Context, repo string, hash []byte) (io.ReadSeekerCloser, *object.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTag", ctx, repo, hash)
	ret0, _ := ret[0].(io.ReadSeekerCloser)
	ret1, _ := ret[1].(*object.Tag)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetTag indicates an expected call of GetTag
func (mr *MockStreamerMockRecorder) GetTag(ctx, repo, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTag", reflect.TypeOf((*MockStreamer)(nil).GetTag), ctx, repo, hash)
}

// OnRequest mocks base method
func (m *MockStreamer) OnRequest(s network.Stream) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnRequest", s)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OnRequest indicates an expected call of OnRequest
func (mr *MockStreamerMockRecorder) OnRequest(s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnRequest", reflect.TypeOf((*MockStreamer)(nil).OnRequest), s)
}

// GetProviders mocks base method
func (m *MockStreamer) GetProviders(ctx context.Context, repoName string, objectHash []byte) ([]peer.AddrInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProviders", ctx, repoName, objectHash)
	ret0, _ := ret[0].([]peer.AddrInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProviders indicates an expected call of GetProviders
func (mr *MockStreamerMockRecorder) GetProviders(ctx, repoName, objectHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProviders", reflect.TypeOf((*MockStreamer)(nil).GetProviders), ctx, repoName, objectHash)
}
