package logic

import (
	"os"

	"github.com/golang/mock/gomock"

	"github.com/makeos/mosdef/crypto"
	"github.com/makeos/mosdef/types"
	"github.com/makeos/mosdef/types/mocks"
	"github.com/makeos/mosdef/util"

	"github.com/makeos/mosdef/config"
	"github.com/makeos/mosdef/storage"
	"github.com/makeos/mosdef/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {
	var appDB, stateTreeDB storage.Engine
	var err error
	var cfg *config.AppConfig
	var logic *Logic
	var txLogic *Transaction
	var ctrl *gomock.Controller
	var mockTickMgr *mocks.MockTicketManager
	var sender = crypto.NewKeyFromIntSeed(1)
	var key2 = crypto.NewKeyFromIntSeed(2)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		cfg, err = testutil.SetTestCfg()
		Expect(err).To(BeNil())
		appDB, stateTreeDB = testutil.GetDB(cfg)
		logic = New(appDB, stateTreeDB, cfg)
		txLogic = &Transaction{logic: logic}
		mockTickMgr = mocks.NewMockTicketManager(ctrl)
	})

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	BeforeEach(func() {
		types.DefaultRepoConfig = types.MakeDefaultRepoConfig()
		err := logic.SysKeeper().SaveBlockInfo(&types.BlockInfo{Height: 1})
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		Expect(appDB.Close()).To(BeNil())
		Expect(stateTreeDB.Close()).To(BeNil())
		err = os.RemoveAll(cfg.DataDir())
		Expect(err).To(BeNil())
	})

	Describe(".execRepoCreate", func() {
		var err error
		var sender = crypto.NewKeyFromIntSeed(1)
		var spk util.Bytes32

		BeforeEach(func() {
			logic.AccountKeeper().Update(sender.Addr(), &types.Account{
				Balance:             util.String("10"),
				Stakes:              types.BareAccountStakes(),
				DelegatorCommission: 10,
			})
		})

		When("successful", func() {
			BeforeEach(func() {
				types.DefaultRepoConfig.Governace.ProposalProposee = types.ProposeeOwner
				spk = sender.PubKey().MustBytes32()
				err = txLogic.execRepoCreate(spk, "repo", "1.5", 0)
				Expect(err).To(BeNil())
			})

			Specify("that the repo was added to the tree", func() {
				repo := txLogic.logic.RepoKeeper().GetRepo("repo")
				Expect(repo.IsNil()).To(BeFalse())
				Expect(repo.Owners).To(HaveKey(sender.Addr().String()))
			})

			Specify("that fee is deducted from sender account", func() {
				acct := logic.AccountKeeper().GetAccount(sender.Addr())
				Expect(acct.GetBalance()).To(Equal(util.String("8.5")))
			})

			Specify("that sender account nonce increased", func() {
				acct := logic.AccountKeeper().GetAccount(sender.Addr())
				Expect(acct.Nonce).To(Equal(uint64(1)))
			})

			When("proposee is not ProposalOwner", func() {
				BeforeEach(func() {
					types.DefaultRepoConfig.Governace.ProposalProposee = types.ProposeeNetStakeholders
					spk = sender.PubKey().MustBytes32()
					err = txLogic.execRepoCreate(spk, "repo", "1.5", 0)
					Expect(err).To(BeNil())
				})

				It("should not add the sender as an owner", func() {
					repo := txLogic.logic.RepoKeeper().GetRepo("repo")
					Expect(repo.Owners).To(BeEmpty())
				})
			})
		})
	})

	Describe(".execRepoProposalVote", func() {
		var err error
		var spk util.Bytes32

		BeforeEach(func() {
			logic.AccountKeeper().Update(sender.Addr(), &types.Account{
				Balance:             util.String("10"),
				Stakes:              types.BareAccountStakes(),
				DelegatorCommission: 10,
			})
		})

		When("proposal tally method is ProposalTallyOneVote", func() {
			var propID = "proposer_id"
			var repoName = "repo"

			BeforeEach(func() {
				repoUpd := types.BareRepository()
				repoUpd.Config = types.DefaultRepoConfig
				repoUpd.AddOwner(sender.Addr().String(), &types.RepoOwner{})
				proposal := &types.RepoProposal{
					Proposee:    types.ProposeeOwner,
					TallyMethod: types.ProposalTallyOneVote,
					Yes:         1,
				}
				repoUpd.Proposals.Add(propID, proposal)
				logic.RepoKeeper().Update(repoName, repoUpd)

				spk = sender.PubKey().MustBytes32()
				err = txLogic.execRepoProposalVote(spk, repoName, propID, types.ProposalVoteYes, "1.5", 0)
				Expect(err).To(BeNil())
			})

			It("should increment proposal.Yes by 1", func() {
				repo := logic.RepoKeeper().GetRepo(repoName)
				Expect(repo.Proposals.Get(propID).Yes).To(Equal(float64(2)))
			})
		})

		When("proposal tally method is ProposalTallyMethodCoinWeighted", func() {
			var propID = "proposer_id"
			var repoName = "repo"

			BeforeEach(func() {
				repoUpd := types.BareRepository()
				repoUpd.Config = types.DefaultRepoConfig
				repoUpd.AddOwner(sender.Addr().String(), &types.RepoOwner{})
				proposal := &types.RepoProposal{
					Proposee:    types.ProposeeOwner,
					TallyMethod: types.ProposalTallyMethodCoinWeighted,
					Yes:         1,
				}
				repoUpd.Proposals.Add(propID, proposal)
				logic.RepoKeeper().Update(repoName, repoUpd)

				spk = sender.PubKey().MustBytes32()
				err = txLogic.execRepoProposalVote(spk, repoName, propID, types.ProposalVoteYes, "1.5", 0)
				Expect(err).To(BeNil())
			})

			It("should increment proposal.Yes by 10", func() {
				repo := logic.RepoKeeper().GetRepo(repoName)
				Expect(repo.Proposals.Get(propID).Yes).To(Equal(float64(11)))
			})
		})

		When("proposal tally method is ProposalTallyMethodNetStakeOfProposer and the voter's non-delegated ticket value=100", func() {
			var propID = "proposer_id"
			var repoName = "repo"

			BeforeEach(func() {
				repoUpd := types.BareRepository()
				repoUpd.Config = types.DefaultRepoConfig
				repoUpd.AddOwner(sender.Addr().String(), &types.RepoOwner{})
				proposal := &types.RepoProposal{
					TallyMethod: types.ProposalTallyMethodNetStakeOfProposer,
					Yes:         0,
				}
				repoUpd.Proposals.Add(propID, proposal)
				logic.RepoKeeper().Update(repoName, repoUpd)

				mockTickMgr.EXPECT().ValueOfNonDelegatedTickets(sender.PubKey().
					MustBytes32(), uint64(0)).Return(float64(100), nil)
				txLogic.logic.SetTicketManager(mockTickMgr)

				spk = sender.PubKey().MustBytes32()
				err = txLogic.execRepoProposalVote(spk, repoName, propID, types.ProposalVoteYes, "1.5", 10)
				Expect(err).To(BeNil())
			})

			It("should increment proposal.Yes by 100", func() {
				repo := logic.RepoKeeper().GetRepo(repoName)
				Expect(repo.Proposals.Get(propID).Yes).To(Equal(float64(100)))
			})
		})

		When("proposal tally method is ProposalTallyMethodNetStakeOfDelegators and the voter's non-delegated ticket value=100", func() {
			var propID = "proposer_id"
			var repoName = "repo"

			BeforeEach(func() {
				repoUpd := types.BareRepository()
				repoUpd.Config = types.DefaultRepoConfig
				repoUpd.AddOwner(sender.Addr().String(), &types.RepoOwner{})
				proposal := &types.RepoProposal{
					TallyMethod: types.ProposalTallyMethodNetStakeOfDelegators,
					Yes:         0,
				}
				repoUpd.Proposals.Add(propID, proposal)
				logic.RepoKeeper().Update(repoName, repoUpd)

				mockTickMgr.EXPECT().ValueOfDelegatedTickets(sender.PubKey().
					MustBytes32(), uint64(0)).Return(float64(100), nil)
				txLogic.logic.SetTicketManager(mockTickMgr)

				spk = sender.PubKey().MustBytes32()
				err = txLogic.execRepoProposalVote(spk, repoName, propID, types.ProposalVoteYes, "1.5", 10)
				Expect(err).To(BeNil())
			})

			It("should increment proposal.Yes by 100", func() {
				repo := logic.RepoKeeper().GetRepo(repoName)
				Expect(repo.Proposals.Get(propID).Yes).To(Equal(float64(100)))
			})
		})

		When("proposal tally method is ProposalTallyMethodNetStake", func() {
			var propID = "proposer_id"
			var repoName = "repo"

			When("ticketA and ticketB are not delegated, with value 10, 20 respectively", func() {
				BeforeEach(func() {
					repoUpd := types.BareRepository()
					repoUpd.Config = types.DefaultRepoConfig
					repoUpd.AddOwner(sender.Addr().String(), &types.RepoOwner{})
					proposal := &types.RepoProposal{
						TallyMethod: types.ProposalTallyMethodNetStake,
						Yes:         0,
					}
					repoUpd.Proposals.Add(propID, proposal)
					logic.RepoKeeper().Update(repoName, repoUpd)

					ticketA := &types.Ticket{Value: "10"}
					ticketB := &types.Ticket{Value: "20"}
					tickets := []*types.Ticket{ticketA, ticketB}

					mockTickMgr.EXPECT().GetNonDecayedTickets(sender.PubKey().
						MustBytes32(), uint64(0)).Return(tickets, nil)
					txLogic.logic.SetTicketManager(mockTickMgr)

					spk = sender.PubKey().MustBytes32()
					err = txLogic.execRepoProposalVote(spk, repoName, propID, types.ProposalVoteYes, "1.5", 10)
					Expect(err).To(BeNil())
				})

				It("should increment proposal.Yes by 30", func() {
					repo := logic.RepoKeeper().GetRepo(repoName)
					Expect(repo.Proposals.Get(propID).Yes).To(Equal(float64(30)))
				})
			})

			When("ticketA and ticketB exist, with value 10, 20 respectively. voter is delegator and proposer of ticketB", func() {
				BeforeEach(func() {
					repoUpd := types.BareRepository()
					repoUpd.Config = types.DefaultRepoConfig
					repoUpd.AddOwner(sender.Addr().String(), &types.RepoOwner{})
					proposal := &types.RepoProposal{
						TallyMethod: types.ProposalTallyMethodNetStake,
						Yes:         0,
					}
					repoUpd.Proposals.Add(propID, proposal)
					logic.RepoKeeper().Update(repoName, repoUpd)

					ticketA := &types.Ticket{Value: "10"}
					ticketB := &types.Ticket{
						Value:          "20",
						ProposerPubKey: sender.PubKey().MustBytes32(),
						Delegator:      sender.Addr().String(),
					}
					tickets := []*types.Ticket{ticketA, ticketB}

					mockTickMgr.EXPECT().GetNonDecayedTickets(sender.PubKey().
						MustBytes32(), uint64(0)).Return(tickets, nil)
					txLogic.logic.SetTicketManager(mockTickMgr)

					spk = sender.PubKey().MustBytes32()
					err = txLogic.execRepoProposalVote(spk, repoName, propID, types.ProposalVoteYes, "1.5", 10)
					Expect(err).To(BeNil())
				})

				It("should increment proposal.Yes by 30", func() {
					repo := logic.RepoKeeper().GetRepo(repoName)
					Expect(repo.Proposals.Get(propID).Yes).To(Equal(float64(30)))
				})
			})

			When("ticketA and ticketB exist, with value 10, 20 respectively. voter is "+
				"proposer of ticketB but someone else is delegator and they have not "+
				"voted on the proposal", func() {
				BeforeEach(func() {
					repoUpd := types.BareRepository()
					repoUpd.Config = types.DefaultRepoConfig
					repoUpd.AddOwner(sender.Addr().String(), &types.RepoOwner{})
					proposal := &types.RepoProposal{
						TallyMethod: types.ProposalTallyMethodNetStake,
						Yes:         0,
					}
					repoUpd.Proposals.Add(propID, proposal)
					logic.RepoKeeper().Update(repoName, repoUpd)

					ticketA := &types.Ticket{Value: "10"}
					ticketB := &types.Ticket{
						Value:          "20",
						ProposerPubKey: sender.PubKey().MustBytes32(),
						Delegator:      key2.Addr().String(),
					}
					tickets := []*types.Ticket{ticketA, ticketB}

					mockTickMgr.EXPECT().GetNonDecayedTickets(sender.PubKey().
						MustBytes32(), uint64(0)).Return(tickets, nil)
					txLogic.logic.SetTicketManager(mockTickMgr)

					spk = sender.PubKey().MustBytes32()
					err = txLogic.execRepoProposalVote(spk, repoName, propID, types.ProposalVoteYes, "1.5", 10)
					Expect(err).To(BeNil())
				})

				It("should increment proposal.Yes by 30", func() {
					repo := logic.RepoKeeper().GetRepo(repoName)
					Expect(repo.Proposals.Get(propID).Yes).To(Equal(float64(30)))
				})
			})

			When("ticketA and ticketB exist, with value 10, 20 respectively. voter is "+
				"proposer of ticketB but someone else is delegator and they have "+
				"voted on the proposal", func() {
				BeforeEach(func() {
					repoUpd := types.BareRepository()
					repoUpd.Config = types.DefaultRepoConfig
					repoUpd.AddOwner(sender.Addr().String(), &types.RepoOwner{})
					proposal := &types.RepoProposal{
						TallyMethod: types.ProposalTallyMethodNetStake,
						Yes:         0,
					}
					repoUpd.Proposals.Add(propID, proposal)
					logic.RepoKeeper().Update(repoName, repoUpd)

					ticketA := &types.Ticket{Value: "10"}
					ticketB := &types.Ticket{
						Value:          "20",
						ProposerPubKey: sender.PubKey().MustBytes32(),
						Delegator:      key2.Addr().String(),
					}
					tickets := []*types.Ticket{ticketA, ticketB}

					mockTickMgr.EXPECT().GetNonDecayedTickets(sender.PubKey().
						MustBytes32(), uint64(0)).Return(tickets, nil)
					txLogic.logic.SetTicketManager(mockTickMgr)

					logic.RepoKeeper().IndexProposalVote(repoName, propID,
						key2.Addr().String(), types.ProposalVoteYes)

					spk = sender.PubKey().MustBytes32()
					err = txLogic.execRepoProposalVote(spk, repoName, propID, types.ProposalVoteYes, "1.5", 10)
					Expect(err).To(BeNil())
				})

				It("should increment proposal.Yes by 10", func() {
					repo := logic.RepoKeeper().GetRepo(repoName)
					Expect(repo.Proposals.Get(propID).Yes).To(Equal(float64(10)))
				})
			})

			When("ticketA and ticketB exist, with value 10, 20 respectively. voter is "+
				"delegator of ticketB but someone else is proposer and they have not "+
				"voted on the proposal", func() {
				BeforeEach(func() {
					repoUpd := types.BareRepository()
					repoUpd.Config = types.DefaultRepoConfig
					repoUpd.AddOwner(sender.Addr().String(), &types.RepoOwner{})
					proposal := &types.RepoProposal{
						TallyMethod: types.ProposalTallyMethodNetStake,
						Yes:         0,
					}
					repoUpd.Proposals.Add(propID, proposal)
					logic.RepoKeeper().Update(repoName, repoUpd)

					ticketA := &types.Ticket{Value: "10"}
					ticketB := &types.Ticket{
						Value:          "20",
						ProposerPubKey: key2.PubKey().MustBytes32(),
						Delegator:      sender.Addr().String(),
					}
					tickets := []*types.Ticket{ticketA, ticketB}

					mockTickMgr.EXPECT().GetNonDecayedTickets(sender.PubKey().
						MustBytes32(), uint64(0)).Return(tickets, nil)
					txLogic.logic.SetTicketManager(mockTickMgr)

					spk = sender.PubKey().MustBytes32()
					err = txLogic.execRepoProposalVote(spk, repoName, propID, types.ProposalVoteYes, "1.5", 10)
					Expect(err).To(BeNil())
				})

				It("should increment proposal.Yes by 30", func() {
					repo := logic.RepoKeeper().GetRepo(repoName)
					Expect(repo.Proposals.Get(propID).Yes).To(Equal(float64(30)))
				})
			})

			When("ticketA and ticketB exist, with value 10, 20 respectively. voter is "+
				"delegator of ticketB but someone else is proposer and they have "+
				"voted 'Yes' on the proposal", func() {
				BeforeEach(func() {
					repoUpd := types.BareRepository()
					repoUpd.Config = types.DefaultRepoConfig
					repoUpd.AddOwner(sender.Addr().String(), &types.RepoOwner{})
					proposal := &types.RepoProposal{
						TallyMethod: types.ProposalTallyMethodNetStake,
						Yes:         100,
					}
					repoUpd.Proposals.Add(propID, proposal)
					logic.RepoKeeper().Update(repoName, repoUpd)

					ticketA := &types.Ticket{Value: "10"}
					ticketB := &types.Ticket{
						Value:          "20",
						ProposerPubKey: key2.PubKey().MustBytes32(),
						Delegator:      sender.Addr().String(),
					}
					tickets := []*types.Ticket{ticketA, ticketB}

					mockTickMgr.EXPECT().GetNonDecayedTickets(sender.PubKey().
						MustBytes32(), uint64(0)).Return(tickets, nil)
					txLogic.logic.SetTicketManager(mockTickMgr)

					logic.RepoKeeper().IndexProposalVote(repoName, propID,
						key2.Addr().String(), types.ProposalVoteYes)

					spk = sender.PubKey().MustBytes32()
					err = txLogic.execRepoProposalVote(spk, repoName, propID, types.ProposalVoteNo, "1.5", 10)
					Expect(err).To(BeNil())
				})

				It("should increment proposal.No by 30", func() {
					repo := logic.RepoKeeper().GetRepo(repoName)
					Expect(repo.Proposals.Get(propID).No).To(Equal(float64(30)))
				})

				Specify("that proposal.Yes is now 80", func() {
					repo := logic.RepoKeeper().GetRepo(repoName)
					Expect(repo.Proposals.Get(propID).Yes).To(Equal(float64(80)))
				})
			})
		})
	})

	Describe(".execRepoUpsertOwner", func() {
		var err error
		var sender = crypto.NewKeyFromIntSeed(1)
		var key2 = crypto.NewKeyFromIntSeed(2)
		var spk util.Bytes32
		var repoUpd *types.Repository

		BeforeEach(func() {
			logic.AccountKeeper().Update(sender.Addr(), &types.Account{
				Balance:             util.String("10"),
				Stakes:              types.BareAccountStakes(),
				DelegatorCommission: 10,
			})
			repoUpd = types.BareRepository()
			repoUpd.Config = types.DefaultRepoConfig
			repoUpd.Config.Governace.ProposalProposee = types.ProposeeOwner
		})

		When("sender is the only owner", func() {
			repoName := "repo"
			address := "owner_address"

			BeforeEach(func() {
				repoUpd.AddOwner(sender.Addr().String(), &types.RepoOwner{})
				logic.RepoKeeper().Update(repoName, repoUpd)

				spk = sender.PubKey().MustBytes32()
				err = txLogic.execRepoUpsertOwner(spk, repoName, address, false, "1.5", 0)
				Expect(err).To(BeNil())
			})

			It("should add the new proposal to the repo", func() {
				repo := logic.RepoKeeper().GetRepo(repoName)
				Expect(repo.Proposals).To(HaveLen(1))
			})

			Specify("that the proposal is finalized and self accepted", func() {
				repo := logic.RepoKeeper().GetRepo(repoName)
				Expect(repo.Proposals).To(HaveLen(1))
				Expect(repo.Proposals.Get("1").IsFinalized()).To(BeTrue())
				Expect(repo.Proposals.Get("1").Yes).To(Equal(float64(1)))
			})

			Specify("that new owner was added", func() {
				repo := logic.RepoKeeper().GetRepo(repoName)
				Expect(repo.Owners).To(HaveLen(2))
			})

			Specify("that fee was deducted", func() {
				acct := logic.AccountKeeper().GetAccount(sender.Addr(), 0)
				Expect(acct.Balance.String()).To(Equal("8.5"))
			})
		})

		When("sender is the only owner and there are multiple addresses", func() {
			repoName := "repo"
			addresses := "owner_address,owner_address2"

			BeforeEach(func() {
				repoUpd.AddOwner(sender.Addr().String(), &types.RepoOwner{})
				logic.RepoKeeper().Update(repoName, repoUpd)

				spk = sender.PubKey().MustBytes32()
				err = txLogic.execRepoUpsertOwner(spk, repoName, addresses, false, "1.5", 0)
				Expect(err).To(BeNil())
			})

			It("should add the new proposal to the repo", func() {
				repo := logic.RepoKeeper().GetRepo(repoName)
				Expect(repo.Proposals).To(HaveLen(1))
			})

			Specify("that the proposal is finalized and self accepted", func() {
				repo := logic.RepoKeeper().GetRepo(repoName)
				Expect(repo.Proposals).To(HaveLen(1))
				Expect(repo.Proposals.Get("1").IsFinalized()).To(BeTrue())
				Expect(repo.Proposals.Get("1").Yes).To(Equal(float64(1)))
				Expect(repo.Proposals.Get("1").Outcome).To(Equal(types.ProposalOutcomeAccepted))
			})

			Specify("that three owners were added", func() {
				repo := logic.RepoKeeper().GetRepo(repoName)
				Expect(repo.Owners).To(HaveLen(3))
			})

			Specify("that fee was deducted", func() {
				acct := logic.AccountKeeper().GetAccount(sender.Addr(), 0)
				Expect(acct.Balance.String()).To(Equal("8.5"))
			})
		})

		When("sender is not the only owner", func() {
			repoName := "repo"
			address := "owner_address"
			var curHeight = uint64(0)

			BeforeEach(func() {
				repoUpd.AddOwner(sender.Addr().String(), &types.RepoOwner{})
				repoUpd.AddOwner(key2.Addr().String(), &types.RepoOwner{})
				logic.RepoKeeper().Update(repoName, repoUpd)

				spk = sender.PubKey().MustBytes32()
				err = txLogic.execRepoUpsertOwner(spk, repoName, address, false, "1.5", curHeight)
				Expect(err).To(BeNil())
			})

			It("should add the new proposal to the repo", func() {
				repo := logic.RepoKeeper().GetRepo(repoName)
				Expect(repo.Proposals).To(HaveLen(1))
			})

			Specify("that the proposal is not finalized or self accepted", func() {
				repo := logic.RepoKeeper().GetRepo(repoName)
				Expect(repo.Proposals).To(HaveLen(1))
				Expect(repo.Proposals.Get("1").IsFinalized()).To(BeFalse())
				Expect(repo.Proposals.Get("1").Yes).To(Equal(float64(0)))
			})

			Specify("that fee was deducted", func() {
				acct := logic.AccountKeeper().GetAccount(sender.Addr(), curHeight)
				Expect(acct.Balance.String()).To(Equal("8.5"))
			})

			Specify("that the proposal was indexed against its end height", func() {
				res := logic.RepoKeeper().GetProposalsEndingAt(repoUpd.Config.Governace.ProposalDur + curHeight + 1)
				Expect(res).To(HaveLen(1))
			})
		})
	})
})
