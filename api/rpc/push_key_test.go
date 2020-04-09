package rpc

import (
	"github.com/golang/mock/gomock"
	"gitlab.com/makeos/mosdef/mocks"
	"gitlab.com/makeos/mosdef/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/makeos/mosdef/rpc"
	"gitlab.com/makeos/mosdef/types/modules"
)

var _ = Describe("PushKey", func() {
	var ctrl *gomock.Controller
	var pushApi *PushKeyAPI
	var mods *modules.Modules

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mods = &modules.Modules{}
		pushApi = &PushKeyAPI{mods}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe(".find", func() {
		testCases := map[string]testCase{
			"when id is not provided": {
				params: map[string]interface{}{},
				err:    &rpc.Err{Code: "60000", Message: "id is required", Data: "id"},
			},
			"when id type is not string": {
				params: map[string]interface{}{"id": 222},
				err:    &rpc.Err{Code: "60000", Message: "wrong value type, want 'string', got string", Data: "id"},
			},
			"when blockHeight is provided but type is not string": {
				params: map[string]interface{}{"id": "push1_abc", "blockHeight": 1},
				err:    &rpc.Err{Code: "60000", Message: "wrong value type, want 'string', got string", Data: "blockHeight"},
			},
			"when push key is successfully returned": {
				params: map[string]interface{}{"id": "push1_abc"},
				result: util.Map{
					"pubKey":  "---BEGIN PUBLIC KEY...",
					"address": "addr1",
				},
				mocker: func(tp testCase) {
					mockPushKeyMod := mocks.NewMockPushKeyModule(ctrl)
					mockPushKeyMod.EXPECT().Get("push1_abc", uint64(0)).Return(util.Map{
						"pubKey":  "---BEGIN PUBLIC KEY...",
						"address": "addr1",
					})
					mods.PushKey = mockPushKeyMod
				},
			},
		}

		for _tc, _tp := range testCases {
			tc, tp := _tc, _tp
			It(tc, func() {
				if tp.mocker != nil {
					tp.mocker(tp)
				}
				resp := pushApi.find(tp.params)
				Expect(resp).To(Equal(&rpc.Response{
					JSONRPCVersion: "2.0", Err: tp.err, Result: tp.result,
				}))
			})
		}
	})

	Describe(".getAccountOfOwner", func() {
		testCases := map[string]testCase{
			"when id is not provided": {
				params: map[string]interface{}{},
				err:    &rpc.Err{Code: "60000", Message: "id is required", Data: "id"},
			},
			"when id type is not string": {
				params: map[string]interface{}{"id": 222},
				err:    &rpc.Err{Code: "60000", Message: "wrong value type, want 'string', got string", Data: "id"},
			},
			"when blockHeight is provided but type is not string": {
				params: map[string]interface{}{"id": "push1_abc", "blockHeight": 1},
				err:    &rpc.Err{Code: "60000", Message: "wrong value type, want 'string', got string", Data: "blockHeight"},
			},
			"when account is successfully returned": {
				params: map[string]interface{}{"id": "push1_abc"},
				result: util.Map{"balance": "100", "nonce": 10, "delegatorCommission": 23},
				mocker: func(tp testCase) {
					mockPushKeyMod := mocks.NewMockPushKeyModule(ctrl)
					mockPushKeyMod.EXPECT().GetAccountOfOwner("push1_abc", uint64(0)).Return(util.Map{
						"balance":             "100",
						"nonce":               10,
						"delegatorCommission": 23,
					})
					mods.PushKey = mockPushKeyMod
				},
			},
		}

		for _tc, _tp := range testCases {
			tc, tp := _tc, _tp
			It(tc, func() {
				if tp.mocker != nil {
					tp.mocker(tp)
				}
				resp := pushApi.getAccountOfOwner(tp.params)
				Expect(resp).To(Equal(&rpc.Response{
					JSONRPCVersion: "2.0", Err: tp.err, Result: tp.result,
				}))
			})
		}
	})
})