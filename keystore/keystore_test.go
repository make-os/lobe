package keystore

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/makeos/mosdef/config"
)

func testPrompt(resp string) promptFunc {
	return func(prompt string, args ...interface{}) string {
		return resp
	}
}

// testPrompt2 will return response with index equal to count
// count is incremented each time the function is called.
func testPrompt2(count *int, responses []string) promptFunc {
	return func(prompt string, args ...interface{}) string {
		resp := responses[*count]
		*count++
		return resp
	}
}

var _ = Describe("AccountMgr", func() {

	path := filepath.Join("./", "test_cfg")
	accountPath := filepath.Join(path, config.KeystoreDirName)

	BeforeEach(func() {
		err := os.MkdirAll(accountPath, 0700)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err := os.RemoveAll(path)
		Expect(err).To(BeNil())
	})

	Describe(".hardenPassword", func() {
		It("should return [215, 59, 34, 12, 157, 105, 253, 31, 243, 128, 41, 222, 216, 93, 165, 77, 67, 179, 85, 192, 127, 47, 171, 121, 32, 117, 125, 119, 109, 243, 32, 95]", func() {
			bs := hardenPassword([]byte("abc"))
			Expect(bs).To(Equal([]byte{215, 59, 34, 12, 157, 105, 253, 31, 243, 128, 41, 222, 216, 93, 165, 77, 67, 179, 85, 192, 127, 47, 171, 121, 32, 117, 125, 119, 109, 243, 32, 95}))
		})
	})

	Describe(".askForPassword", func() {
		am := New(accountPath)

		It("should return err = 'Passphrases did not match' when passphrase and repeat passphrase don't match", func() {
			count := 0
			am.getPassword = testPrompt2(&count, []string{"passAbc", "passAb"})
			_, err := am.AskForPassword()
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("Passphrases did not match"))
		})

		It("should return input even when no passphrase is provided the first time", func() {
			count := 0
			am.getPassword = testPrompt2(&count, []string{"", "passAb", "passAb"})
			passphrase, err := am.AskForPassword()
			Expect(err).To(BeNil())
			Expect(passphrase).To(Equal("passAb"))
		})
	})

	Describe(".askForPasswordOnce", func() {
		am := New(accountPath)

		It("should return the first input received", func() {
			count := 0
			am.getPassword = testPrompt2(&count, []string{"", "", "passAb"})
			passphrase := am.AskForPasswordOnce()
			Expect(passphrase).To(Equal("passAb"))
		})
	})
})