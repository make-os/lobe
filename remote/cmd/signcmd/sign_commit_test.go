package signcmd

import (
	"fmt"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	restclient "gitlab.com/makeos/mosdef/api/rest/client"
	"gitlab.com/makeos/mosdef/api/rpc/client"
	"gitlab.com/makeos/mosdef/config"
	"gitlab.com/makeos/mosdef/crypto"
	"gitlab.com/makeos/mosdef/mocks"
	"gitlab.com/makeos/mosdef/remote/repo"
	"gitlab.com/makeos/mosdef/testutil"
	"gitlab.com/makeos/mosdef/types"
	"gitlab.com/makeos/mosdef/types/core"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var testGetNextNonce = func(pushKeyID string, rpcClient *client.RPCClient, remoteClients []restclient.RestClient) (string, error) {
	return "1", nil
}

func testPushKeyUnlocker(key core.StoredKey, err error) func(cfg *config.AppConfig, pushKeyID,
	defaultPassphrase string, targetRepo core.BareRepo) (core.StoredKey, error) {
	return func(cfg *config.AppConfig, pushKeyID, defaultPassphrase string, targetRepo core.BareRepo) (core.StoredKey, error) {
		return key, err
	}
}

func testRemoteURLTokenUpdater(token string, err error) func(targetRepo core.BareRepo, targetRemote string,
	txDetail *types.TxDetail, pushKey core.StoredKey, reset bool) (string, error) {
	return func(targetRepo core.BareRepo, targetRemote string, txDetail *types.TxDetail, pushKey core.StoredKey, reset bool) (string, error) {
		return token, err
	}
}

var _ = Describe("SignCommit", func() {
	var err error
	var cfg *config.AppConfig
	var ctrl *gomock.Controller
	var mockRepo *mocks.MockBareRepo
	var key *crypto.Key

	BeforeEach(func() {
		cfg, err = testutil.SetTestCfg()
		Expect(err).To(BeNil())
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockBareRepo(ctrl)
		key = crypto.NewKeyFromIntSeed(1)
	})

	AfterEach(func() {
		ctrl.Finish()
		err = os.RemoveAll(cfg.DataDir())
		Expect(err).To(BeNil())
	})

	Describe(".SignCommitCmd", func() {
		It("should return error when push key ID is not provided", func() {
			mockRepo.EXPECT().GetConfig("user.signingKey").Return("")
			args := &SignCommitArgs{}
			err := SignCommitCmd(cfg, mockRepo, args)
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(ErrMissingPushKeyID))
		})

		It("should return error when unable to find and unlock push key", func() {
			mockRepo.EXPECT().GetConfig("user.signingKey").Return(key.PushAddr().String())
			args := &SignCommitArgs{}
			args.PushKeyUnlocker = testPushKeyUnlocker(nil, fmt.Errorf("error"))
			err := SignCommitCmd(cfg, mockRepo, args)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("unable to unlock push key: error"))
		})

		It("should return error when mergeID is set but invalid", func() {
			mockRepo.EXPECT().GetConfig("user.signingKey").Return(key.PushAddr().String())
			args := &SignCommitArgs{MergeID: "abc123_invalid"}
			mockStoredKey := mocks.NewMockStoredKey(ctrl)
			args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
			err := SignCommitCmd(cfg, mockRepo, args)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("merge id must be numeric"))
			args.MergeID = "12345678910"
			err = SignCommitCmd(cfg, mockRepo, args)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("merge proposal id exceeded 8 bytes limit"))
		})

		It("should return error when unable to get next nonce", func() {
			mockRepo.EXPECT().GetConfig("user.signingKey").Return(key.PushAddr().String())
			args := &SignCommitArgs{GetNextNonce: func(pushKeyID string, rpcClient *client.RPCClient, remoteClients []restclient.RestClient) (string, error) {
				return "", fmt.Errorf("error")
			},
			}
			mockStoredKey := mocks.NewMockStoredKey(ctrl)
			args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
			err := SignCommitCmd(cfg, mockRepo, args)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("error"))
		})

		It("should return error when unable to get local repo HEAD", func() {
			mockRepo.EXPECT().GetConfig("user.signingKey").Return(key.PushAddr().String())
			args := &SignCommitArgs{GetNextNonce: testGetNextNonce}
			mockRepo.EXPECT().Head().Return("", fmt.Errorf("error"))
			mockStoredKey := mocks.NewMockStoredKey(ctrl)
			args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
			err := SignCommitCmd(cfg, mockRepo, args)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("failed to get HEAD"))
		})

		When("args.Branch is set", func() {
			It("should return error when unable to checkout branch", func() {
				mockRepo.EXPECT().GetConfig("user.signingKey").Return(key.PushAddr().String())
				args := &SignCommitArgs{GetNextNonce: testGetNextNonce, Branch: "refs/heads/dev"}
				mockStoredKey := mocks.NewMockStoredKey(ctrl)
				args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
				mockRepo.EXPECT().Head().Return("refs/heads/master", nil)
				mockRepo.EXPECT().Checkout("dev", false, args.ForceCheckout).Return(fmt.Errorf("error"))
				err := SignCommitCmd(cfg, mockRepo, args)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("failed to checkout branch (refs/heads/dev): error"))
			})

			It("should return error when unable to checkout branch", func() {
				mockRepo.EXPECT().GetConfig("user.signingKey").Return(key.PushAddr().String())
				args := &SignCommitArgs{GetNextNonce: testGetNextNonce, Branch: "refs/heads/dev"}
				mockStoredKey := mocks.NewMockStoredKey(ctrl)
				args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
				mockRepo.EXPECT().Head().Return("refs/heads/master", nil)
				mockRepo.EXPECT().Checkout("dev", false, args.ForceCheckout).Return(fmt.Errorf("error"))
				err := SignCommitCmd(cfg, mockRepo, args)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("failed to checkout branch (refs/heads/dev): error"))
			})
		})

		It("should return error when unable to create and set push tokens to remote URLs", func() {
			args := &SignCommitArgs{Fee: "1", PushKeyID: key.PushAddr().String(), GetNextNonce: testGetNextNonce}
			args.RemoteURLTokenUpdater = testRemoteURLTokenUpdater("", fmt.Errorf("error"))
			mockStoredKey := mocks.NewMockStoredKey(ctrl)
			args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
			mockRepo.EXPECT().Head().Return("refs/heads/master", nil)
			err := SignCommitCmd(cfg, mockRepo, args)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("error"))
		})

		When("previous commit amendment is not required (AmendCommit=false", func() {
			It("should attempt to create a new commit and return error on failure", func() {
				args := &SignCommitArgs{Fee: "1", PushKeyID: key.PushAddr().String(), Message: "some message", GetNextNonce: testGetNextNonce, AmendCommit: false}
				args.RemoteURLTokenUpdater = testRemoteURLTokenUpdater("", nil)
				mockStoredKey := mocks.NewMockStoredKey(ctrl)
				args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
				mockRepo.EXPECT().Head().Return("refs/heads/master", nil)
				mockRepo.EXPECT().CreateSignedEmptyCommit(args.Message, args.PushKeyID).Return(fmt.Errorf("error"))
				err := SignCommitCmd(cfg, mockRepo, args)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("error"))
			})

			It("should attempt to create a new commit and return nil on success", func() {
				args := &SignCommitArgs{Fee: "1", PushKeyID: key.PushAddr().String(), Message: "some message", GetNextNonce: testGetNextNonce, AmendCommit: false}
				args.RemoteURLTokenUpdater = testRemoteURLTokenUpdater("", nil)
				mockStoredKey := mocks.NewMockStoredKey(ctrl)
				args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
				mockRepo.EXPECT().Head().Return("refs/heads/master", nil)
				mockRepo.EXPECT().CreateSignedEmptyCommit(args.Message, args.PushKeyID).Return(nil)
				err := SignCommitCmd(cfg, mockRepo, args)
				Expect(err).To(BeNil())
			})
		})

		When("args.Head is set to 'refs/heads/some_branch'", func() {
			It("the generated tx detail should set Reference to 'refs/heads/some_branch'", func() {
				args := &SignCommitArgs{Fee: "1", PushKeyID: key.PushAddr().String(), Message: "some message", GetNextNonce: testGetNextNonce, AmendCommit: false, Head: "refs/heads/some_branch"}
				args.RemoteURLTokenUpdater = func(targetRepo core.BareRepo, targetRemote string, txDetail *types.TxDetail, pushKey core.StoredKey, reset bool) (string, error) {
					Expect(txDetail.Reference).To(Equal("refs/heads/some_branch"))
					return "", nil
				}

				mockStoredKey := mocks.NewMockStoredKey(ctrl)
				args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
				mockRepo.EXPECT().Head().Return("refs/heads/master", nil)
				mockRepo.EXPECT().CreateSignedEmptyCommit(args.Message, args.PushKeyID).Return(nil)
				err := SignCommitCmd(cfg, mockRepo, args)
				Expect(err).To(BeNil())
			})
		})

		When("amend of previous commit is required (AmendCommit=true)", func() {
			It("should return error when unable to get recent commit hash", func() {
				args := &SignCommitArgs{Fee: "1", PushKeyID: key.PushAddr().String(), Message: "some message", GetNextNonce: testGetNextNonce, AmendCommit: true}
				args.RemoteURLTokenUpdater = testRemoteURLTokenUpdater("", nil)

				mockStoredKey := mocks.NewMockStoredKey(ctrl)
				args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
				mockRepo.EXPECT().Head().Return("refs/heads/master", nil)
				mockRepo.EXPECT().GetRecentCommitHash().Return("", fmt.Errorf("error"))
				err := SignCommitCmd(cfg, mockRepo, args)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("error"))
			})

			It("should return error when unable to get recent commit hash", func() {
				args := &SignCommitArgs{Fee: "1", PushKeyID: key.PushAddr().String(), Message: "some message", GetNextNonce: testGetNextNonce, AmendCommit: true}
				args.RemoteURLTokenUpdater = testRemoteURLTokenUpdater("", nil)
				mockStoredKey := mocks.NewMockStoredKey(ctrl)
				args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
				mockRepo.EXPECT().Head().Return("refs/heads/master", nil)
				mockRepo.EXPECT().GetRecentCommitHash().Return("", fmt.Errorf("error"))
				err := SignCommitCmd(cfg, mockRepo, args)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("error"))
			})

			It("should return error when unable to get recent commit due to ErrNoCommits", func() {
				args := &SignCommitArgs{Fee: "1", PushKeyID: key.PushAddr().String(), Message: "some message", GetNextNonce: testGetNextNonce, AmendCommit: true}
				args.RemoteURLTokenUpdater = testRemoteURLTokenUpdater("", nil)
				mockStoredKey := mocks.NewMockStoredKey(ctrl)
				args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
				mockRepo.EXPECT().Head().Return("refs/heads/master", nil)
				mockRepo.EXPECT().GetRecentCommitHash().Return("", repo.ErrNoCommits)
				err := SignCommitCmd(cfg, mockRepo, args)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("no commits have been created yet"))
			})

			It("should use previous commit message if args.Message is not set", func() {
				args := &SignCommitArgs{Fee: "1", PushKeyID: key.PushAddr().String(), Message: "", GetNextNonce: testGetNextNonce, AmendCommit: true}
				args.RemoteURLTokenUpdater = testRemoteURLTokenUpdater("", nil)
				mockStoredKey := mocks.NewMockStoredKey(ctrl)
				args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
				mockRepo.EXPECT().Head().Return("refs/heads/master", nil)
				mockRepo.EXPECT().GetRecentCommitHash().Return("8975220bda0eb48354c5868f8a1b310758eb4591", nil)
				recentCommit := &object.Commit{Message: "This is a commit"}
				mockRepo.EXPECT().CommitObject(plumbing.NewHash("8975220bda0eb48354c5868f8a1b310758eb4591")).Return(recentCommit, nil)
				mockRepo.EXPECT().UpdateRecentCommitMsg(recentCommit.Message, args.PushKeyID).Return(nil)
				err := SignCommitCmd(cfg, mockRepo, args)
				Expect(err).To(BeNil())
				Expect(args.Message).To(Equal(recentCommit.Message))
			})

			It("should return error when unable to get recent commit object", func() {
				args := &SignCommitArgs{Fee: "1", PushKeyID: key.PushAddr().String(), Message: "", GetNextNonce: testGetNextNonce, AmendCommit: true}
				args.RemoteURLTokenUpdater = testRemoteURLTokenUpdater("", nil)
				mockStoredKey := mocks.NewMockStoredKey(ctrl)
				args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
				mockRepo.EXPECT().Head().Return("refs/heads/master", nil)
				mockRepo.EXPECT().GetRecentCommitHash().Return("8975220bda0eb48354c5868f8a1b310758eb4591", nil)
				mockRepo.EXPECT().CommitObject(plumbing.NewHash("8975220bda0eb48354c5868f8a1b310758eb4591")).Return(nil, fmt.Errorf("error getting commit"))
				err := SignCommitCmd(cfg, mockRepo, args)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("error getting commit"))
			})

			It("should return error when unable to update recent commit", func() {
				args := &SignCommitArgs{Fee: "1", PushKeyID: key.PushAddr().String(), Message: "", GetNextNonce: testGetNextNonce, AmendCommit: true}
				args.RemoteURLTokenUpdater = testRemoteURLTokenUpdater("", nil)
				mockStoredKey := mocks.NewMockStoredKey(ctrl)
				args.PushKeyUnlocker = testPushKeyUnlocker(mockStoredKey, nil)
				mockRepo.EXPECT().Head().Return("refs/heads/master", nil)
				mockRepo.EXPECT().GetRecentCommitHash().Return("8975220bda0eb48354c5868f8a1b310758eb4591", nil)
				recentCommit := &object.Commit{Message: "This is a commit"}
				mockRepo.EXPECT().CommitObject(plumbing.NewHash("8975220bda0eb48354c5868f8a1b310758eb4591")).Return(recentCommit, nil)
				mockRepo.EXPECT().UpdateRecentCommitMsg(recentCommit.Message, args.PushKeyID).Return(fmt.Errorf("error"))
				err := SignCommitCmd(cfg, mockRepo, args)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("error"))
			})
		})
	})
})
