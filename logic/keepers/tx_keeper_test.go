package keepers

import (
	"fmt"
	"os"

	"github.com/makeos/mosdef/types"

	storagemocks "github.com/makeos/mosdef/storage/mocks"

	"github.com/golang/mock/gomock"

	"github.com/makeos/mosdef/config"
	"github.com/makeos/mosdef/storage"
	"github.com/makeos/mosdef/testutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TxKeeper", func() {
	var c storage.Engine
	var err error
	var cfg *config.EngineConfig
	var txKeeper *TxKeeper
	var ctrl *gomock.Controller

	BeforeEach(func() {
		cfg, err = testutil.SetTestCfg()
		Expect(err).To(BeNil())
		c = storage.NewBadger(cfg)
		Expect(c.Init()).To(BeNil())
		txKeeper = NewTxKeeper(c)
	})

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	AfterEach(func() {
		Expect(c.Close()).To(BeNil())
		err = os.RemoveAll(cfg.DataDir())
		Expect(err).To(BeNil())
	})

	Describe(".Index", func() {
		When("db operation failed", func() {
			BeforeEach(func() {
				mockDB := storagemocks.NewMockEngine(ctrl)
				mockDB.EXPECT().Put(gomock.Any()).Return(fmt.Errorf("error"))
				txKeeper.db = mockDB
			})

			It("should return err='failed to index tx: error'", func() {
				tx := types.NewBareTx(types.TxTypeTransferCoin)
				err := txKeeper.Index(tx)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("failed to index tx: error"))
			})
		})

		When("index is successful", func() {
			tx := types.NewBareTx(types.TxTypeTransferCoin)

			BeforeEach(func() {
				err := txKeeper.Index(tx)
				Expect(err).To(BeNil())
			})

			It("should return nil", func() {
				rec, err := txKeeper.db.Get(MakeTxKey(tx.GetID()))
				Expect(err).To(BeNil())
				Expect(rec.Value).To(Equal(tx.Bytes()))
			})
		})
	})

	Describe(".GetTx", func() {
		When("db operation failed", func() {
			BeforeEach(func() {
				mockDB := storagemocks.NewMockEngine(ctrl)
				mockDB.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("error"))
				txKeeper.db = mockDB
			})

			It("should return err='failed to get tx: error'", func() {
				tx := types.NewBareTx(types.TxTypeTransferCoin)
				_, err := txKeeper.GetTx(tx.GetID())
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("failed to get tx: error"))
			})
		})

		When("tx is found", func() {
			tx := types.NewBareTx(types.TxTypeTransferCoin)

			BeforeEach(func() {
				err := txKeeper.Index(tx)
				Expect(err).To(BeNil())
			})

			It("should return tx", func() {
				res, err := txKeeper.GetTx(tx.GetID())
				Expect(err).To(BeNil())
				Expect(res.Bytes()).To(Equal(tx.Bytes()))
			})
		})

		When("tx is not found", func() {
			It("should return tx", func() {
				_, err := txKeeper.GetTx("unknown")
				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(ErrTxNotFound))
			})
		})
	})
})
