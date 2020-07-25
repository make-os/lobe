package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	"github.com/imroc/req"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/themakeos/lobe/util"
)

var _ = Describe("Account", func() {
	var ctrl *gomock.Controller
	var client *ClientV1

	BeforeEach(func() {
		client = &ClientV1{apiRoot: ""}
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe(".SendTxPayload", func() {
		It("should send payload in request and receive tx hash from server", func() {
			client.post = func(endpoint string, params map[string]interface{}) (resp *req.Resp, err error) {
				Expect(endpoint).To(Equal("/v1/tx/send-payload"))

				mockReqHandler := func(w http.ResponseWriter, r *http.Request) {
					data, _ := json.Marshal(util.Map{"hash": "0x12345"})
					w.WriteHeader(201)
					w.Write(data)
				}
				ts := httptest.NewServer(http.HandlerFunc(mockReqHandler))
				resp, _ = req.Get(ts.URL)

				return resp, nil
			}
			resp, err := client.SendTxPayload(map[string]interface{}{"type": 1})
			Expect(err).To(BeNil())
			Expect(resp.Hash).To(Equal("0x12345"))
		})

		It("should return response as error when status code is not 201", func() {
			client.post = func(endpoint string, params map[string]interface{}) (resp *req.Resp, err error) {
				Expect(endpoint).To(Equal("/v1/tx/send-payload"))

				mockReqHandler := func(w http.ResponseWriter, r *http.Request) {
					data, _ := json.Marshal(util.Map{"hash": "0x12345"})
					w.WriteHeader(400)
					w.Write(data)
				}
				ts := httptest.NewServer(http.HandlerFunc(mockReqHandler))
				resp, _ = req.Get(ts.URL)

				return resp, nil
			}
			_, err := client.SendTxPayload(map[string]interface{}{"type": 1})
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(`{"hash":"0x12345"}`))
		})
	})

	Describe(".GetTransaction", func() {
		It("should return error if request failed", func() {
			client.get = func(endpoint string, params map[string]interface{}) (resp *req.Resp, err error) {
				Expect(endpoint).To(Equal("/v1/tx/get"))
				Expect(params).To(HaveKey("hash"))
				Expect(params["hash"]).To(Equal("0x123"))
				return resp, fmt.Errorf("error")
			}
			_, err := client.GetTransaction("0x123")
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("error"))
		})

		It("should return response on success", func() {
			client.get = func(endpoint string, params map[string]interface{}) (resp *req.Resp, err error) {
				mockReqHandler := func(w http.ResponseWriter, r *http.Request) {
					data, _ := json.Marshal(map[string]interface{}{"value": "10.4"})
					w.WriteHeader(200)
					w.Write(data)
				}
				ts := httptest.NewServer(http.HandlerFunc(mockReqHandler))
				resp, _ = req.Get(ts.URL)
				return resp, nil
			}
			resp, err := client.GetTransaction("repo1")
			Expect(err).To(BeNil())
			Expect(resp).To(Equal(map[string]interface{}{
				"value": "10.4",
			}))
		})
	})
})