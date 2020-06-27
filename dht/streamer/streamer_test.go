package streamer_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/makeos/mosdef/config"
	"gitlab.com/makeos/mosdef/dht"
	"gitlab.com/makeos/mosdef/dht/streamer"
	types3 "gitlab.com/makeos/mosdef/dht/streamer/types"
	"gitlab.com/makeos/mosdef/mocks"
	"gitlab.com/makeos/mosdef/remote/plumbing"
	"gitlab.com/makeos/mosdef/remote/repo"
	"gitlab.com/makeos/mosdef/remote/types"
	"gitlab.com/makeos/mosdef/testutil"
	types2 "gitlab.com/makeos/mosdef/types"
	io2 "gitlab.com/makeos/mosdef/util/io"
	plumb "gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type fakePackfile struct {
	name string
}

func (f *fakePackfile) Read(p []byte) (n int, err error) {
	return
}

func (f *fakePackfile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f *fakePackfile) Close() error {
	return nil
}

var _ = Describe("BasicObjectStreamer", func() {
	var err error
	var cfg *config.AppConfig
	var ctrl *gomock.Controller
	var mockHost *mocks.MockHost
	var mockDHT *mocks.MockDHT
	var cs *streamer.BasicObjectStreamer
	var hash = plumb.NewHash("6fe5e981f7defdfb907c1237e2e8427696adafa7")
	var parentHash = plumb.NewHash("7a561e23f4e81c61df1b0dc63a89ae9c8d5680cd")

	BeforeEach(func() {
		cfg, err = testutil.SetTestCfg()
		Expect(err).To(BeNil())
		cfg.Node.GitBinPath = "/usr/bin/git"
		ctrl = gomock.NewController(GinkgoT())
		mockHost = mocks.NewMockHost(ctrl)
		mockDHT = mocks.NewMockDHT(ctrl)
	})

	BeforeEach(func() {
		mockHost.EXPECT().SetStreamHandler(gomock.Any(), gomock.Any())
		mockDHT.EXPECT().Host().Return(mockHost)
		cs = streamer.NewObjectStreamer(mockDHT, cfg)
	})

	AfterEach(func() {
		ctrl.Finish()
		err = os.RemoveAll(cfg.DataDir())
		Expect(err).To(BeNil())
	})

	Describe(".NewObjectStreamer", func() {
		It("should register commit stream protocol handler", func() {
			mockHost.EXPECT().SetStreamHandler(streamer.ObjectStreamerProtocolID, gomock.Any())
			mockDHT.EXPECT().Host().Return(mockHost)
			streamer.NewObjectStreamer(mockDHT, cfg)
		})
	})

	Describe(".Announce", func() {
		It("should announce commit hash", func() {
			mockDHT.EXPECT().Announce(dht.MakeObjectKey(hash[:]), gomock.Any())
			cs.Announce(hash[:], nil)
		})

		It("should return error when announce attempt failed", func() {
			mockDHT.EXPECT().Announce(dht.MakeObjectKey(hash[:]), gomock.Any()).Do(func(key []byte, doneCB func(error)) {
				doneCB(fmt.Errorf("error"))
			})
			cs.Announce(hash[:], func(err error) {
				Expect(err).ToNot(BeNil())
				Expect(err).To(MatchError("error"))
			})
		})
	})

	Describe(".OnRequest", func() {
		It("should return error when unable to read stream", func() {
			mockStream := mocks.NewMockStream(ctrl)
			mockStream.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
				return 0, fmt.Errorf("read error")
			})
			_, err := cs.OnRequest(mockStream)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to read request: read error"))
		})

		It("should return ErrUnknownMsgType when message type is unknown", func() {
			mockStream := mocks.NewMockStream(ctrl)
			mockStream.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
				msg := []byte("unknown")
				copy(p, msg)
				return len(msg), nil
			})
			_, err := cs.OnRequest(mockStream)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError(streamer.ErrUnknownMsgType))
		})

		It("should call 'Want' handler when message is MsgTypeWant", func() {
			msg := []byte(dht.MsgTypeWant)
			mockStream := mocks.NewMockStream(ctrl)
			mockStream.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
				copy(p, msg)
				return len(msg), nil
			})
			cs.OnWantHandler = func(m []byte, s network.Stream) error {
				Expect(msg).To(Equal(msg))
				return nil
			}
			_, err := cs.OnRequest(mockStream)
			Expect(err).To(BeNil())
		})

		It("should call 'Send' handler when message is MsgTypeSend", func() {
			msg := []byte(dht.MsgTypeSend)
			mockStream := mocks.NewMockStream(ctrl)
			mockStream.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
				copy(p, msg)
				return len(msg), nil
			})
			cs.OnSendHandler = func(m []byte, s network.Stream) error {
				Expect(msg).To(Equal(msg))
				return nil
			}
			success, err := cs.OnRequest(mockStream)
			Expect(err).To(BeNil())
			Expect(success).To(BeTrue())
		})
	})

	Describe(".OnWant", func() {
		It("should return error if msg could not be parsed", func() {
			mockStream := mocks.NewMockStream(ctrl)
			mockStream.EXPECT().Reset()
			err := cs.OnWant([]byte(""), mockStream)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("malformed message"))
		})

		It("should return error if unable to get local repository", func() {
			mockStream := mocks.NewMockStream(ctrl)
			mockStream.EXPECT().Reset()
			cs.RepoGetter = func(string, string) (types.LocalRepo, error) {
				return nil, fmt.Errorf("failed to get repo")
			}
			err := cs.OnWant(dht.MakeWantMsg("repo1", hash[:]), mockStream)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to get repo"))
		})

		It("should return error if extracted commit key is malformed", func() {
			mockStream := mocks.NewMockStream(ctrl)
			mockStream.EXPECT().Reset()
			mockRepo := mocks.NewMockLocalRepo(ctrl)
			cs.RepoGetter = func(string, string) (types.LocalRepo, error) {
				return mockRepo, nil
			}
			err := cs.OnWant(dht.MakeWantMsg("repo1", hash[:]), mockStream)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("malformed commit key"))
		})

		It("should return write 'NOPE' message to stream and return ErrObjNotFound if object does not exist", func() {
			mockStream := mocks.NewMockStream(ctrl)
			mockStream.EXPECT().Reset()
			mockRepo := mocks.NewMockLocalRepo(ctrl)
			mockRepo.EXPECT().ObjectExist(hash.String()).Return(false)
			mockStream.EXPECT().Write(dht.MakeNopeMsg())
			cs.RepoGetter = func(string, string) (types.LocalRepo, error) {
				return mockRepo, nil
			}
			key := dht.MakeObjectKey(hash[:])
			err := cs.OnWant(dht.MakeWantMsg("repo1", key), mockStream)
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(dht.ErrObjNotFound))
		})

		It("should return when unable to write 'NOPE' message to stream when object does not exist", func() {
			mockStream := mocks.NewMockStream(ctrl)
			mockStream.EXPECT().Reset()
			mockRepo := mocks.NewMockLocalRepo(ctrl)
			mockRepo.EXPECT().ObjectExist(hash.String()).Return(false)
			mockStream.EXPECT().Write(dht.MakeNopeMsg()).Return(0, fmt.Errorf("write error"))
			cs.RepoGetter = func(string, string) (types.LocalRepo, error) {
				return mockRepo, nil
			}
			key := dht.MakeObjectKey(hash[:])
			err := cs.OnWant(dht.MakeWantMsg("repo1", key), mockStream)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to write 'nope' message: write error"))
		})

		When("commit object exist in local repo", func() {
			It("should return error when writing 'HAVE' response failed", func() {
				mockStream := mocks.NewMockStream(ctrl)
				mockStream.EXPECT().Reset()
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				mockRepo.EXPECT().ObjectExist(hash.String()).Return(true)
				mockStream.EXPECT().Write(dht.MakeHaveMsg()).Return(0, fmt.Errorf("write error"))
				cs.RepoGetter = func(string, string) (types.LocalRepo, error) {
					return mockRepo, nil
				}
				key := dht.MakeObjectKey(hash[:])
				err := cs.OnWant(dht.MakeWantMsg("repo1", key), mockStream)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("write error"))
			})

			It("should return no error when writing 'HAVE' response succeeds", func() {
				mockStream := mocks.NewMockStream(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				mockRepo.EXPECT().ObjectExist(hash.String()).Return(true)
				mockStream.EXPECT().Write(dht.MakeHaveMsg()).Return(0, nil)
				cs.RepoGetter = func(string, string) (types.LocalRepo, error) {
					return mockRepo, nil
				}
				key := dht.MakeObjectKey(hash[:])
				err := cs.OnWant(dht.MakeWantMsg("repo1", key), mockStream)
				Expect(err).To(BeNil())
			})
		})
	})

	Describe(".OnSend", func() {

		It("should return error if msg could not be parsed", func() {
			mockStream := mocks.NewMockStream(ctrl)
			mockStream.EXPECT().Reset()
			err := cs.OnSend([]byte(""), mockStream)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("malformed message"))
		})

		It("should return error if unable to get local repository", func() {
			mockStream := mocks.NewMockStream(ctrl)
			mockStream.EXPECT().Reset()
			cs.RepoGetter = func(string, string) (types.LocalRepo, error) {
				return nil, fmt.Errorf("failed to get repo")
			}
			err := cs.OnSend(dht.MakeWantMsg("repo1", hash[:]), mockStream)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to get repo"))
		})

		It("should return error if extracted commit key is malformed", func() {
			mockStream := mocks.NewMockStream(ctrl)
			mockStream.EXPECT().Reset()
			mockRepo := mocks.NewMockLocalRepo(ctrl)
			cs.RepoGetter = func(string, string) (types.LocalRepo, error) {
				return mockRepo, nil
			}
			err := cs.OnSend(dht.MakeWantMsg("repo1", hash[:]), mockStream)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("malformed commit key"))
		})

		It("should return error when non-ErrObjectNotFound is returned when getting commit from local repo", func() {
			mockStream := mocks.NewMockStream(ctrl)
			mockStream.EXPECT().Reset()
			mockRepo := mocks.NewMockLocalRepo(ctrl)
			mockRepo.EXPECT().GetObject(hash.String()).Return(nil, fmt.Errorf("unexpected error"))
			cs.RepoGetter = func(string, string) (types.LocalRepo, error) {
				return mockRepo, nil
			}
			key := dht.MakeObjectKey(hash[:])
			err := cs.OnSend(dht.MakeWantMsg("repo1", key), mockStream)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("unexpected error"))
		})

		When("ErrObjectNotFound is returned when getting commit from local repo", func() {
			It("should return error when writing a 'NOPE' response failed", func() {
				mockStream := mocks.NewMockStream(ctrl)
				mockStream.EXPECT().Reset()
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				mockRepo.EXPECT().GetObject(hash.String()).Return(nil, plumb.ErrObjectNotFound)
				mockStream.EXPECT().Write(dht.MakeNopeMsg()).Return(0, fmt.Errorf("write error"))
				cs.RepoGetter = func(string, string) (types.LocalRepo, error) {
					return mockRepo, nil
				}
				key := dht.MakeObjectKey(hash[:])
				err := cs.OnSend(dht.MakeWantMsg("repo1", key), mockStream)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("failed to write 'nope' message: write error"))
			})
		})

		When("commit object exist in local repo", func() {
			It("should return error when generating a packfile for the commit failed", func() {
				mockStream := mocks.NewMockStream(ctrl)
				mockStream.EXPECT().Reset()
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				mockRepo.EXPECT().GetObject(hash.String()).Return(nil, nil)

				mockConn := mocks.NewMockConn(ctrl)
				mockConn.EXPECT().RemotePeer().Return(peer.ID("peer-id"))
				mockStream.EXPECT().Conn().Return(mockConn)

				cs.RepoGetter = func(string, string) (types.LocalRepo, error) {
					return mockRepo, nil
				}
				cs.PackObject = func(repo types.LocalRepo, args *plumbing.PackObjectArgs) (io.Reader, []plumb.Hash, error) {
					return nil, nil, fmt.Errorf("error")
				}
				key := dht.MakeObjectKey(hash[:])
				err := cs.OnSend(dht.MakeWantMsg("repo1", key), mockStream)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("failed to generate commit packfile: error"))
			})

			It("should return no error", func() {
				mockStream := mocks.NewMockStream(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				mockRepo.EXPECT().GetObject(hash.String()).Return(nil, nil)

				peerID := peer.ID("peer-id")
				mockConn := mocks.NewMockConn(ctrl)
				mockConn.EXPECT().RemotePeer().Return(peerID)
				mockStream.EXPECT().Conn().Return(mockConn)

				cs.RepoGetter = func(string, string) (types.LocalRepo, error) {
					return mockRepo, nil
				}
				objs := []plumb.Hash{
					plumb.NewHash("9f00445ef94ed0f78f95fb40a96c5eba22ab1f03"),
					plumb.NewHash("ba751747e0de82408417600288daa79221eda714"),
				}
				cs.PackObject = func(repo types.LocalRepo, args *plumbing.PackObjectArgs) (io.Reader, []plumb.Hash, error) {
					return bytes.NewReader(nil), objs, nil
				}
				mockStream.EXPECT().Close()

				key := dht.MakeObjectKey(hash[:])
				err := cs.OnSend(dht.MakeWantMsg("repo1", key), mockStream)
				Expect(err).To(BeNil())

				// It should add packed objects to the peer's HaveCache.
				cache := cs.HaveCache.GetCache(peerID.Pretty())
				Expect(cache.Has(hash))
				for _, obj := range objs {
					Expect(cache.Has(obj)).To(BeTrue())
				}
			})
		})
	})

	Describe(".GetCommit", func() {
		var ctx = context.Background()
		var repoName = "repo1"

		It("should return error when unable to get providers", func() {
			mockDHT.EXPECT().GetProviders(ctx, dht.MakeObjectKey(hash[:])).Return(nil, fmt.Errorf("error"))
			_, _, err := cs.GetCommit(ctx, repoName, hash[:])
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to get providers: error"))
		})

		It("should return ErrNoProviderFound when no provider is found", func() {
			mockDHT.EXPECT().GetProviders(ctx, dht.MakeObjectKey(hash[:])).Return(nil, nil)
			_, _, err := cs.GetCommit(ctx, repoName, hash[:])
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(streamer.ErrNoProviderFound))
		})

		It("should return error when request failed", func() {
			mockDHT.EXPECT().Host().Return(mockHost)

			prov := peer.AddrInfo{ID: "id", Addrs: []multiaddr.Multiaddr{multiaddr.StringCast("/ip4/127.0.0.1")}}
			mockDHT.EXPECT().GetProviders(ctx, dht.MakeObjectKey(hash[:])).Return([]peer.AddrInfo{prov}, nil)
			mockReq := mocks.NewMockObjectRequester(ctrl)
			mockReq.EXPECT().Do(ctx).Return(nil, fmt.Errorf("request error"))
			cs.MakeRequester = func(args streamer.RequestArgs) streamer.ObjectRequester {
				return mockReq
			}
			_, _, err := cs.GetCommit(ctx, repoName, hash[:])
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("request failed: request error"))
		})

		It("should return error when unable to get target object in packfile", func() {
			mockDHT.EXPECT().Host().Return(mockHost)

			prov := peer.AddrInfo{ID: "id", Addrs: []multiaddr.Multiaddr{multiaddr.StringCast("/ip4/127.0.0.1")}}
			mockDHT.EXPECT().GetProviders(ctx, dht.MakeObjectKey(hash[:])).Return([]peer.AddrInfo{prov}, nil)
			mockReq := mocks.NewMockObjectRequester(ctrl)
			mockReq.EXPECT().Do(ctx).Return(&streamer.PackResult{}, nil)
			cs.MakeRequester = func(args streamer.RequestArgs) streamer.ObjectRequester {
				return mockReq
			}
			cs.PackObjectGetter = func(io.ReadSeeker, string) (res object.Object, err error) {
				return nil, fmt.Errorf("error")
			}
			_, _, err := cs.GetCommit(ctx, repoName, hash[:])
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to get target commit from packfile: error"))
		})

		It("should return error when get target object does not exist in packfile", func() {
			mockDHT.EXPECT().Host().Return(mockHost)

			prov := peer.AddrInfo{ID: "id", Addrs: []multiaddr.Multiaddr{multiaddr.StringCast("/ip4/127.0.0.1")}}
			mockDHT.EXPECT().GetProviders(ctx, dht.MakeObjectKey(hash[:])).Return([]peer.AddrInfo{prov}, nil)
			mockReq := mocks.NewMockObjectRequester(ctrl)
			mockReq.EXPECT().Do(ctx).Return(&streamer.PackResult{}, nil)
			cs.MakeRequester = func(args streamer.RequestArgs) streamer.ObjectRequester {
				return mockReq
			}
			cs.PackObjectGetter = func(io.ReadSeeker, string) (res object.Object, err error) {
				return nil, nil
			}
			_, _, err := cs.GetCommit(ctx, repoName, hash[:])
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("target commit not found in the packfile"))
		})

		It("should return packfile on success", func() {
			mockDHT.EXPECT().Host().Return(mockHost)

			prov := peer.AddrInfo{ID: "id", Addrs: []multiaddr.Multiaddr{multiaddr.StringCast("/ip4/127.0.0.1")}}
			mockDHT.EXPECT().GetProviders(ctx, dht.MakeObjectKey(hash[:])).Return([]peer.AddrInfo{prov}, nil)
			mockReq := mocks.NewMockObjectRequester(ctrl)

			pack, err := ioutil.TempFile(os.TempDir(), "")
			Expect(err).To(BeNil())
			defer pack.Close()
			mockReq.EXPECT().Do(ctx).Return(&streamer.PackResult{Pack: pack}, nil)
			cs.MakeRequester = func(args streamer.RequestArgs) streamer.ObjectRequester {
				return mockReq
			}
			commit := object.Commit{Hash: hash}
			cs.PackObjectGetter = func(io.ReadSeeker, string) (res object.Object, err error) {
				return &commit, nil
			}
			res, _, err := cs.GetCommit(ctx, repoName, hash[:])
			Expect(err).To(BeNil())
			Expect(res).To(Equal(pack))
		})
	})

	Describe(".GetTag", func() {
		var ctx = context.Background()
		var repoName = "repo1"

		It("should return error when unable to get providers", func() {
			mockDHT.EXPECT().GetProviders(ctx, dht.MakeObjectKey(hash[:])).Return(nil, fmt.Errorf("error"))
			_, _, err := cs.GetTag(ctx, repoName, hash[:])
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to get providers: error"))
		})

		It("should return ErrNoProviderFound when no provider is found", func() {
			mockDHT.EXPECT().GetProviders(ctx, dht.MakeObjectKey(hash[:])).Return(nil, nil)
			_, _, err := cs.GetTag(ctx, repoName, hash[:])
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(streamer.ErrNoProviderFound))
		})

		It("should return error when request failed", func() {
			mockDHT.EXPECT().Host().Return(mockHost)

			prov := peer.AddrInfo{ID: "id", Addrs: []multiaddr.Multiaddr{multiaddr.StringCast("/ip4/127.0.0.1")}}
			mockDHT.EXPECT().GetProviders(ctx, dht.MakeObjectKey(hash[:])).Return([]peer.AddrInfo{prov}, nil)
			mockReq := mocks.NewMockObjectRequester(ctrl)
			mockReq.EXPECT().Do(ctx).Return(nil, fmt.Errorf("request error"))
			cs.MakeRequester = func(args streamer.RequestArgs) streamer.ObjectRequester {
				return mockReq
			}
			_, _, err := cs.GetTag(ctx, repoName, hash[:])
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("request failed: request error"))
		})

		It("should return error when packfile failed validation", func() {
			mockDHT.EXPECT().Host().Return(mockHost)

			prov := peer.AddrInfo{ID: "id", Addrs: []multiaddr.Multiaddr{multiaddr.StringCast("/ip4/127.0.0.1")}}
			mockDHT.EXPECT().GetProviders(ctx, dht.MakeObjectKey(hash[:])).Return([]peer.AddrInfo{prov}, nil)
			mockReq := mocks.NewMockObjectRequester(ctrl)
			mockReq.EXPECT().Do(ctx).Return(&streamer.PackResult{}, nil)
			cs.MakeRequester = func(args streamer.RequestArgs) streamer.ObjectRequester {
				return mockReq
			}
			cs.PackObjectGetter = func(io.ReadSeeker, string) (res object.Object, err error) {
				return nil, fmt.Errorf("error")
			}
			_, _, err := cs.GetTag(ctx, repoName, hash[:])
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to get target tag from packfile: error"))
		})

		It("should return packfile on success", func() {
			mockDHT.EXPECT().Host().Return(mockHost)

			prov := peer.AddrInfo{ID: "id", Addrs: []multiaddr.Multiaddr{multiaddr.StringCast("/ip4/127.0.0.1")}}
			mockDHT.EXPECT().GetProviders(ctx, dht.MakeObjectKey(hash[:])).Return([]peer.AddrInfo{prov}, nil)
			mockReq := mocks.NewMockObjectRequester(ctrl)

			pack, err := ioutil.TempFile(os.TempDir(), "")
			Expect(err).To(BeNil())
			defer pack.Close()
			mockReq.EXPECT().Do(ctx).Return(&streamer.PackResult{Pack: pack}, nil)
			cs.MakeRequester = func(args streamer.RequestArgs) streamer.ObjectRequester {
				return mockReq
			}
			tag := object.Tag{Hash: hash}
			cs.PackObjectGetter = func(io.ReadSeeker, string) (res object.Object, err error) {
				return &tag, nil
			}
			res, _, err := cs.GetTag(ctx, repoName, hash[:])
			Expect(err).To(BeNil())
			Expect(res).To(Equal(pack))
		})
	})

	Describe(".GetTaggedCommitWithAncestors", func() {
		var ctx = context.Background()
		var repoName = "repo1"

		It("should return error when unable to get target repository", func() {
			cs := mocks.NewMockObjectStreamer(ctrl)
			_, err := streamer.GetTaggedCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
				return nil, fmt.Errorf("error")
			}, types3.GetAncestorArgs{})
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to get repo: error"))
		})

		When("end commit hash is provided", func() {
			It("should return error if end commit object does not exist locally", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				mockRepo.EXPECT().GetObject(plumbing.BytesToHex(hash[:])).Return(nil, plumb.ErrObjectNotFound)
				_, err := streamer.GetTaggedCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					EndHash: hash[:],
				})
				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(streamer.ErrEndObjMustExistLocally))
			})

			It("should return error if unable to get end hash object locally", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				mockRepo.EXPECT().GetObject(plumbing.BytesToHex(hash[:])).Return(nil, fmt.Errorf("error"))
				_, err := streamer.GetTaggedCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					EndHash: hash[:],
				})
				Expect(err).ToNot(BeNil())
				Expect(err).To(MatchError("error"))
			})

			It("should return error end object is not a tag", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				obj := object.Commit{}
				mockRepo.EXPECT().GetObject(plumbing.BytesToHex(hash[:])).Return(&obj, nil)
				_, err := streamer.GetTaggedCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					EndHash: hash[:],
				})
				Expect(err).ToNot(BeNil())
				Expect(err).To(MatchError("end hash must be a tag object"))
			})

			It("should return error end object is a tag that does not point to a commit or a tag", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				obj := object.Tag{TargetType: plumb.BlobObject}
				mockRepo.EXPECT().GetObject(plumbing.BytesToHex(hash[:])).Return(&obj, nil)
				_, err := streamer.GetTaggedCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					EndHash: hash[:],
				})
				Expect(err).ToNot(BeNil())
				Expect(err).To(MatchError("end tag must point to a tag or commit object"))
			})

			When("end object is a tag that points to another tag", func() {
				Specify("that the pointed tag's target is recursively checked", func() {
					cs := mocks.NewMockObjectStreamer(ctrl)
					mockRepo := mocks.NewMockLocalRepo(ctrl)

					tag2Hash := plumb.NewHash("6081bfcf869e310ed06304641fdf7c365a03ac56")
					tag1 := object.Tag{TargetType: plumb.TagObject, Target: tag2Hash}
					commitHash := plumb.NewHash("3114383fe03a7b441ce5a0a6ac43a1f83622ba1a")
					tag2 := object.Tag{TargetType: plumb.CommitObject, Target: commitHash}

					mockRepo.EXPECT().GetObject(plumbing.BytesToHex(hash[:])).Return(&tag1, nil)
					mockRepo.EXPECT().GetObject(tag2Hash.String()).Return(&tag2, nil)

					cs.EXPECT().GetTag(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, fmt.Errorf("error"))
					streamer.GetTaggedCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
						return mockRepo, nil
					}, types3.GetAncestorArgs{
						EndHash: hash[:],
					})
				})
			})
		})

		When("start hash is an existing tag", func() {
			It("should return error when unable to get start tag ", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				cs.EXPECT().GetTag(ctx, repoName, hash[:]).Return(nil, nil, fmt.Errorf("error"))
				_, err := streamer.GetTaggedCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					RepoName:  repoName,
					StartHash: hash[:],
				})
				Expect(err).ToNot(BeNil())
				Expect(err).To(MatchError("error"))
			})

			It("should return no error and start tag packfile when tag does not point to a commit or tag", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				targetTag := &object.Tag{TargetType: plumb.BlobObject}
				tagPackfile := &fakePackfile{"pack-1"}
				cs.EXPECT().GetTag(ctx, repoName, hash[:]).Return(tagPackfile, targetTag, nil)
				packfiles, err := streamer.GetTaggedCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					RepoName:  repoName,
					StartHash: hash[:],
				})
				Expect(err).To(BeNil())
				Expect(packfiles).To(HaveLen(1))
				Expect(packfiles[0]).To(Equal(tagPackfile))
			})

			When("tag points to another tag", func() {
				It("should return error if unable to get pointed tag", func() {
					targetHash := "7a561e23f4e81c61df1b0dc63a89ae9c8d5680cd"
					cs := mocks.NewMockObjectStreamer(ctrl)
					mockRepo := mocks.NewMockLocalRepo(ctrl)
					targetTag := &object.Tag{TargetType: plumb.TagObject, Target: plumb.NewHash(targetHash)}
					cs.EXPECT().GetTag(ctx, repoName, hash[:]).Return(&fakePackfile{"pack-1"}, targetTag, nil)
					cs.EXPECT().GetTag(ctx, repoName, plumbing.HashToBytes(targetHash)).Return(nil, nil, fmt.Errorf("error"))
					_, err := streamer.GetTaggedCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
						return mockRepo, nil
					}, types3.GetAncestorArgs{
						RepoName:  repoName,
						StartHash: hash[:],
					})
					Expect(err).ToNot(BeNil())
					Expect(err).To(MatchError("error"))
				})
			})

			When("tag points to another tag", func() {
				It("should try to get pointed tag by calling GetTaggedCommitWithAncestors", func() {
					targetHash := "7a561e23f4e81c61df1b0dc63a89ae9c8d5680cd"
					cs := mocks.NewMockObjectStreamer(ctrl)
					mockRepo := mocks.NewMockLocalRepo(ctrl)
					targetTag := &object.Tag{TargetType: plumb.TagObject, Target: plumb.NewHash(targetHash)}
					cs.EXPECT().GetTag(ctx, repoName, hash[:]).Return(&fakePackfile{"pack-1"}, targetTag, nil)
					cs.EXPECT().GetTag(ctx, repoName, plumbing.HashToBytes(targetHash)).Return(nil, nil, fmt.Errorf("error"))
					_, err := streamer.GetTaggedCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
						return mockRepo, nil
					}, types3.GetAncestorArgs{
						RepoName:  repoName,
						StartHash: hash[:],
					})
					Expect(err).ToNot(BeNil())
					Expect(err).To(MatchError("error"))
				})
			})

			When("tag points to another commit", func() {
				It("should try to get ancestor of pointed commit by calling GetCommitWithAncestors", func() {
					targetHash := "7a561e23f4e81c61df1b0dc63a89ae9c8d5680cd"
					endHash := "c8ecc929fc8ef7964ef9d445a03e85e9f88c9d99"
					cs := mocks.NewMockObjectStreamer(ctrl)
					mockRepo := mocks.NewMockLocalRepo(ctrl)
					targetTag := &object.Tag{TargetType: plumb.CommitObject, Target: plumb.NewHash(targetHash)}
					endTag := &object.Tag{TargetType: plumb.CommitObject, Target: plumb.NewHash(endHash)}

					mockRepo.EXPECT().GetObject(gomock.Any()).Return(endTag, nil)
					cs.EXPECT().GetTag(ctx, repoName, hash[:]).Return(&fakePackfile{"pack-1"}, targetTag, nil)

					cs.EXPECT().GetCommitWithAncestors(ctx, types3.GetAncestorArgs{
						RepoName:  repoName,
						StartHash: targetTag.Target[:],
						EndHash:   endTag.Target[:],
					}).Return(nil, fmt.Errorf("error"))

					streamer.GetTaggedCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
						return mockRepo, nil
					}, types3.GetAncestorArgs{
						RepoName:  repoName,
						StartHash: hash[:],
						EndHash:   plumbing.HashToBytes(endHash),
					})

				})
			})
		})
	})

	Describe(".GetCommitWithAncestors", func() {
		var ctx = context.Background()
		var repoName = "repo1"

		It("should return error when unable to get target repository", func() {
			cs := mocks.NewMockObjectStreamer(ctrl)
			_, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
				return nil, fmt.Errorf("error")
			}, types3.GetAncestorArgs{})
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to get repo: error"))
		})

		When("end commit hash is provided", func() {
			It("should return error if end commit object does not exist locally", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				mockRepo.EXPECT().ObjectExist(hash.String()).Return(false)
				_, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					EndHash: hash[:],
				})
				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(streamer.ErrEndObjMustExistLocally))
			})
		})

		It("should return error on failed attempt to get start object locally", func() {
			cs := mocks.NewMockObjectStreamer(ctrl)
			mockRepo := mocks.NewMockLocalRepo(ctrl)
			mockRepo.EXPECT().CommitObject(hash).Return(nil, fmt.Errorf("error"))

			_, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
				return mockRepo, nil
			}, types3.GetAncestorArgs{
				StartHash: hash[:],
				RepoName:  repoName,
			})
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("error"))
		})

		When("start commit exist locally", func() {
			It("should not attempt to get start commit from the DHT", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				mockRepo.EXPECT().CommitObject(hash).Return(&object.Commit{}, nil)

				_, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					StartHash: hash[:],
					RepoName:  repoName,
				})

				Expect(err).To(BeNil())
			})

			It("should add parent to the waiting list and fetch it from the DHT", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				startCommit := &object.Commit{ParentHashes: []plumb.Hash{parentHash}}
				mockRepo.EXPECT().CommitObject(hash).Return(startCommit, nil)

				// Mock expectations for parent
				mockRepo.EXPECT().CommitObject(parentHash).Return(nil, plumb.ErrObjectNotFound)
				cs.EXPECT().GetCommit(ctx, repoName, parentHash[:]).Return(&fakePackfile{}, &object.Commit{}, nil)

				packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					StartHash: hash[:],
					RepoName:  repoName,
				})

				Expect(err).To(BeNil())
				Expect(packfiles).To(HaveLen(1))
			})
		})

		It("should return error if end unable to get start hash from DHT", func() {
			cs := mocks.NewMockObjectStreamer(ctrl)
			mockRepo := mocks.NewMockLocalRepo(ctrl)

			mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
			cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(nil, nil, fmt.Errorf("error"))

			_, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
				return mockRepo, nil
			}, types3.GetAncestorArgs{
				StartHash: hash[:],
				RepoName:  repoName,
			})
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("error"))
		})

		When("start commit hash and end commit hash match", func() {
			It("should return start commit pack file when ExcludeEndCommit is false", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				startCommit := &object.Commit{Hash: hash}
				startCommitPackfile := &fakePackfile{"pack-1"}
				mockRepo.EXPECT().ObjectExist(hash.String()).Return(true)
				mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
				cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)

				packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					StartHash:        hash[:],
					EndHash:          hash[:],
					RepoName:         repoName,
					ExcludeEndCommit: false,
				})
				Expect(err).To(BeNil())
				Expect(packfiles).To(HaveLen(1))
				Expect(packfiles[0]).To(Equal(startCommitPackfile))
			})

			It("should not return start commit pack file when ExcludeEndCommit is true", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				startCommit := &object.Commit{Hash: hash}
				startCommitPackfile := &fakePackfile{"pack-1"}
				mockRepo.EXPECT().ObjectExist(hash.String()).Return(true)
				mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
				cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)

				packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					StartHash:        hash[:],
					EndHash:          hash[:],
					RepoName:         repoName,
					ExcludeEndCommit: true,
				})
				Expect(err).To(BeNil())
				Expect(packfiles).To(HaveLen(0))
			})
		})

		When("start commit does not have parents", func() {
			It("should return start commit pack file alone", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				startCommit := &object.Commit{Hash: hash}
				startCommitPackfile := &fakePackfile{"pack-1"}
				mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
				cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)

				packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					StartHash: hash[:],
					RepoName:  repoName,
				})
				Expect(err).To(BeNil())
				Expect(packfiles).To(HaveLen(1))
				Expect(packfiles[0]).To(Equal(startCommitPackfile))
			})
		})

		When("start commit has one parent with matching hash", func() {

			It("should return only start commit and packfile", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				startCommit := &object.Commit{Hash: hash}
				startCommit.ParentHashes = append(startCommit.ParentHashes, hash)
				startCommitPackfile := &fakePackfile{"pack-1"}

				mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
				cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)
				mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)

				packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					StartHash: hash[:],
					RepoName:  repoName,
				})
				Expect(err).To(BeNil())
				Expect(packfiles).To(HaveLen(1))
				Expect(packfiles[0]).To(Equal(startCommitPackfile))
			})
		})

		When("start commit has one parent", func() {
			var parentHash = plumb.NewHash("7a561e23f4e81c61df1b0dc63a89ae9c8d5680cd")

			It("should return start commit and its parent commit packfiles", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				startCommit := &object.Commit{Hash: hash}
				startCommit.ParentHashes = append(startCommit.ParentHashes, parentHash)
				startCommitPackfile := &fakePackfile{"pack-1"}

				mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
				cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)

				parentCommit := &object.Commit{Hash: parentHash}
				parentCommitPackfile := &fakePackfile{"pack-2"}
				cs.EXPECT().GetCommit(ctx, repoName, parentHash[:]).Return(parentCommitPackfile, parentCommit, nil)

				mockRepo.EXPECT().CommitObject(parentHash).Return(nil, plumb.ErrObjectNotFound)

				packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					StartHash: hash[:],
					RepoName:  repoName,
				})
				Expect(err).To(BeNil())
				Expect(packfiles).To(HaveLen(2))
				Expect(packfiles[0]).To(Equal(startCommitPackfile))
				Expect(packfiles[1]).To(Equal(parentCommitPackfile))
			})

			When("unable to get start commit's parent object", func() {
				It("should return error and start commit packfile", func() {
					cs := mocks.NewMockObjectStreamer(ctrl)
					mockRepo := mocks.NewMockLocalRepo(ctrl)
					startCommit := &object.Commit{Hash: hash}
					startCommit.ParentHashes = append(startCommit.ParentHashes, parentHash)
					startCommitPackfile := &fakePackfile{"pack-1"}

					mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
					cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)

					mockRepo.EXPECT().CommitObject(parentHash).Return(nil, fmt.Errorf("error"))

					packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
						return mockRepo, nil
					}, types3.GetAncestorArgs{
						StartHash: hash[:],
						RepoName:  repoName,
					})
					Expect(err).ToNot(BeNil())
					Expect(err).To(MatchError("error"))
					Expect(packfiles).To(HaveLen(1))
					Expect(packfiles[0]).To(Equal(startCommitPackfile))
				})
			})

			// Test that when the start commit has a parent that exist locally,
			// that parent's parent will be added to the wantlist if it does not exists locally.
			When("start commit parent already exist", func() {
				It("should add the start commit's grand-parent to the waitlist if the grand-parent does not exist locally", func() {
					cs := mocks.NewMockObjectStreamer(ctrl)
					mockRepo := mocks.NewMockLocalRepo(ctrl)

					startCommit := &object.Commit{Hash: hash}
					startCommit.ParentHashes = append(startCommit.ParentHashes, parentHash)
					startCommitPackfile := &fakePackfile{"pack-1"}
					mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
					cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)

					grandParentHash := plumb.NewHash("d9dbe0e59248c7f0505dd5d80ed470fb43f82521")
					parentCommit := &object.Commit{Hash: parentHash, ParentHashes: []plumb.Hash{grandParentHash}}
					mockRepo.EXPECT().CommitObject(parentHash).Return(parentCommit, nil)
					mockRepo.EXPECT().CommitObject(grandParentHash).Return(nil, nil)

					grandParentCommit := &object.Commit{Hash: grandParentHash}
					grandParentCommitPackfile := &fakePackfile{"pack-2"}
					cs.EXPECT().GetCommit(ctx, repoName, grandParentHash[:]).Return(grandParentCommitPackfile, grandParentCommit, nil)

					packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
						return mockRepo, nil
					}, types3.GetAncestorArgs{
						StartHash: hash[:],
						RepoName:  repoName,
					})
					Expect(err).To(BeNil())
					Expect(packfiles).To(HaveLen(2))
					Expect(packfiles[0]).To(Equal(startCommitPackfile))
					Expect(packfiles[1]).To(Equal(grandParentCommitPackfile))
				})
			})

			It("should not add parent to wantlist if parent is the end commit and ExcludeEndCommit=true", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				startCommit := &object.Commit{Hash: hash}
				startCommit.ParentHashes = append(startCommit.ParentHashes, parentHash)
				startCommitPackfile := &fakePackfile{"pack-1"}
				mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
				cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)
				mockRepo.EXPECT().ObjectExist(parentHash.String()).Return(true)

				packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					StartHash:        hash[:],
					RepoName:         repoName,
					EndHash:          parentHash[:],
					ExcludeEndCommit: true,
				})
				Expect(err).To(BeNil())
				Expect(packfiles).To(HaveLen(1))
				Expect(packfiles[0]).To(Equal(startCommitPackfile))
			})

			When("start commit parent is the end commit, it exists locally and ExcludeEndCommit=false", func() {
				It("should add start commit parent to wantlist", func() {
					cs := mocks.NewMockObjectStreamer(ctrl)
					mockRepo := mocks.NewMockLocalRepo(ctrl)
					startCommit := &object.Commit{Hash: hash}
					startCommit.ParentHashes = append(startCommit.ParentHashes, parentHash)
					startCommitPackfile := &fakePackfile{"pack-1"}
					mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
					cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)
					mockRepo.EXPECT().ObjectExist(parentHash.String()).Return(true)

					parentCommit := &object.Commit{Hash: parentHash}
					parentCommitPackfile := &fakePackfile{"pack-2"}
					cs.EXPECT().GetCommit(ctx, repoName, parentHash[:]).Return(parentCommitPackfile, parentCommit, nil)

					packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
						return mockRepo, nil
					}, types3.GetAncestorArgs{
						StartHash:        hash[:],
						RepoName:         repoName,
						EndHash:          parentHash[:],
						ExcludeEndCommit: false,
					})
					Expect(err).To(BeNil())
					Expect(packfiles).To(HaveLen(2))
					Expect(packfiles[0]).To(Equal(startCommitPackfile))
					Expect(packfiles[1]).To(Equal(parentCommitPackfile))
				})
			})
		})

		When("start commit has two parents", func() {
			var parentHash = plumb.NewHash("7a561e23f4e81c61df1b0dc63a89ae9c8d5680cd")
			var parent2Hash = plumb.NewHash("c988dcc9fc47958626c8bd1b956817e5b5bb0105")

			It("should return start commit and its parents commit packfiles", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				startCommit := &object.Commit{Hash: hash}
				startCommit.ParentHashes = append(startCommit.ParentHashes, parentHash, parent2Hash)
				startCommitPackfile := &fakePackfile{"pack-1"}
				mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
				cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)

				parentCommit := &object.Commit{Hash: parentHash}
				parentCommitPackfile := &fakePackfile{"pack-2"}
				cs.EXPECT().GetCommit(ctx, repoName, parentHash[:]).Return(parentCommitPackfile, parentCommit, nil)
				mockRepo.EXPECT().CommitObject(parentHash).Return(nil, plumb.ErrObjectNotFound)

				parent2Commit := &object.Commit{Hash: parent2Hash}
				parent2CommitPackfile := &fakePackfile{"pack-3"}
				cs.EXPECT().GetCommit(ctx, repoName, parent2Hash[:]).Return(parent2CommitPackfile, parent2Commit, nil)
				mockRepo.EXPECT().CommitObject(parent2Hash).Return(nil, plumb.ErrObjectNotFound)

				packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					StartHash: hash[:],
					RepoName:  repoName,
				})
				Expect(err).To(BeNil())
				Expect(packfiles).To(HaveLen(3))
				Expect(packfiles[0]).To(Equal(startCommitPackfile))
				Expect(packfiles[1]).To(Equal(parentCommitPackfile))
				Expect(packfiles[2]).To(Equal(parent2CommitPackfile))
			})
		})

		When("end commit has been fetched/seen but wantlist is not empty", func() {
			var parentHash = plumb.NewHash("7a561e23f4e81c61df1b0dc63a89ae9c8d5680cd")
			var parent2Hash = plumb.NewHash("c988dcc9fc47958626c8bd1b956817e5b5bb0105")

			It("should add object to wantlist if ErrObjectNotFound is returned while performing ancestor check", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				startCommit := &object.Commit{Hash: hash}
				startCommit.ParentHashes = append(startCommit.ParentHashes, parentHash, parent2Hash)
				startCommitPackfile := &fakePackfile{"pack-1"}

				mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
				cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)

				mockRepo.EXPECT().ObjectExist(parentHash.String()).Return(true)
				mockRepo.EXPECT().CommitObject(parent2Hash).Return(nil, plumb.ErrObjectNotFound)

				parentCommit := &object.Commit{Hash: parentHash}
				parentCommitPackfile := &fakePackfile{"pack-2"}
				cs.EXPECT().GetCommit(ctx, repoName, parentHash[:]).Return(parentCommitPackfile, parentCommit, nil)

				parent2Commit := &object.Commit{Hash: parent2Hash}
				parent2CommitPackfile := &fakePackfile{"pack-3"}
				cs.EXPECT().GetCommit(ctx, repoName, parent2Hash[:]).Return(parent2CommitPackfile, parent2Commit, nil)

				mockRepo.EXPECT().IsAncestor(parent2Hash.String(), parentHash.String()).Return(plumb.ErrObjectNotFound)

				packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					StartHash: hash[:],
					RepoName:  repoName,
					EndHash:   parentHash[:],
				})
				Expect(err).To(BeNil())
				Expect(packfiles).To(HaveLen(3))
				Expect(packfiles[0]).To(Equal(startCommitPackfile))
				Expect(packfiles[1]).To(Equal(parentCommitPackfile))
				Expect(packfiles[2]).To(Equal(parent2CommitPackfile))
			})

			It("should add object to wantlist if ErrNotAnAncestor is returned while performing ancestor check", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				startCommit := &object.Commit{Hash: hash}
				startCommit.ParentHashes = append(startCommit.ParentHashes, parentHash)
				startCommit.ParentHashes = append(startCommit.ParentHashes, parent2Hash)
				startCommitPackfile := &fakePackfile{"pack-1"}

				mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
				cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)

				mockRepo.EXPECT().ObjectExist(parentHash.String()).Return(true)
				mockRepo.EXPECT().CommitObject(parent2Hash).Return(nil, plumb.ErrObjectNotFound)

				parentCommit := &object.Commit{Hash: parentHash}
				parentCommitPackfile := &fakePackfile{"pack-2"}
				cs.EXPECT().GetCommit(ctx, repoName, parentHash[:]).Return(parentCommitPackfile, parentCommit, nil)

				parent2Commit := &object.Commit{Hash: parent2Hash}
				parent2CommitPackfile := &fakePackfile{"pack-3"}
				cs.EXPECT().GetCommit(ctx, repoName, parent2Hash[:]).Return(parent2CommitPackfile, parent2Commit, nil)

				mockRepo.EXPECT().IsAncestor(parent2Hash.String(), parentHash.String()).Return(repo.ErrNotAnAncestor)

				packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					StartHash: hash[:],
					RepoName:  repoName,
					EndHash:   parentHash[:],
				})
				Expect(err).To(BeNil())
				Expect(packfiles).To(HaveLen(3))
				Expect(packfiles[0]).To(Equal(startCommitPackfile))
				Expect(packfiles[1]).To(Equal(parentCommitPackfile))
				Expect(packfiles[2]).To(Equal(parent2CommitPackfile))
			})

			It("should return error and current packfiles result if a non-ErrNotAnAncestor is returned while performing ancestor check", func() {
				cs := mocks.NewMockObjectStreamer(ctrl)
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				startCommit := &object.Commit{Hash: hash}
				startCommit.ParentHashes = append(startCommit.ParentHashes, parentHash)
				startCommit.ParentHashes = append(startCommit.ParentHashes, parent2Hash)
				startCommitPackfile := &fakePackfile{"pack-1"}

				mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
				cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)

				mockRepo.EXPECT().ObjectExist(parentHash.String()).Return(true)
				mockRepo.EXPECT().CommitObject(parent2Hash).Return(nil, plumb.ErrObjectNotFound)

				mockRepo.EXPECT().IsAncestor(parent2Hash.String(), parentHash.String()).Return(fmt.Errorf("bad error"))

				packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
					return mockRepo, nil
				}, types3.GetAncestorArgs{
					StartHash: hash[:],
					RepoName:  repoName,
					EndHash:   parentHash[:],
				})
				Expect(err).ToNot(BeNil())
				Expect(err).To(MatchError("failed to perform ancestor check: bad error"))
				Expect(packfiles).To(HaveLen(1))
				Expect(packfiles[0]).To(Equal(startCommitPackfile))
			})
		})

		Context("use callback to collect result, instead method returned result", func() {
			var parentHash = plumbing.HashToBytes("7a561e23f4e81c61df1b0dc63a89ae9c8d5680cd")

			When("ResultCB is provided", func() {
				It("should pass result to the callback and zero packfiles must be returned from the method", func() {
					cs := mocks.NewMockObjectStreamer(ctrl)
					mockRepo := mocks.NewMockLocalRepo(ctrl)
					startCommit := &object.Commit{Hash: hash}
					startCommitPackfile := &fakePackfile{"pack-1"}

					mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
					cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)

					cbPackfiles := []io2.ReadSeekerCloser{}
					packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
						return mockRepo, nil
					}, types3.GetAncestorArgs{
						StartHash: hash[:],
						RepoName:  repoName,
						ResultCB: func(packfile io2.ReadSeekerCloser, hash string) error {
							cbPackfiles = append(cbPackfiles, packfile)
							return nil
						},
					})
					Expect(err).To(BeNil())
					Expect(packfiles).To(HaveLen(0))
					Expect(cbPackfiles).To(HaveLen(1))
					Expect(cbPackfiles[0]).To(Equal(startCommitPackfile))
				})
			})

			When("callback returns non-ErrExit error", func() {
				It("should return start commit and its parent commit packfiles", func() {
					cs := mocks.NewMockObjectStreamer(ctrl)
					mockRepo := mocks.NewMockLocalRepo(ctrl)
					startCommit := &object.Commit{Hash: hash}
					startCommit.ParentHashes = append(startCommit.ParentHashes, plumb.NewHash(plumbing.BytesToHex(parentHash)))
					startCommitPackfile := &fakePackfile{"pack-1"}

					mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
					cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)

					cbPackfiles := []io2.ReadSeekerCloser{}
					packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
						return mockRepo, nil
					}, types3.GetAncestorArgs{
						StartHash: hash[:],
						RepoName:  repoName,
						ResultCB: func(packfile io2.ReadSeekerCloser, hash string) error {
							cbPackfiles = append(cbPackfiles, packfile)
							return fmt.Errorf("error")
						},
					})
					Expect(err).ToNot(BeNil())
					Expect(err).To(MatchError("error"))
					Expect(packfiles).To(HaveLen(0))
					Expect(cbPackfiles).To(HaveLen(1))
					Expect(cbPackfiles[0]).To(Equal(startCommitPackfile))
				})
			})

			When("callback returns ErrExit error", func() {
				It("should return start commit and its parent commit packfiles", func() {
					cs := mocks.NewMockObjectStreamer(ctrl)
					mockRepo := mocks.NewMockLocalRepo(ctrl)
					startCommit := &object.Commit{Hash: hash}
					startCommit.ParentHashes = append(startCommit.ParentHashes, plumb.NewHash(plumbing.BytesToHex(parentHash)))
					startCommitPackfile := &fakePackfile{"pack-1"}

					mockRepo.EXPECT().CommitObject(hash).Return(nil, plumb.ErrObjectNotFound)
					cs.EXPECT().GetCommit(ctx, repoName, hash[:]).Return(startCommitPackfile, startCommit, nil)

					cbPackfiles := []io2.ReadSeekerCloser{}
					packfiles, err := streamer.GetCommitWithAncestors(ctx, cs, func(gitBinPath, path string) (types.LocalRepo, error) {
						return mockRepo, nil
					}, types3.GetAncestorArgs{
						StartHash: hash[:],
						RepoName:  repoName,
						ResultCB: func(packfile io2.ReadSeekerCloser, hash string) error {
							cbPackfiles = append(cbPackfiles, packfile)
							return types2.ErrExit
						},
					})
					Expect(err).To(BeNil())
					Expect(packfiles).To(HaveLen(0))
					Expect(cbPackfiles).To(HaveLen(1))
					Expect(cbPackfiles[0]).To(Equal(startCommitPackfile))
				})
			})
		})
	})

})