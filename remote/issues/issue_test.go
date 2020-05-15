package issues_test

import (
	"fmt"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/makeos/mosdef/mocks"
	"gitlab.com/makeos/mosdef/remote/issues"
	"gitlab.com/makeos/mosdef/remote/plumbing"
	"gitlab.com/makeos/mosdef/remote/repo"
	"gitlab.com/makeos/mosdef/util"

	"gitlab.com/makeos/mosdef/config"
	"gitlab.com/makeos/mosdef/testutil"
)

var _ = Describe("Issue", func() {
	var err error
	var cfg *config.AppConfig
	var ctrl *gomock.Controller

	BeforeEach(func() {
		cfg, err = testutil.SetTestCfg()
		Expect(err).To(BeNil())
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
		err = os.RemoveAll(cfg.DataDir())
		Expect(err).To(BeNil())
	})

	Describe(".CreateIssueComment", func() {
		var mockRepo *mocks.MockLocalRepo

		BeforeEach(func() {
			mockRepo = mocks.NewMockLocalRepo(ctrl)
		})

		When("issue number is not provided", func() {
			It("should return err when unable to get free issue number", func() {
				mockRepo.EXPECT().GetFreeIssueNum(1).Return(0, fmt.Errorf("error"))
				_, _, err := issues.CreateIssueComment(mockRepo, 0, "", false)
				Expect(err).ToNot(BeNil())
				Expect(err).To(MatchError("failed to find free issue number: error"))
			})
		})

		It("should return error when unable to get issue reference", func() {
			mockRepo.EXPECT().GetFreeIssueNum(1).Return(1, nil)
			mockRepo.EXPECT().RefGet(plumbing.MakeIssueReference(1)).Return("", fmt.Errorf("error"))
			_, _, err := issues.CreateIssueComment(mockRepo, 0, "", false)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to check issue existence: error"))
		})

		When("comment is request but issue does not exist", func() {
			It("should return err", func() {
				mockRepo.EXPECT().GetFreeIssueNum(1).Return(1, nil)
				mockRepo.EXPECT().RefGet(plumbing.MakeIssueReference(1)).Return("", repo.ErrRefNotFound)
				_, _, err := issues.CreateIssueComment(mockRepo, 0, "", true)
				Expect(err).ToNot(BeNil())
				Expect(err).To(MatchError("can't add comment to a non-existing issue (1)"))
			})
		})

		It("should return err when unable to create a single file commit", func() {
			hash := util.RandString(40)
			mockRepo.EXPECT().GetFreeIssueNum(1).Return(1, nil)
			mockRepo.EXPECT().RefGet(plumbing.MakeIssueReference(1)).Return(hash, nil)
			mockRepo.EXPECT().CreateSingleFileCommit("body", "body content", "", hash).Return("", fmt.Errorf("error"))
			_, _, err := issues.CreateIssueComment(mockRepo, 0, "body content", true)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to create issue commit: error"))
		})

		It("should return err when unable to update issue reference target hash", func() {
			refname := plumbing.MakeIssueReference(1)
			hash := util.RandString(40)
			issueHash := util.RandString(40)
			mockRepo.EXPECT().GetFreeIssueNum(1).Return(1, nil)
			mockRepo.EXPECT().RefGet(refname).Return(hash, nil)
			mockRepo.EXPECT().CreateSingleFileCommit("body", "body content", "", hash).Return(issueHash, nil)
			mockRepo.EXPECT().RefUpdate(refname, issueHash).Return(fmt.Errorf("error"))
			_, _, err := issues.CreateIssueComment(mockRepo, 0, "body content", true)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to update issue reference target hash: error"))
		})

		It("should return no error when successful", func() {
			refname := plumbing.MakeIssueReference(1)
			hash := util.RandString(40)
			issueHash := util.RandString(40)
			mockRepo.EXPECT().GetFreeIssueNum(1).Return(1, nil)
			mockRepo.EXPECT().RefGet(refname).Return(hash, nil)
			mockRepo.EXPECT().CreateSingleFileCommit("body", "body content", "", hash).Return(issueHash, nil)
			mockRepo.EXPECT().RefUpdate(refname, issueHash).Return(nil)
			isNewIssue, issueReference, err := issues.CreateIssueComment(mockRepo, 0, "body content", true)
			Expect(err).To(BeNil())
			Expect(isNewIssue).To(BeFalse())
			Expect(issueReference).To(Equal(refname))
		})
	})
})
