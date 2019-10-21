package keepers

import (
	"fmt"
	"os"

	"github.com/makeos/mosdef/params"

	"github.com/golang/mock/gomock"

	"github.com/makeos/mosdef/config"
	"github.com/makeos/mosdef/storage"
	storagemocks "github.com/makeos/mosdef/storage/mocks"
	"github.com/makeos/mosdef/testutil"
	"github.com/makeos/mosdef/types"
	"github.com/makeos/mosdef/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SystemKeeper", func() {
	var c storage.Engine
	var err error
	var cfg *config.EngineConfig
	var valKeeper *ValidatorKeeper
	var ctrl *gomock.Controller
	var pubKey = types.HexBytes([]byte("pubkey"))

	BeforeEach(func() {
		cfg, err = testutil.SetTestCfg()
		Expect(err).To(BeNil())
		c = storage.NewBadger(cfg)
		Expect(c.Init()).To(BeNil())
		valKeeper = NewValidatorKeeper(c.F(true, true))
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

	Describe(".getByHeight", func() {
		When("no result is found", func() {
			It("should return err=nil and empty map", func() {
				res, err := valKeeper.getByHeight(1)
				Expect(err).To(BeNil())
				Expect(res).To(BeEmpty())
			})
		})

		When("db error occurred", func() {
			BeforeEach(func() {
				mockDB := storagemocks.NewMockFunctions(ctrl)
				mockDB.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("error"))
				valKeeper.db = mockDB
			})

			It("should return err='error' and nil result", func() {
				res, err := valKeeper.getByHeight(1)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("error"))
				Expect(res).To(BeNil())
			})
		})

		When("record exist", func() {
			height := int64(1)
			rec := map[string]*types.Validator{"pubkey": &types.Validator{Power: 1}}
			BeforeEach(func() {
				key := MakeBlockValidatorsKey(height)
				err := valKeeper.db.Put(storage.NewFromKeyValue(key, util.ObjectToBytes(rec)))
				Expect(err).To(BeNil())
			})

			It("should return err=nil and expected result", func() {
				res, err := valKeeper.getByHeight(height)
				Expect(err).To(BeNil())
				Expect(res).To(Equal(types.BlockValidators(rec)))
			})
		})
	})

	Describe(".GetByHeight", func() {
		When("one validator is stored at height=1, search height = 1", func() {
			rec := map[string]*types.Validator{"pubkey": &types.Validator{Power: 1}}
			BeforeEach(func() {
				key := MakeBlockValidatorsKey(1)
				err := valKeeper.db.Put(storage.NewFromKeyValue(key, util.ObjectToBytes(rec)))
				Expect(err).To(BeNil())
			})

			It("should return err=nil and one validator", func() {
				res, err := valKeeper.GetByHeight(1)
				Expect(err).To(BeNil())
				Expect(res).To(Equal(types.BlockValidators(rec)))
			})
		})

		When("two two validators exist; valset1 at height 1, valset2 at height 2; argument height = 0", func() {
			valset := map[string]*types.Validator{"pubkey": &types.Validator{Power: 1}}
			valset2 := map[string]*types.Validator{"pubkey": &types.Validator{Power: 2}}
			BeforeEach(func() {
				err := valKeeper.db.Put(storage.NewFromKeyValue(MakeBlockValidatorsKey(1), util.ObjectToBytes(valset)))
				Expect(err).To(BeNil())
				err = valKeeper.db.Put(storage.NewFromKeyValue(MakeBlockValidatorsKey(2), util.ObjectToBytes(valset2)))
				Expect(err).To(BeNil())
			})

			It("should return valset2 since it is the most recent", func() {
				res, err := valKeeper.GetByHeight(0)
				Expect(err).To(BeNil())
				Expect(res).To(Equal(types.BlockValidators(valset2)))
			})
		})

		When("two validators exist; valset1 at height 2, valset2 at height 4; argument height = 9; blocks per epoch = 2", func() {
			valset := map[string]*types.Validator{"pubkey": &types.Validator{Power: 1}}
			valset2 := map[string]*types.Validator{"pubkey": &types.Validator{Power: 2}}
			BeforeEach(func() {
				params.NumBlocksPerEpoch = 2
				err := valKeeper.db.Put(storage.NewFromKeyValue(MakeBlockValidatorsKey(2), util.ObjectToBytes(valset)))
				Expect(err).To(BeNil())
				err = valKeeper.db.Put(storage.NewFromKeyValue(MakeBlockValidatorsKey(4), util.ObjectToBytes(valset2)))
				Expect(err).To(BeNil())
			})

			It("should return valset2 since it is the most recent set", func() {
				res, err := valKeeper.GetByHeight(9)
				Expect(err).To(BeNil())
				Expect(res).To(Equal(types.BlockValidators(valset2)))
			})
		})

		When("two validators exist; valset1 at height 2, valset2 at height 4; argument height = 10; blocks per epoch = 2", func() {
			valset := map[string]*types.Validator{"pubkey": &types.Validator{Power: 1}}
			valset2 := map[string]*types.Validator{"pubkey": &types.Validator{Power: 2}}
			BeforeEach(func() {
				params.NumBlocksPerEpoch = 2
				err := valKeeper.db.Put(storage.NewFromKeyValue(MakeBlockValidatorsKey(2), util.ObjectToBytes(valset)))
				Expect(err).To(BeNil())
				err = valKeeper.db.Put(storage.NewFromKeyValue(MakeBlockValidatorsKey(4), util.ObjectToBytes(valset2)))
				Expect(err).To(BeNil())
			})

			It("should return valset2 since it is the most recent set", func() {
				res, err := valKeeper.GetByHeight(10)
				Expect(err).To(BeNil())
				Expect(res).To(Equal(types.BlockValidators(valset2)))
			})
		})
	})

	Describe(".Index", func() {
		var err error
		When("no issues with database", func() {
			BeforeEach(func() {
				vals := []*types.Validator{&types.Validator{PubKey: pubKey, Power: 1}}
				err = valKeeper.Index(1, vals)
			})

			It("should successfully index validators", func() {
				Expect(err).To(BeNil())
			})

			Specify("that key exist in db", func() {
				rec, err := valKeeper.db.Get(MakeBlockValidatorsKey(1))
				Expect(err).To(BeNil())
				Expect(rec).ToNot(BeNil())
			})
		})

		When("db.Put returns an error", func() {
			BeforeEach(func() {
				mockDB := storagemocks.NewMockFunctions(ctrl)
				mockDB.EXPECT().Put(gomock.Any()).Return(fmt.Errorf("error"))
				valKeeper.db = mockDB
			})

			BeforeEach(func() {
				vals := []*types.Validator{&types.Validator{PubKey: pubKey, Power: 1}}
				err = valKeeper.Index(1, vals)
			})

			It("should successfully index validators", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("failed to index validators: error"))
			})

		})
	})
})
