package refsync_test

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"gitlab.com/makeos/mosdef/config"
	"gitlab.com/makeos/mosdef/mocks"
	"gitlab.com/makeos/mosdef/remote/push"
	"gitlab.com/makeos/mosdef/remote/push/types"
	repo3 "gitlab.com/makeos/mosdef/remote/repo"
	testutil2 "gitlab.com/makeos/mosdef/remote/testutil"
	types2 "gitlab.com/makeos/mosdef/remote/types"
	"gitlab.com/makeos/mosdef/testutil"
	"gitlab.com/makeos/mosdef/types/txns"
	"gitlab.com/makeos/mosdef/util"
	"gopkg.in/src-d/go-git.v4/plumbing"

	. "github.com/onsi/gomega"
	. "gitlab.com/makeos/mosdef/remote/refsync"
)

var _ = Describe("RefSync", func() {
	var err error
	var cfg *config.AppConfig
	var rs RefSyncer
	var ctrl *gomock.Controller
	var mockFetcher *mocks.MockObjectFetcher
	var repoName string
	var oldHash = "5b9ba1de20344b12cce76256b67cff9bb31e77b2"
	var newHash = "8d998c7de21bbe561f7992bb983cef4b1554993b"
	var path string

	BeforeEach(func() {
		cfg, err = testutil.SetTestCfg()
		Expect(err).To(BeNil())
		cfg.Node.GitBinPath = "/usr/bin/git"
		ctrl = gomock.NewController(GinkgoT())
		mockFetcher = mocks.NewMockObjectFetcher(ctrl)
		rs = New(cfg, mockFetcher, 1)

		repoName = util.RandString(5)
		testutil2.ExecGit(cfg.GetRepoRoot(), "init", repoName)
		path = filepath.Join(cfg.GetRepoRoot(), repoName)
	})

	Describe(".Start", func() {
		It("should panic if called twice", func() {
			rs.Start()
			Expect(func() { rs.Start() }).To(Panic())
		})
	})

	Describe(".IsRunning", func() {
		It("should return false if not running", func() {
			Expect(rs.IsRunning()).To(BeFalse())
		})

		It("should return true if running", func() {
			rs.Start()
			Expect(rs.IsRunning()).To(BeTrue())
		})
	})

	Describe(".Stop", func() {
		It("should set running state to false", func() {
			rs.Start()
			Expect(rs.IsRunning()).To(BeTrue())
			rs.Stop()
			Expect(rs.IsRunning()).To(BeFalse())
		})
	})

	Describe(".HasTax", func() {
		It("should return false when task queue is empty", func() {
			Expect(rs.QueueSize()).To(BeZero())
			Expect(rs.HasTask()).To(BeFalse())
		})
	})

	Describe(".OnNewTx", func() {
		It("should add new task to queue", func() {
			rs.OnNewTx(&txns.TxPush{Note: &types.Note{
				References: []*types.PushedReference{{Name: "master", Nonce: 1}},
			}})
			Expect(rs.HasTask()).To(BeTrue())
			Expect(rs.QueueSize()).To(Equal(1))
		})

		It("should not add new task to queue when task with matching ID already exist in queue", func() {
			rs.OnNewTx(&txns.TxPush{Note: &types.Note{References: []*types.PushedReference{{Name: "master", Nonce: 1}}}})
			rs.OnNewTx(&txns.TxPush{Note: &types.Note{References: []*types.PushedReference{{Name: "master", Nonce: 1}}}})
			Expect(rs.HasTask()).To(BeTrue())
			Expect(rs.QueueSize()).To(Equal(1))
		})

		It("should not add task if reference new hash is zero-hash", func() {
			rs.OnNewTx(&txns.TxPush{Note: &types.Note{References: []*types.PushedReference{
				{Name: "master", Nonce: 1, NewHash: plumbing.ZeroHash.String()},
			}}})
			Expect(rs.HasTask()).To(BeFalse())
			Expect(rs.QueueSize()).To(Equal(0))
		})

		It("should add two tasks if push transaction contains 2 different references", func() {
			rs.OnNewTx(&txns.TxPush{Note: &types.Note{
				References: []*types.PushedReference{
					{Name: "refs/heads/master", Nonce: 1},
					{Name: "refs/heads/dev", Nonce: 1},
				},
			}})
			Expect(rs.HasTask()).To(BeTrue())
			Expect(rs.QueueSize()).To(Equal(2))
		})
	})

	Describe(".Start", func() {
		It("should set the status to start", func() {
			rs.Start()
			Expect(rs.IsRunning()).To(BeTrue())
			rs.Stop()
		})
	})

	Describe(".Do", func() {
		It("should append task back to queue when task's has a future next run time", func() {
			task := &Task{Ref: &types.PushedReference{}, NextRunTime: time.Now().Add(1 * time.Second)}
			err := Do(rs.(*RefSync), task, 0)
			Expect(err).To(BeNil())
			Expect(rs.QueueSize()).To(Equal(1))
		})

		It("should append task back to queue when another task with matching reference name is being finalized", func() {
			task := &Task{Ref: &types.PushedReference{Name: "refs/heads/master"}}
			rs.(*RefSync).FinalizingRefs[task.Ref.Name] = struct{}{}
			err := Do(rs.(*RefSync), task, 0)
			Expect(err).To(BeNil())
			Expect(rs.QueueSize()).To(Equal(1))
		})

		It("should return error when repo does not exist locally", func() {
			task := &Task{RepoName: "unknown", Ref: &types.PushedReference{Name: "refs/heads/master"}}
			err := Do(rs.(*RefSync), task, 0)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to get target repo: repository does not exist"))
		})

		It("should return error when unable to get reference from local repo", func() {
			task := &Task{RepoName: "unknown", Ref: &types.PushedReference{Name: "refs/heads/master"}}
			mockRepo := mocks.NewMockLocalRepo(ctrl)
			rs.(*RefSync).RepoGetter = func(gitBinPath, path string) (types2.LocalRepo, error) { return mockRepo, nil }
			mockRepo.EXPECT().RefGet(task.Ref.Name).Return("", fmt.Errorf("error"))
			err := Do(rs.(*RefSync), task, 0)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to get reference from target repo: error"))
		})

		When("local reference hash and task reference old hash do not match", func() {
			var localHash = "2630fe8660633d5c543d4484769d148fae255b3e"

			It("should add task back to queue, set next run time and increment compat retry count", func() {
				task := &Task{RepoName: "unknown", Ref: &types.PushedReference{Name: "refs/heads/master", OldHash: oldHash, NewHash: newHash}}
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				rs.(*RefSync).RepoGetter = func(gitBinPath, path string) (types2.LocalRepo, error) { return mockRepo, nil }
				mockRepo.EXPECT().RefGet(task.Ref.Name).Return(localHash, nil)
				mockFetcher.EXPECT().Fetch(gomock.Any(), gomock.Any())
				err := Do(rs.(*RefSync), task, 0)
				Expect(err).To(BeNil())
				Expect(rs.QueueSize()).To(Equal(1))
				Expect(time.Now().Before(task.NextRunTime)).To(BeTrue())
				Expect(task.CompatRetryCount).To(Equal(1))
			})

			When("compat retry count has reached the max", func() {
				var err error
				var task *Task

				BeforeEach(func() {
					MaxCompatRetries = 1
					task = &Task{RepoName: "unknown", Ref: &types.PushedReference{Name: "refs/heads/master", OldHash: oldHash, NewHash: newHash}}
					mockRepo := mocks.NewMockLocalRepo(ctrl)
					rs.(*RefSync).RepoGetter = func(gitBinPath, path string) (types2.LocalRepo, error) { return mockRepo, nil }
					mockRepo.EXPECT().RefGet(task.Ref.Name).Return(localHash, nil)
					mockFetcher.EXPECT().Fetch(gomock.Any(), gomock.Any())
					err = Do(rs.(*RefSync), task, 0)
				})

				It("should return error", func() {
					Expect(err).ToNot(BeNil())
					Expect(err).To(MatchError("reference is not compatible with local state"))
				})

				It("should add task back to queue and increment compat retry count", func() {
					Expect(rs.QueueSize()).To(Equal(1))
					Expect(time.Now().Before(task.NextRunTime)).To(BeFalse())
					Expect(task.CompatRetryCount).To(Equal(1))
				})

				It("should update task's old hash to the local reference old hash", func() {
					Expect(task.Ref.OldHash).To(Equal(localHash))
				})
			})
		})

		When("local reference hash and task reference old hash match", func() {
			It("should attempt to fetch objects and update repo", func() {
				task := &Task{RepoName: "unknown", Ref: &types.PushedReference{Name: "refs/heads/master", OldHash: oldHash, NewHash: newHash}}
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				rs.(*RefSync).RepoGetter = func(gitBinPath, path string) (types2.LocalRepo, error) { return mockRepo, nil }
				mockRepo.EXPECT().RefGet(task.Ref.Name).Return(oldHash, nil)
				updated := false
				rs.(*RefSync).UpdateRepoUsingNote = func(string, push.ReferenceUpdateRequestPackMaker, types.PushNote) error {
					updated = true
					return nil
				}
				mockFetcher.EXPECT().Fetch(gomock.Any(), gomock.Any()).Do(func(note types.PushNote, cb func(err error)) {
					Expect(note).ToNot(BeNil())
					cb(nil)
				})
				err := Do(rs.(*RefSync), task, 0)
				Expect(err).To(BeNil())
				Expect(updated).To(BeTrue())
			})

			It("should not attempt to update repo if fetch attempt failed", func() {
				task := &Task{RepoName: "unknown", Ref: &types.PushedReference{Name: "refs/heads/master", OldHash: oldHash, NewHash: newHash}}
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				rs.(*RefSync).RepoGetter = func(gitBinPath, path string) (types2.LocalRepo, error) { return mockRepo, nil }
				mockRepo.EXPECT().RefGet(task.Ref.Name).Return(oldHash, nil)
				updated := false
				rs.(*RefSync).UpdateRepoUsingNote = func(string, push.ReferenceUpdateRequestPackMaker, types.PushNote) error {
					updated = true
					return nil
				}
				mockFetcher.EXPECT().Fetch(gomock.Any(), gomock.Any()).Do(func(note types.PushNote, cb func(err error)) {
					Expect(note).ToNot(BeNil())
					cb(fmt.Errorf("error"))
				})
				err := Do(rs.(*RefSync), task, 0)
				Expect(err).To(BeNil())
				Expect(updated).To(BeFalse())
			})

			It("should update repo without fetching objects if node is the creator of the push note", func() {
				key, _ := cfg.G().PrivVal.GetKey()
				task := &Task{
					RepoName:    "unknown",
					Ref:         &types.PushedReference{Name: "refs/heads/master", OldHash: oldHash, NewHash: newHash},
					NoteCreator: key.PubKey().MustBytes32(),
				}
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				rs.(*RefSync).RepoGetter = func(gitBinPath, path string) (types2.LocalRepo, error) { return mockRepo, nil }
				mockRepo.EXPECT().RefGet(task.Ref.Name).Return(oldHash, nil)
				updated := false
				rs.(*RefSync).UpdateRepoUsingNote = func(string, push.ReferenceUpdateRequestPackMaker, types.PushNote) error {
					updated = true
					return nil
				}
				err := Do(rs.(*RefSync), task, 0)
				Expect(err).To(BeNil())
				Expect(updated).To(BeTrue())
			})

			It("should update repo without fetching objects if node is an endorser of the push note", func() {
				key, _ := cfg.G().PrivVal.GetKey()
				task := &Task{
					RepoName: "unknown",
					Ref:      &types.PushedReference{Name: "refs/heads/master", OldHash: oldHash, NewHash: newHash},
					Endorsements: []*types.PushEndorsement{
						{EndorserPubKey: key.PubKey().MustBytes32()},
					},
				}
				mockRepo := mocks.NewMockLocalRepo(ctrl)
				rs.(*RefSync).RepoGetter = func(gitBinPath, path string) (types2.LocalRepo, error) { return mockRepo, nil }
				mockRepo.EXPECT().RefGet(task.Ref.Name).Return(oldHash, nil)
				updated := false
				rs.(*RefSync).UpdateRepoUsingNote = func(string, push.ReferenceUpdateRequestPackMaker, types.PushNote) error {
					updated = true
					return nil
				}
				err := Do(rs.(*RefSync), task, 0)
				Expect(err).To(BeNil())
				Expect(updated).To(BeTrue())
			})
		})
	})

	Describe(".UpdateRepoUsingNote", func() {
		var repo types2.LocalRepo

		BeforeEach(func() {
			repo, err = repo3.GetWithLiteGit(cfg.Node.GitBinPath, path)
			Expect(err).To(BeNil())
		})

		It("should return error if unable to create packfile from note", func() {
			note := &types.Note{}
			err := UpdateRepoUsingNote(cfg.Node.GitBinPath, func(tx types.PushNote) (io.ReadSeeker, error) {
				return nil, fmt.Errorf("error")
			}, note)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to create packfile from push note: error"))
		})

		It("should return failed to run git-receive if repo path is invalid", func() {
			mockRepo := mocks.NewMockLocalRepo(ctrl)
			mockRepo.EXPECT().GetPath().Return("invalid/path/to/repo")
			note := &types.Note{TargetRepo: mockRepo}
			buf := strings.NewReader("invalid")
			err := UpdateRepoUsingNote(cfg.Node.GitBinPath, func(tx types.PushNote) (io.ReadSeeker, error) {
				return buf, nil
			}, note)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(MatchRegexp("failed to start git-receive-pack command"))
		})

		It("should return error when generated packfile is invalid", func() {
			note := &types.Note{TargetRepo: repo}
			buf := strings.NewReader("invalid")
			err := UpdateRepoUsingNote(cfg.Node.GitBinPath, func(tx types.PushNote) (io.ReadSeeker, error) {
				return buf, nil
			}, note)
			Expect(err).ToNot(BeNil())
		})

		It("should return nil when packfile is valid", func() {
			testutil2.AppendCommit(path, "file.txt", "some text", "commit msg")
			commitHash := testutil2.GetRecentCommitHash(path, "refs/heads/master")
			note := &types.Note{
				TargetRepo: repo,
				References: []*types.PushedReference{
					{Name: "refs/heads/master", NewHash: commitHash, OldHash: plumbing.ZeroHash.String()},
				},
			}
			packfile, err := push.MakeReferenceUpdateRequestPack(note)
			Expect(err).To(BeNil())
			err = UpdateRepoUsingNote(cfg.Node.GitBinPath, func(tx types.PushNote) (io.ReadSeeker, error) {
				return packfile, nil
			}, note)
			Expect(err).To(BeNil())
		})
	})
})