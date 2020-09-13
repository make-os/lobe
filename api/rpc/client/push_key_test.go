package client

import (
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/make-os/lobe/api/types"
	"github.com/make-os/lobe/crypto"
	"github.com/make-os/lobe/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var client *RPCClient
	var ctrl *gomock.Controller
	var key = crypto.NewKeyFromIntSeed(1)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		client = NewClient(&Options{Host: "127.0.0.1", Port: 8000})
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe(".GetOwner", func() {
		It("should return ReqError when call failed", func() {
			client.call = func(method string, params interface{}) (res util.Map, statusCode int, err error) {
				Expect(method).To(Equal("pk_getOwner"))
				return nil, 0, fmt.Errorf("error")
			}
			_, err := client.PushKey().GetOwner("pk1_abc", 100)
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(&util.ReqError{
				Code:     ErrCodeUnexpected,
				HttpCode: 0,
				Msg:      "error",
				Field:    "",
			}))
		})

		It("should return error when RPC call response could not be decoded", func() {
			client.call = func(method string, params interface{}) (res util.Map, statusCode int, err error) {
				Expect(method).To(Equal("pk_getOwner"))
				return util.Map{"balance": 1000}, 0, nil
			}
			_, err := client.PushKey().GetOwner("pk1_abc", 100)
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(&util.ReqError{
				Code:     "decode_error",
				HttpCode: 500,
				Msg:      "field:balance, msg:invalid value type: has int, wants string",
				Field:    "",
			}))
		})

		It("should return expected result on success", func() {
			client.call = func(method string, params interface{}) (res util.Map, statusCode int, err error) {
				Expect(method).To(Equal("pk_getOwner"))
				return util.Map{"balance": "1000"}, 0, nil
			}
			acct, err := client.PushKey().GetOwner("pk1_abc", 100)
			Expect(err).To(BeNil())
			Expect(acct.Balance).To(Equal(util.String("1000")))
		})
	})

	Describe(".Register()", func() {
		It("should return ReqError when signing key is not provided", func() {
			_, err := client.PushKey().Register(&types.BodyRegisterPushKey{
				SigningKey: nil,
			})
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(&util.ReqError{
				Code:     ErrCodeBadParam,
				HttpCode: 400,
				Msg:      "signing key is required",
				Field:    "signingKey",
			}))
		})

		It("should return ReqError when call failed", func() {
			client.call = func(method string, params interface{}) (res util.Map, statusCode int, err error) {
				Expect(params).To(And(
					HaveKey("senderPubKey"),
					HaveKey("feeCap"),
					HaveKey("pubKey"),
					HaveKey("sig"),
					HaveKey("timestamp"),
					HaveKey("type"),
					HaveKey("scopes"),
					HaveKey("nonce"),
					HaveKey("fee"),
				))
				return nil, 0, fmt.Errorf("error")
			}
			_, err := client.PushKey().Register(&types.BodyRegisterPushKey{
				Nonce:      100,
				Fee:        1,
				Scopes:     []string{"scope1"},
				FeeCap:     1.2,
				PublicKey:  key.PubKey().ToPublicKey(),
				SigningKey: key,
			})
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(&util.ReqError{
				Code:     ErrCodeUnexpected,
				HttpCode: 0,
				Msg:      "error",
				Field:    "",
			}))
		})

		It("should return expected address and transaction hash on success", func() {
			client.call = func(method string, params interface{}) (res util.Map, statusCode int, err error) {
				Expect(method).To(Equal("pk_register"))
				return util.Map{"address": "pk1abc", "hash": "0x123"}, 0, nil
			}
			resp, err := client.PushKey().Register(&types.BodyRegisterPushKey{
				Nonce:      100,
				Fee:        1,
				SigningKey: key,
			})
			Expect(err).To(BeNil())
			Expect(resp.Address).To(Equal("pk1abc"))
			Expect(resp.Hash).To(Equal("0x123"))
		})
	})
})
