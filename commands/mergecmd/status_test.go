package mergecmd_test

import (
	"bytes"
	"fmt"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/themakeos/lobe/commands/mergecmd"
	"github.com/themakeos/lobe/config"
	"github.com/themakeos/lobe/mocks"
	"github.com/themakeos/lobe/remote/plumbing"
	"github.com/themakeos/lobe/remote/types"
	"github.com/themakeos/lobe/testutil"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var _ = Describe("MergeReqStatus", func() {
	var err error
	var cfg *config.AppConfig
	var ctrl *gomock.Controller
	var mockRepo *mocks.MockLocalRepo

	BeforeEach(func() {
		cfg, err = testutil.SetTestCfg()
		Expect(err).To(BeNil())
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockLocalRepo(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
		err = os.RemoveAll(cfg.DataDir())
		Expect(err).To(BeNil())
	})

	Describe(".MergeReqStatusCmd", func() {
		It("should return error when unable to get reference", func() {
			ref := plumbing.MakeMergeRequestReference(1)
			mockRepo.EXPECT().RefGet(ref).Return("", fmt.Errorf("error"))
			err := mergecmd.MergeReqStatusCmd(mockRepo, &mergecmd.MergeReqStatusArgs{Reference: ref})
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("error"))
		})

		It("should return error when merge request reference does not exist", func() {
			ref := plumbing.MakeMergeRequestReference(1)
			mockRepo.EXPECT().RefGet(ref).Return("", plumbing.ErrRefNotFound)
			err := mergecmd.MergeReqStatusCmd(mockRepo, &mergecmd.MergeReqStatusArgs{Reference: ref})
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("merge request not found"))
		})

		It("should return error when unable to read recent commit post body", func() {
			ref := plumbing.MakeMergeRequestReference(1)
			hash := "e31992a88829f3cb70ab5f5e964597a6c8f17047"
			mockRepo.EXPECT().RefGet(ref).Return(hash, nil)
			err := mergecmd.MergeReqStatusCmd(mockRepo, &mergecmd.MergeReqStatusArgs{
				Reference: ref,
				ReadPostBody: func(repo types.LocalRepo, hash string) (*plumbing.PostBody, *object.Commit, error) {
					return nil, nil, fmt.Errorf("error")
				},
			})
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to read recent comment: error"))
		})

		It("should print 'opened' when post body includes a closed=false directive", func() {
			ref := plumbing.MakeMergeRequestReference(1)
			hash := "e31992a88829f3cb70ab5f5e964597a6c8f17047"
			mockRepo.EXPECT().RefGet(ref).Return(hash, nil)
			buf := bytes.NewBuffer(nil)
			err := mergecmd.MergeReqStatusCmd(mockRepo, &mergecmd.MergeReqStatusArgs{
				Reference: ref,
				ReadPostBody: func(repo types.LocalRepo, hash string) (*plumbing.PostBody, *object.Commit, error) {
					closed := false
					return &plumbing.PostBody{Close: &closed}, nil, nil
				},
				StdOut: buf,
			})
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal("opened\n"))
		})

		It("should print 'closed' when post body includes a closed=true directive", func() {
			ref := plumbing.MakeMergeRequestReference(1)
			hash := "e31992a88829f3cb70ab5f5e964597a6c8f17047"
			mockRepo.EXPECT().RefGet(ref).Return(hash, nil)
			buf := bytes.NewBuffer(nil)
			err := mergecmd.MergeReqStatusCmd(mockRepo, &mergecmd.MergeReqStatusArgs{
				Reference: ref,
				ReadPostBody: func(repo types.LocalRepo, hash string) (*plumbing.PostBody, *object.Commit, error) {
					closed := true
					return &plumbing.PostBody{Close: &closed}, nil, nil
				},
				StdOut: buf,
			})
			Expect(err).To(BeNil())
			Expect(buf.String()).To(Equal("closed\n"))
		})
	})
})