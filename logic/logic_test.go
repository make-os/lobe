package logic_test

import (
	"os"

	"github.com/makeos/mosdef/util"

	"github.com/makeos/mosdef/config"
	"github.com/makeos/mosdef/storage"
	"github.com/makeos/mosdef/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	l "github.com/makeos/mosdef/logic"
)

var _ = Describe("Logic", func() {
	var appDB, stateTreeDB storage.Engine
	var err error
	var cfg *config.EngineConfig
	var logic *l.Logic

	BeforeEach(func() {
		cfg, err = testutil.SetTestCfg()
		Expect(err).To(BeNil())
		appDB, stateTreeDB = testutil.GetDB(cfg)
		logic = l.New(appDB, stateTreeDB, cfg)
	})

	AfterEach(func() {
		Expect(appDB.Close()).To(BeNil())
		Expect(stateTreeDB.Close()).To(BeNil())
		err = os.RemoveAll(cfg.DataDir())
		Expect(err).To(BeNil())
	})

	Describe(".WriteGenesisState", func() {
		var testGenAccts = []*config.GenAccount{
			&config.GenAccount{Address: "addr1", Balance: "100"},
			&config.GenAccount{Address: "addr2", Balance: "200"},
		}

		BeforeEach(func() {
			cfg.GenesisAccounts = testGenAccts
			for _, a := range testGenAccts {
				res := logic.AccountKeeper().GetAccount(util.String(a.Address))
				Expect(res.Balance).To(Equal(util.String("0")))
				Expect(res.Nonce).To(Equal(uint64(0)))
			}
		})

		It("should successfully add genesis accounts", func() {
			err := logic.WriteGenesisState()
			Expect(err).To(BeNil())
			addr1Res := logic.AccountKeeper().GetAccount(util.String(testGenAccts[0].Address))
			Expect(addr1Res.Balance).To(Equal(util.String("100")))
			addr2Res := logic.AccountKeeper().GetAccount(util.String(testGenAccts[1].Address))
			Expect(addr2Res.Balance).To(Equal(util.String("200")))
		})
	})
})
