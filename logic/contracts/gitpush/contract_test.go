package gitpush_test

import (
	"os"
	"strings"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/makeos/mosdef/config"
	"gitlab.com/makeos/mosdef/crypto"
	logic2 "gitlab.com/makeos/mosdef/logic"
	"gitlab.com/makeos/mosdef/logic/contracts/gitpush"
	"gitlab.com/makeos/mosdef/mocks"
	"gitlab.com/makeos/mosdef/remote/plumbing"
	"gitlab.com/makeos/mosdef/remote/pushpool/types"
	types2 "gitlab.com/makeos/mosdef/remote/types"
	"gitlab.com/makeos/mosdef/remote/types/common"
	"gitlab.com/makeos/mosdef/storage"
	"gitlab.com/makeos/mosdef/testutil"
	"gitlab.com/makeos/mosdef/types/core"
	"gitlab.com/makeos/mosdef/types/state"
	"gitlab.com/makeos/mosdef/types/txns"
	"gitlab.com/makeos/mosdef/util"
)

var _ = Describe("AcquireNamespaceContract", func() {
	var appDB, stateTreeDB storage.Engine
	var err error
	var cfg *config.AppConfig
	var logic *logic2.Logic
	var ctrl *gomock.Controller
	var sender = crypto.NewKeyFromIntSeed(1)
	var mockRepoMgr *mocks.MockRemoteServer

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		cfg, err = testutil.SetTestCfg()
		Expect(err).To(BeNil())
		appDB, stateTreeDB = testutil.GetDB(cfg)
		logic = logic2.New(appDB, stateTreeDB, cfg)
		err := logic.SysKeeper().SaveBlockInfo(&core.BlockInfo{Height: 1})
		mockRepoMgr = mocks.NewMockRemoteServer(ctrl)
		logic.SetRemoteServer(mockRepoMgr)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		ctrl.Finish()
		Expect(appDB.Close()).To(BeNil())
		Expect(stateTreeDB.Close()).To(BeNil())
		err = os.RemoveAll(cfg.DataDir())
		Expect(err).To(BeNil())
	})

	Describe(".CanExec", func() {
		It("should return true when able to execute tx type", func() {
			ct := gitpush.NewContract()
			Expect(ct.CanExec(txns.TxTypePush)).To(BeTrue())
			Expect(ct.CanExec(txns.TxTypeCoinTransfer)).To(BeFalse())
		})
	})

	Describe(".Exec", func() {
		var repo = "repo1"
		var pushKeyID = crypto.CreatePushKeyID(crypto.StrToPublicKey("pushKeyID"))
		var creator = util.MustDecodePushKeyID(pushKeyID)
		var pushKeyID2 = crypto.CreatePushKeyID(crypto.StrToPublicKey("pushKeyID2"))
		var refs []*types.PushedReference
		var issueRef1 = plumbing.MakeIssueReference(1)
		var issue1RefData = &state.ReferenceData{Labels: []string{"label1", "label2"}, Assignees: []string{"key1", "key2"}}

		BeforeEach(func() {
			logic.AccountKeeper().Update(sender.Addr(), &state.Account{Balance: "10", Nonce: 1})
			logic.PushKeyKeeper().Update(pushKeyID, &state.PushKey{PubKey: crypto.StrToPublicKey("pub_key"), Address: sender.Addr()})
			logic.RepoKeeper().Update(repo, &state.Repository{
				References: map[string]*state.Reference{
					"refs/heads/master": {Nonce: 1, Creator: creator},
					issueRef1:           {Nonce: 1, Creator: creator, Data: issue1RefData},
				},
			})
		})

		When("pushed reference did not previously exist (new reference)", func() {
			It("should add pusher as creator of the reference", func() {
				mockRepoMgr.EXPECT().ExecTxPush(gomock.Any())
				refs = []*types.PushedReference{{Name: "refs/heads/dev", Data: &types2.ReferenceData{}}}
				rawPkID := util.MustDecodePushKeyID(pushKeyID2)
				err = gitpush.NewContract().Init(logic, &txns.TxPush{
					TxCommon: &txns.TxCommon{Fee: "1", SenderPubKey: sender.PubKey().ToPublicKey()},
					PushNote: &types.PushNote{RepoName: repo, References: refs, PushKeyID: rawPkID},
				}, 0).Exec()
				Expect(err).To(BeNil())
				rep := logic.RepoKeeper().Get(repo)
				Expect(rep.References.Get("refs/heads/dev").Creator).To(Equal(crypto.PushKey(rawPkID)))
			})
		})

		When("pushed reference new hash is zero (meaning delete is required)", func() {
			It("should remove the reference from the repo", func() {
				mockRepoMgr.EXPECT().ExecTxPush(gomock.Any())
				refs = []*types.PushedReference{{Name: "refs/heads/master", NewHash: strings.Repeat("0", 40)}}
				rawPkID := util.MustDecodePushKeyID(pushKeyID2)
				err = gitpush.NewContract().Init(logic, &txns.TxPush{
					TxCommon: &txns.TxCommon{Fee: "1", SenderPubKey: sender.PubKey().ToPublicKey()},
					PushNote: &types.PushNote{RepoName: repo, References: refs, PushKeyID: rawPkID},
				}, 0).Exec()
				Expect(err).To(BeNil())
				rep := logic.RepoKeeper().Get(repo)
				Expect(rep.References.Has("refs/heads/master")).To(BeFalse())
			})
		})

		When("pushed reference already exist", func() {
			It("should not update reference creator", func() {
				mockRepoMgr.EXPECT().ExecTxPush(gomock.Any())
				refs = []*types.PushedReference{{Name: "refs/heads/master", Data: (&types2.ReferenceData{})}}
				rawPkID := util.MustDecodePushKeyID(pushKeyID2)
				err = gitpush.NewContract().Init(logic, &txns.TxPush{
					TxCommon: &txns.TxCommon{Fee: "1", SenderPubKey: sender.PubKey().ToPublicKey()},
					PushNote: &types.PushNote{RepoName: repo, References: refs, PushKeyID: rawPkID},
				}, 0).Exec()
				Expect(err).To(BeNil())
				rep := logic.RepoKeeper().Get(repo)
				Expect(rep.References.Get("refs/heads/master").Creator).ToNot(Equal(rawPkID))
				actual := util.MustDecodePushKeyID(pushKeyID)
				Expect(rep.References.Get("refs/heads/master").Creator).To(Equal(crypto.PushKey(actual)))
			})
		})

		When("reference has nonce = 1. After a successful exec:", func() {
			BeforeEach(func() {
				mockRepoMgr.EXPECT().ExecTxPush(gomock.Any())
				refs = []*types.PushedReference{{Name: "refs/heads/master", Data: (&types2.ReferenceData{})}}
				rawPkID := util.MustDecodePushKeyID(pushKeyID)
				err = gitpush.NewContract().Init(logic, &txns.TxPush{
					TxCommon: &txns.TxCommon{Fee: "1", SenderPubKey: sender.PubKey().ToPublicKey()},
					PushNote: &types.PushNote{RepoName: repo, References: refs, PushKeyID: rawPkID},
				}, 0).Exec()
				Expect(err).To(BeNil())
			})

			Specify("that the reference's nonce was incremented", func() {
				repo := logic.RepoKeeper().Get(repo)
				Expect(repo.References.Get("refs/heads/master").Nonce).To(Equal(uint64(2)))
			})

			Specify("that fee was deducted from pusher account", func() {
				acct := logic.AccountKeeper().Get(sender.Addr())
				Expect(acct.Balance).To(Equal(util.String("9")))
			})

			Specify("that sender account nonce was incremented", func() {
				acct := logic.AccountKeeper().Get(sender.Addr())
				Expect(acct.Nonce).To(Equal(uint64(2)))
			})
		})

		When("pushed reference is an issue reference", func() {
			It("should add issue data from pushed reference", func() {
				mockRepoMgr.EXPECT().ExecTxPush(gomock.Any())
				ref := plumbing.MakeIssueReference(2)
				labels := []string{"lbl1"}
				assignees := []string{"key1"}
				cls := true
				refs = []*types.PushedReference{{Name: ref, Data: &types2.ReferenceData{Close: &cls, IssueFields: common.IssueFields{Labels: &labels, Assignees: &assignees}}}}
				rawPkID := util.MustDecodePushKeyID(pushKeyID2)
				err = gitpush.NewContract().Init(logic, &txns.TxPush{
					TxCommon: &txns.TxCommon{Fee: "1", SenderPubKey: sender.PubKey().ToPublicKey()},
					PushNote: &types.PushNote{RepoName: repo, References: refs, PushKeyID: rawPkID},
				}, 0).Exec()
				Expect(err).To(BeNil())
				rep := logic.RepoKeeper().Get(repo)
				Expect(rep.References.Get(ref).Creator).To(Equal(crypto.PushKey(rawPkID)))
				Expect(rep.References[ref].Data.Closed).To(BeTrue())
				Expect(rep.References[ref].Data.Labels).To(Equal(labels))
				Expect(rep.References[ref].Data.Assignees).To(Equal(assignees))
			})

			It("should not alter issue data fields not specified in pushed reference", func() {
				mockRepoMgr.EXPECT().ExecTxPush(gomock.Any())
				cls := true
				refs = []*types.PushedReference{{Name: issueRef1, Data: (&types2.ReferenceData{Close: &cls})}}
				err = gitpush.NewContract().Init(logic, &txns.TxPush{
					TxCommon: &txns.TxCommon{Fee: "1", SenderPubKey: sender.PubKey().ToPublicKey()},
					PushNote: &types.PushNote{RepoName: repo, References: refs, PushKeyID: creator},
				}, 0).Exec()
				Expect(err).To(BeNil())
				rep := logic.RepoKeeper().Get(repo)
				Expect(rep.References.Get(issueRef1).Creator).To(Equal(crypto.PushKey(creator)))
				Expect(rep.References[issueRef1].Data.Closed).To(BeTrue())
				Expect(rep.References[issueRef1].Data.Labels).To(Equal(issue1RefData.Labels))
				Expect(rep.References[issueRef1].Data.Assignees).To(Equal(issue1RefData.Assignees))
			})

			When("pushed reference label include a negated entry", func() {
				It("should remove the negated entry from the reference", func() {
					mockRepoMgr.EXPECT().ExecTxPush(gomock.Any())
					refs = []*types.PushedReference{{Name: issueRef1, Data: &types2.ReferenceData{IssueFields: common.IssueFields{Labels: &[]string{"-label1"}}}}}
					err = gitpush.NewContract().Init(logic, &txns.TxPush{
						TxCommon: &txns.TxCommon{Fee: "1", SenderPubKey: sender.PubKey().ToPublicKey()},
						PushNote: &types.PushNote{RepoName: repo, References: refs, PushKeyID: creator},
					}, 0).Exec()
					Expect(err).To(BeNil())
					rep := logic.RepoKeeper().Get(repo)
					Expect(rep.References[issueRef1].Data.Labels).To(HaveLen(1))
					Expect(rep.References[issueRef1].Data.Labels).To(ContainElement(issue1RefData.Labels[1]))
				})
			})

			When("pushed reference assignee include a negated entry", func() {
				It("should remove the negated entry from the reference", func() {
					mockRepoMgr.EXPECT().ExecTxPush(gomock.Any())
					refs = []*types.PushedReference{{Name: issueRef1, Data: &types2.ReferenceData{IssueFields: common.IssueFields{Assignees: &[]string{"-key1"}}}}}
					err = gitpush.NewContract().Init(logic, &txns.TxPush{
						TxCommon: &txns.TxCommon{Fee: "1", SenderPubKey: sender.PubKey().ToPublicKey()},
						PushNote: &types.PushNote{RepoName: repo, References: refs, PushKeyID: creator},
					}, 0).Exec()
					Expect(err).To(BeNil())
					rep := logic.RepoKeeper().Get(repo)
					Expect(rep.References[issueRef1].Data.Assignees).To(HaveLen(1))
					Expect(rep.References[issueRef1].Data.Assignees).To(ContainElement(issue1RefData.Assignees[1]))
				})
			})
		})
	})
})
