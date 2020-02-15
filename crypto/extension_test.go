package crypto

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WrappedPV", func() {
	Describe(".GetKey", func() {
		It("should return expected key and no error", func() {
			wpv := GenerateWrappedPV([]byte("abc"))
			key, err := wpv.GetKey()
			Expect(err).To(BeNil())
			Expect(key.Addr().String()).To(Equal("maker1m4aaslnzmdp4k3g52tk6eh94ghr547exleegs7"))
		})
	})
})
