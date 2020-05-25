package mergerequest_test

import (
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/makeos/mosdef/config"
	"gitlab.com/makeos/mosdef/crypto"
	logic2 "gitlab.com/makeos/mosdef/logic"
	"gitlab.com/makeos/mosdef/logic/contracts/mergerequest"
	"gitlab.com/makeos/mosdef/storage"
	"gitlab.com/makeos/mosdef/testutil"
	"gitlab.com/makeos/mosdef/types/constants"
	"gitlab.com/makeos/mosdef/types/core"
	"gitlab.com/makeos/mosdef/types/state"
	"gitlab.com/makeos/mosdef/types/txns"
	"gitlab.com/makeos/mosdef/util"
)

var _ = Describe("MergeRequestContract", func() {
	var appDB, stateTreeDB storage.Engine
	var err error
	var cfg *config.AppConfig
	var logic *logic2.Logic
	var ctrl *gomock.Controller
	var sender = crypto.NewKeyFromIntSeed(1)
	var key2 = crypto.NewKeyFromIntSeed(2)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		cfg, err = testutil.SetTestCfg()
		Expect(err).To(BeNil())
		appDB, stateTreeDB = testutil.GetDB(cfg)
		logic = logic2.New(appDB, stateTreeDB, cfg)
		err := logic.SysKeeper().SaveBlockInfo(&core.BlockInfo{Height: 1})
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
			data := &mergerequest.MergeRequestData{}
			ct := mergerequest.NewContract(data)
			Expect(ct.CanExec(txns.MergeRequestProposalAction)).To(BeTrue())
			Expect(ct.CanExec(txns.TxTypeRepoProposalSendFee)).To(BeFalse())
		})
	})

	Describe(".Exec", func() {
		var err error
		var repo *state.Repository

		BeforeEach(func() {
			logic.AccountKeeper().Update(sender.Addr(), &state.Account{
				Balance:             "10",
				Stakes:              state.BareAccountStakes(),
				DelegatorCommission: 10,
			})
			repo = state.BareRepository()
			repo.Config = state.DefaultRepoConfig
			repo.Config.Governance.Voter = state.VoterOwner
		})

		When("sender is the only owner", func() {
			repoName := "repo"
			proposalFee := util.String("1")
			id := "1"

			BeforeEach(func() {
				repo.AddOwner(sender.Addr().String(), &state.RepoOwner{})

				err = mergerequest.NewContract(&mergerequest.MergeRequestData{
					Repo:             repo,
					RepoName:         repoName,
					ProposalID:       id,
					ProposerFee:      proposalFee,
					Fee:              "1.5",
					CreatorAddress:   sender.Addr(),
					BaseBranch:       "base",
					BaseBranchHash:   "baseHash",
					TargetBranch:     "target",
					TargetBranchHash: "targetHash",
				}).Init(logic, nil, 0).Exec()
				Expect(err).To(BeNil())
			})

			It("should add the new proposal to the repo", func() {
				Expect(repo.Proposals).To(HaveLen(1))
			})

			Specify("that the proposal is finalized and self accepted", func() {
				Expect(repo.Proposals).To(HaveLen(1))
				propID := mergerequest.MakeMergeRequestID(id)
				Expect(repo.Proposals.Get(propID).IsFinalized()).To(BeTrue())
				Expect(repo.Proposals.Get(propID).Yes).To(Equal(float64(1)))
			})

			Specify("that no fee was deducted", func() {
				acct := logic.AccountKeeper().Get(sender.Addr(), 0)
				Expect(acct.Balance.String()).To(Equal("10"))
			})

			Specify("that the proposal fee by the sender is registered on the proposal", func() {
				propID := mergerequest.MakeMergeRequestID(id)
				Expect(repo.Proposals.Get(propID).Fees).To(HaveLen(1))
				Expect(repo.Proposals.Get(propID).Fees).To(HaveKey(sender.Addr().String()))
				Expect(repo.Proposals.Get(propID).Fees[sender.Addr().String()]).To(Equal(id))
			})
		})

		When("sender is not the only owner", func() {
			repoName := "repo"
			curHeight := uint64(0)
			proposalFee := util.String("1")
			id := "1"

			BeforeEach(func() {
				repo.AddOwner(sender.Addr().String(), &state.RepoOwner{})
				repo.AddOwner(key2.Addr().String(), &state.RepoOwner{})

				err = mergerequest.NewContract(&mergerequest.MergeRequestData{
					Repo:             repo,
					RepoName:         repoName,
					ProposalID:       id,
					ProposerFee:      proposalFee,
					Fee:              "1.5",
					CreatorAddress:   sender.Addr(),
					BaseBranch:       "base",
					BaseBranchHash:   "baseHash",
					TargetBranch:     "target",
					TargetBranchHash: "targetHash",
				}).Init(logic, nil, curHeight).Exec()
				Expect(err).To(BeNil())
			})

			It("should add the new proposal to the repo", func() {
				Expect(repo.Proposals).To(HaveLen(1))
			})

			Specify("that the proposal is not finalized or self accepted", func() {
				Expect(repo.Proposals).To(HaveLen(1))
				propID := mergerequest.MakeMergeRequestID(id)
				Expect(repo.Proposals.Get(propID).IsFinalized()).To(BeFalse())
				Expect(repo.Proposals.Get(propID).Yes).To(Equal(float64(0)))
			})

			Specify("that no fee was deducted", func() {
				acct := logic.AccountKeeper().Get(sender.Addr(), curHeight)
				Expect(acct.Balance.String()).To(Equal("10"))
			})

			Specify("that the proposal fee by the sender is registered on the proposal", func() {
				propID := mergerequest.MakeMergeRequestID(id)
				Expect(repo.Proposals.Get(propID).Fees).To(HaveLen(1))
				Expect(repo.Proposals.Get(propID).Fees).To(HaveKey(sender.Addr().String()))
				Expect(repo.Proposals.Get(propID).Fees[sender.Addr().String()]).To(Equal(id))
			})

			Specify("that the proposal was indexed against its end height", func() {
				res := logic.RepoKeeper().GetProposalsEndingAt(repo.Config.Governance.ProposalDuration + curHeight + 1)
				Expect(res).To(HaveLen(1))
			})
		})

		When("the proposal already exist and is not finalized", func() {
			repoName := "repo"
			curHeight := uint64(0)
			proposalFee := util.String("1")
			id := "1"

			BeforeEach(func() {
				repo.AddOwner(sender.Addr().String(), &state.RepoOwner{})
				repo.Proposals.Add(mergerequest.MakeMergeRequestID(id), &state.RepoProposal{
					ActionData: map[string][]byte{
						constants.ActionDataKeyBaseBranch:   []byte("base"),
						constants.ActionDataKeyBaseHash:     []byte("baseHash"),
						constants.ActionDataKeyTargetBranch: []byte("target"),
						constants.ActionDataKeyTargetHash:   []byte("targetHash"),
					},
				})

				err = mergerequest.NewContract(&mergerequest.MergeRequestData{
					Repo:             repo,
					RepoName:         repoName,
					ProposalID:       id,
					ProposerFee:      proposalFee,
					Fee:              "1.5",
					CreatorAddress:   sender.Addr(),
					BaseBranch:       "base2",
					BaseBranchHash:   "baseHash2",
					TargetBranch:     "target2",
					TargetBranchHash: "targetHash2",
				}).Init(logic, nil, curHeight).Exec()
				Expect(err).To(BeNil())
			})

			It("should not add a new proposal to the repo", func() {
				Expect(repo.Proposals).To(HaveLen(1))
			})

			It("should update proposal action data", func() {
				Expect(repo.Proposals).To(HaveLen(1))
				id := mergerequest.MakeMergeRequestID(id)
				Expect(repo.Proposals.Get(id).ActionData[constants.ActionDataKeyBaseBranch]).To(Equal([]byte("base2")))
				Expect(repo.Proposals.Get(id).ActionData[constants.ActionDataKeyBaseHash]).To(Equal([]byte("baseHash2")))
				Expect(repo.Proposals.Get(id).ActionData[constants.ActionDataKeyTargetBranch]).To(Equal([]byte("target2")))
				Expect(repo.Proposals.Get(id).ActionData[constants.ActionDataKeyTargetHash]).To(Equal([]byte("targetHash2")))
			})
		})

		When("the proposal already exist and is finalized", func() {
			repoName := "repo"
			curHeight := uint64(0)
			proposalFee := util.String("1")
			id := "1"

			BeforeEach(func() {
				repo.AddOwner(sender.Addr().String(), &state.RepoOwner{})
				repo.Proposals.Add(mergerequest.MakeMergeRequestID(id), &state.RepoProposal{
					Outcome: state.ProposalOutcomeAccepted,
					ActionData: map[string][]byte{
						constants.ActionDataKeyBaseBranch:   []byte("base"),
						constants.ActionDataKeyBaseHash:     []byte("baseHash"),
						constants.ActionDataKeyTargetBranch: []byte("target"),
						constants.ActionDataKeyTargetHash:   []byte("targetHash"),
					},
				})

				err = mergerequest.NewContract(&mergerequest.MergeRequestData{
					Repo:             repo,
					RepoName:         repoName,
					ProposalID:       id,
					ProposerFee:      proposalFee,
					Fee:              "1.5",
					CreatorAddress:   sender.Addr(),
					BaseBranch:       "base2",
					BaseBranchHash:   "baseHash2",
					TargetBranch:     "target2",
					TargetBranchHash: "targetHash2",
				}).Init(logic, nil, curHeight).Exec()
				Expect(err).To(BeNil())
			})

			It("should not add a new proposal to the repo", func() {
				Expect(repo.Proposals).To(HaveLen(1))
			})

			It("should not update proposal action data", func() {
				Expect(repo.Proposals).To(HaveLen(1))
				id := mergerequest.MakeMergeRequestID(id)
				Expect(repo.Proposals.Get(id).ActionData[constants.ActionDataKeyBaseBranch]).To(Equal([]byte("base")))
				Expect(repo.Proposals.Get(id).ActionData[constants.ActionDataKeyBaseHash]).To(Equal([]byte("baseHash")))
				Expect(repo.Proposals.Get(id).ActionData[constants.ActionDataKeyTargetBranch]).To(Equal([]byte("target")))
				Expect(repo.Proposals.Get(id).ActionData[constants.ActionDataKeyTargetHash]).To(Equal([]byte("targetHash")))
			})
		})
	})
})