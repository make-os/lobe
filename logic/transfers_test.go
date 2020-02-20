package logic

import (
	types2 "gitlab.com/makeos/mosdef/logic/types"
	"os"

	"github.com/golang/mock/gomock"

	"gitlab.com/makeos/mosdef/crypto"

	"gitlab.com/makeos/mosdef/util"

	"gitlab.com/makeos/mosdef/types"

	"gitlab.com/makeos/mosdef/config"
	"gitlab.com/makeos/mosdef/storage"
	"gitlab.com/makeos/mosdef/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transfers", func() {
	var appDB, stateTreeDB storage.Engine
	var err error
	var cfg *config.AppConfig
	var logic *Logic
	var txLogic *Transaction
	var ctrl *gomock.Controller

	BeforeEach(func() {
		cfg, err = testutil.SetTestCfg()
		Expect(err).To(BeNil())
		appDB, stateTreeDB = testutil.GetDB(cfg)
		logic = New(appDB, stateTreeDB, cfg)
		txLogic = &Transaction{logic: logic}
	})

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	BeforeEach(func() {
		err := logic.SysKeeper().SaveBlockInfo(&types2.BlockInfo{Height: 1})
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		Expect(appDB.Close()).To(BeNil())
		Expect(stateTreeDB.Close()).To(BeNil())
		err = os.RemoveAll(cfg.DataDir())
		Expect(err).To(BeNil())
	})

	Describe("CanExecCoinTransfer", func() {
		var sender = crypto.NewKeyFromIntSeed(1)

		Context("when sender variable type is not valid", func() {
			It("should panic", func() {
				Expect(func() {
					txLogic.CanExecCoinTransfer(123, util.String("100"), util.String("0"), 3, 1)
				}).Should(Panic())
			})
		})

		Context("when sender account has insufficient spendable balance", func() {
			It("should not return err='sender's spendable account balance is insufficient'", func() {
				err := txLogic.CanExecCoinTransfer(sender.PubKey(), util.String("100"), util.String("0"), 1, 1)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("field:value, error:sender's spendable account balance is insufficient"))
			})
		})

		Context("when nonce is invalid", func() {
			It("should return no error", func() {
				err := txLogic.CanExecCoinTransfer(sender.PubKey(), util.String("100"), util.String("0"), 3, 1)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("field:value, error:tx has invalid nonce (3), expected (1)"))
			})
		})

		Context("when sender account has sufficient spendable balance", func() {
			BeforeEach(func() {
				logic.AccountKeeper().Update(sender.Addr(), &types.Account{
					Balance: util.String("1000"),
					Stakes:  types.BareAccountStakes(),
				})
			})

			It("should return no error", func() {
				err := txLogic.CanExecCoinTransfer(sender.PubKey(), util.String("100"), util.String("0"), 1, 0)
				Expect(err).To(BeNil())
			})
		})
	})

	Describe(".execCoinTransfer", func() {
		var sender = crypto.NewKeyFromIntSeed(1)
		var recipientKey = crypto.NewKeyFromIntSeed(2)

		Context("when sender has bal=100, recipient has bal=10", func() {
			BeforeEach(func() {
				logic.AccountKeeper().Update(sender.Addr(), &types.Account{
					Balance: util.String("100"),
					Stakes:  types.BareAccountStakes(),
				})
				logic.AccountKeeper().Update(recipientKey.Addr(), &types.Account{
					Balance: util.String("10"),
					Stakes:  types.BareAccountStakes(),
				})
			})

			Context("sender creates a tx with value=10, fee=1", func() {
				BeforeEach(func() {
					senderPubKey := sender.PubKey().MustBytes32()
					err := txLogic.execCoinTransfer(senderPubKey, recipientKey.Addr(), util.String("10"), util.String("1"), 0)
					Expect(err).To(BeNil())
				})

				Specify("that sender balance is equal to 89 and nonce=1", func() {
					senderAcct := logic.AccountKeeper().GetAccount(sender.Addr())
					Expect(senderAcct.Balance).To(Equal(util.String("89")))
					Expect(senderAcct.Nonce).To(Equal(uint64(1)))
				})

				Specify("that recipient balance is equal to 20 and nonce=0", func() {
					recipientAcct := logic.AccountKeeper().GetAccount(recipientKey.Addr())
					Expect(recipientAcct.Balance).To(Equal(util.String("20")))
					Expect(recipientAcct.Nonce).To(Equal(uint64(0)))
				})
			})
		})

		Context("when sender and recipient are the same; with bal=100", func() {
			BeforeEach(func() {
				logic.AccountKeeper().Update(sender.Addr(), &types.Account{
					Balance: util.String("100"),
					Stakes:  types.BareAccountStakes(),
				})
			})

			Context("sender creates a tx with value=10, fee=1", func() {
				BeforeEach(func() {
					senderPubKey := sender.PubKey().MustBytes32()
					err := txLogic.execCoinTransfer(senderPubKey, sender.Addr(), util.String("10"), util.String("1"), 0)
					Expect(err).To(BeNil())
				})

				Specify("that sender balance is equal to 99 and nonce=1", func() {
					senderAcct := logic.AccountKeeper().GetAccount(sender.Addr())
					Expect(senderAcct.Balance).To(Equal(util.String("99")))
					Expect(senderAcct.Nonce).To(Equal(uint64(1)))
				})
			})
		})

		When("recipient address is a namespaced URI with a user account target", func() {
			ns := "namespace"
			var senderNamespaceURI = util.String(ns + "/domain")

			BeforeEach(func() {
				logic.AccountKeeper().Update(sender.Addr(), &types.Account{
					Balance: util.String("100"),
					Stakes:  types.BareAccountStakes(),
				})

				logic.NamespaceKeeper().Update(util.Hash20Hex([]byte(ns)), &types.Namespace{
					Domains: map[string]string{
						"domain": "a/" + sender.Addr().String(),
					},
				})
			})

			Context("sender creates a tx with value=10, fee=1", func() {
				BeforeEach(func() {
					senderPubKey := sender.PubKey().MustBytes32()
					err := txLogic.execCoinTransfer(senderPubKey, senderNamespaceURI, util.String("10"), util.String("1"), 0)
					Expect(err).To(BeNil())
				})

				Specify("that sender balance is equal to 99 and nonce=1", func() {
					senderAcct := logic.AccountKeeper().GetAccount(sender.Addr())
					Expect(senderAcct.Balance).To(Equal(util.String("99")))
					Expect(senderAcct.Nonce).To(Equal(uint64(1)))
				})
			})
		})

		When("recipient address is a prefixed user account address", func() {
			BeforeEach(func() {
				logic.AccountKeeper().Update(sender.Addr(), &types.Account{
					Balance: util.String("100"),
					Stakes:  types.BareAccountStakes(),
				})
			})

			Context("sender creates a tx with value=10, fee=1", func() {
				BeforeEach(func() {
					recipient := "a/" + sender.Addr()
					senderPubKey := sender.PubKey().MustBytes32()
					err := txLogic.execCoinTransfer(senderPubKey, recipient, util.String("10"), util.String("1"), 0)
					Expect(err).To(BeNil())
				})

				Specify("that sender balance is equal to 99 and nonce=1", func() {
					senderAcct := logic.AccountKeeper().GetAccount(sender.Addr())
					Expect(senderAcct.Balance).To(Equal(util.String("99")))
					Expect(senderAcct.Nonce).To(Equal(uint64(1)))
				})
			})
		})

		When("recipient address is a namespaced URI with a repo account target", func() {
			ns := "namespace"
			var senderNamespaceURI = util.String(ns + "/domain")
			var repoName = "repo1"

			BeforeEach(func() {
				logic.AccountKeeper().Update(sender.Addr(), &types.Account{
					Balance: util.String("100"),
					Stakes:  types.BareAccountStakes(),
				})

				logic.NamespaceKeeper().Update(util.Hash20Hex([]byte(ns)), &types.Namespace{
					Domains: map[string]string{
						"domain": "r/" + repoName,
					},
				})
			})

			Context("sender creates a tx with value=10, fee=1", func() {
				BeforeEach(func() {
					senderPubKey := sender.PubKey().MustBytes32()
					err := txLogic.execCoinTransfer(senderPubKey, senderNamespaceURI, util.String("10"), util.String("1"), 0)
					Expect(err).To(BeNil())
				})

				Specify("that sender balance is equal to 89 and nonce=1", func() {
					senderAcct := logic.AccountKeeper().GetAccount(sender.Addr())
					Expect(senderAcct.Balance).To(Equal(util.String("89")))
					Expect(senderAcct.Nonce).To(Equal(uint64(1)))
				})

				Specify("that the repo has a balance=10", func() {
					repo := logic.RepoKeeper().GetRepo(repoName)
					Expect(repo.GetBalance()).To(Equal(util.String("10")))
				})
			})
		})

		When("recipient address is a prefixed user account address", func() {
			BeforeEach(func() {
				logic.AccountKeeper().Update(sender.Addr(), &types.Account{
					Balance: util.String("100"),
					Stakes:  types.BareAccountStakes(),
				})
			})

			Context("sender creates a tx with value=10, fee=1", func() {
				var repoName = "repo1"

				BeforeEach(func() {
					recipient := util.String("r/" + repoName)
					senderPubKey := sender.PubKey().MustBytes32()
					err := txLogic.execCoinTransfer(senderPubKey, recipient, util.String("10"), util.String("1"), 0)
					Expect(err).To(BeNil())
				})

				Specify("that sender balance is equal to 89 and nonce=1", func() {
					senderAcct := logic.AccountKeeper().GetAccount(sender.Addr())
					Expect(senderAcct.Balance).To(Equal(util.String("89")))
					Expect(senderAcct.Nonce).To(Equal(uint64(1)))
				})

				Specify("that the repo has a balance=10", func() {
					repo := logic.RepoKeeper().GetRepo(repoName)
					Expect(repo.GetBalance()).To(Equal(util.String("10")))
				})
			})
		})
	})
})
