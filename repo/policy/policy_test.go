package policy

import (
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/makeos/mosdef/config"
	"gitlab.com/makeos/mosdef/crypto"
	"gitlab.com/makeos/mosdef/mocks"
	"gitlab.com/makeos/mosdef/testutil"
	"gitlab.com/makeos/mosdef/types"
	"gitlab.com/makeos/mosdef/types/core"
	"gitlab.com/makeos/mosdef/types/state"
)

func testCheckTxDetail(err error) func(params *types.TxDetail, keepers core.Keepers, index int) error {
	return func(params *types.TxDetail, keepers core.Keepers, index int) error { return err }
}

var _ = Describe("Auth", func() {
	var err error
	var cfg *config.AppConfig
	var ctrl *gomock.Controller
	var mockLogic *mocks.MockLogic
	var key, key2 *crypto.Key

	BeforeEach(func() {
		cfg, err = testutil.SetTestCfg()
		Expect(err).To(BeNil())

		key = crypto.NewKeyFromIntSeed(1)
		key2 = crypto.NewKeyFromIntSeed(2)

		ctrl = gomock.NewController(GinkgoT())
		mocksObjs := testutil.MockLogic(ctrl)
		mockLogic = mocksObjs.Logic
	})

	AfterEach(func() {
		ctrl.Finish()
		err = os.RemoveAll(cfg.DataDir())
		Expect(err).To(BeNil())
	})

	Describe(".MakePusherPolicyGroups", func() {
		var polGroups [][]*state.Policy
		var repoPolicy *state.Policy
		var namespacePolicy *state.ContributorPolicy
		var contribPolicy *state.ContributorPolicy
		var targetPusherAddr string

		BeforeEach(func() {
			targetPusherAddr = key.PushAddr().String()
		})

		When("repo config, repo namespace and repo contributor entry has policies", func() {
			BeforeEach(func() {

				// Add target pusher repo config policies
				repoState := state.BareRepository()
				repoPolicy = &state.Policy{Subject: targetPusherAddr, Object: "refs/heads/master", Action: "update"}
				repoState.Config.Policies = append(repoState.Config.Policies, repoPolicy)

				// Add target pusher namespace policies
				namespacePolicy = &state.ContributorPolicy{Object: "refs/heads/about", Action: "update"}
				ns := &state.Namespace{Contributors: map[string]*state.BaseContributor{
					key.PushAddr().String(): {Policies: []*state.ContributorPolicy{namespacePolicy}},
				}}

				// Add target pusher address repo contributor policies
				contribPolicy = &state.ContributorPolicy{Object: "refs/heads/dev", Action: "delete"}
				repoState.Contributors[key.PushAddr().String()] = &state.RepoContributor{
					Policies: []*state.ContributorPolicy{contribPolicy},
				}

				polGroups = MakePusherPolicyGroups(key.PushAddr().String(), repoState, ns)
			})

			Specify("that each policy group is not empty", func() {
				Expect(polGroups).To(HaveLen(3))
				Expect(polGroups[0]).To(HaveLen(1))
				Expect(polGroups[1]).To(HaveLen(1))
				Expect(polGroups[2]).To(HaveLen(1))
			})

			Specify("that index 0 includes pusher's repo contributor policy", func() {
				Expect(polGroups[0]).To(ContainElement(&state.Policy{
					Object:  "refs/heads/dev",
					Action:  "delete",
					Subject: key.PushAddr().String(),
				}))
			})

			Specify("that index 1 includes the pusher's namespace contributor policy", func() {
				Expect(polGroups[1]).To(ContainElement(&state.Policy{
					Object:  "refs/heads/about",
					Action:  "update",
					Subject: key.PushAddr().String(),
				}))
			})

			Specify("that index 1 includes the pusher's repo config policy", func() {
				Expect(polGroups[2]).To(ContainElement(&state.Policy{
					Object:  "refs/heads/master",
					Action:  "update",
					Subject: key.PushAddr().String(),
				}))
			})
		})

		When("repo config policies include a policy whose subject is not a push key ID or 'all'", func() {
			BeforeEach(func() {
				repoState := state.BareRepository()
				repoPolicy = &state.Policy{Subject: "some_subject", Object: "refs/heads/master", Action: "update"}
				repoState.Config.Policies = append(repoState.Config.Policies, repoPolicy)
				polGroups = MakePusherPolicyGroups(key.PushAddr().String(), repoState, state.BareNamespace())
			})

			It("should not include the policy", func() {
				Expect(polGroups).To(HaveLen(3))
				Expect(polGroups[2]).To(HaveLen(0))
			})
		})

		When("repo config policies include a policy whose subject is 'all'", func() {
			BeforeEach(func() {
				repoState := state.BareRepository()
				repoPolicy = &state.Policy{Subject: "all", Object: "refs/heads/master", Action: "update"}
				repoState.Config.Policies = append(repoState.Config.Policies, repoPolicy)
				polGroups = MakePusherPolicyGroups(key.PushAddr().String(), repoState, state.BareNamespace())
			})

			It("should include the policy", func() {
				Expect(polGroups).To(HaveLen(3))
				Expect(polGroups[2]).To(HaveLen(1))
			})
		})

		When("repo config policies include a policy whose object is not a recognized reference name", func() {
			BeforeEach(func() {
				repoState := state.BareRepository()
				repoPolicy = &state.Policy{Subject: "all", Object: "master", Action: "update"}
				repoState.Config.Policies = append(repoState.Config.Policies, repoPolicy)
				polGroups = MakePusherPolicyGroups(key.PushAddr().String(), repoState, state.BareNamespace())
			})

			It("should not include the policy", func() {
				Expect(polGroups).To(HaveLen(3))
				Expect(polGroups[2]).To(HaveLen(0))
			})
		})
	})

	Describe(".CheckPolicy", func() {
		It("should return error when reference type is unknown", func() {
			enforcer := GetPolicyEnforcer([][]*state.Policy{{{Object: "obj", Subject: "sub", Action: "ac"}}})
			err := CheckPolicy(enforcer, key.PushAddr().String(), "refs/unknown/xyz", "update")
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("unknown reference (refs/unknown/xyz)"))
		})

		Context("with 'update' action", func() {
			var allowAction = "update"
			var denyAction = "deny-" + allowAction
			var enforcer EnforcerFunc
			var pushAddrA string

			BeforeEach(func() {
				pushAddrA = key.PushAddr().String()
			})

			When("action is allowed on any level", func() {
				It("should return nil at level 0", func() {
					policies := [][]*state.Policy{{{Subject: pushAddrA, Object: "refs/heads/master", Action: allowAction}}}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err).To(BeNil())
				})
				It("should return nil at level 1", func() {
					policies := [][]*state.Policy{{}, {{Subject: pushAddrA, Object: "refs/heads/master", Action: allowAction}}}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err).To(BeNil())
				})
				It("should return nil at level 2", func() {
					policies := [][]*state.Policy{{}, {}, {{Subject: pushAddrA, Object: "refs/heads/master", Action: allowAction}}}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err).To(BeNil())
				})
			})

			When("action does not have a policy", func() {
				It("should return err", func() {
					policies := [][]*state.Policy{{}, {}, {}}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err.Error()).To(Equal("reference (refs/heads/master): not authorized to perform 'update' action"))
				})
			})

			When("action is allowed on level 0 and denied on level 0", func() {
				It("should return err", func() {
					policies := [][]*state.Policy{
						{
							{Subject: pushAddrA, Object: "refs/heads/master", Action: allowAction},
							{Subject: pushAddrA, Object: "refs/heads/master", Action: denyAction},
						},
					}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("reference (refs/heads/master): not authorized to perform 'update' action"))
				})
			})

			When("action is allowed on level 0 and denied on level 1", func() {
				It("should return err", func() {
					policies := [][]*state.Policy{
						{{Subject: pushAddrA, Object: "refs/heads/master", Action: allowAction}},
						{{Subject: pushAddrA, Object: "refs/heads/master", Action: denyAction}},
					}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err).To(BeNil())
				})
			})

			When("action is denied on level 0 and allowed on level 1", func() {
				It("should return err", func() {
					policies := [][]*state.Policy{
						{{Subject: pushAddrA, Object: "refs/heads/master", Action: denyAction}},
						{{Subject: pushAddrA, Object: "refs/heads/master", Action: allowAction}},
					}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err.Error()).To(Equal("reference (refs/heads/master): not authorized to perform 'update' action"))
				})
			})

			When("action is denied on level 1 and allowed on level 2", func() {
				It("should return err", func() {
					policies := [][]*state.Policy{
						{},
						{{Subject: pushAddrA, Object: "refs/heads/master", Action: denyAction}},
						{{Subject: pushAddrA, Object: "refs/heads/master", Action: allowAction}},
					}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err.Error()).To(Equal("reference (refs/heads/master): not authorized to perform 'update' action"))
				})
			})

			When("action is allowed for subject:'all' on level 2", func() {
				It("should return nil", func() {
					policies := [][]*state.Policy{
						{}, {},
						{{Subject: "all", Object: "refs/heads/master", Action: allowAction}},
					}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err).To(BeNil())
				})
			})

			When("action is denied for subject:'all' on level 2", func() {
				It("should return error", func() {
					policies := [][]*state.Policy{
						{}, {},
						{{Subject: "all", Object: "refs/heads/master", Action: denyAction}},
					}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("reference (refs/heads/master): not authorized to perform 'update' action"))
				})
			})

			When("action is denied for subject:'all' on level 2 and allowed at level 1", func() {
				It("should return nil", func() {
					policies := [][]*state.Policy{
						{},
						{{Subject: pushAddrA, Object: "refs/heads/master", Action: allowAction}},
						{{Subject: "all", Object: "refs/heads/master", Action: denyAction}},
					}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err).To(BeNil())
				})
			})

			When("action is denied for subject:'pushAddrA' on level 2 and allowed for subject:all level 2", func() {
				It("should not authorize pushAddrA by returning error", func() {
					policies := [][]*state.Policy{
						{}, {},
						{
							{Subject: "all", Object: "refs/heads/master", Action: allowAction},
							{Subject: pushAddrA, Object: "refs/heads/master", Action: denyAction},
						},
					}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("reference (refs/heads/master): not authorized to perform 'update' action"))
				})
			})

			When("action is denied for subject:'all' on level 1 and allowed for subject:'pushAddrA' level 2", func() {
				It("should not authorize pushAddrA by returning error", func() {
					policies := [][]*state.Policy{
						{},
						{{Subject: "all", Object: "refs/heads/master", Action: denyAction}},
						{{Subject: pushAddrA, Object: "refs/heads/master", Action: allowAction}},
					}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("reference (refs/heads/master): not authorized to perform 'update' action"))
				})
			})

			When("action is denied for subject:'all' on level 1 and allowed for subject:'pushAddrA' level 0", func() {
				It("should return nil", func() {
					policies := [][]*state.Policy{
						{{Subject: pushAddrA, Object: "refs/heads/master", Action: allowAction}},
						{{Subject: "all", Object: "refs/heads/master", Action: denyAction}},
						{},
					}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err).To(BeNil())
				})
			})

			When("action is denied on dir:refs/heads as subject:'all' on level 0 and allowed on refs/heads/master on level 1", func() {
				It("should not authorize pushAddrA by returning error", func() {
					policies := [][]*state.Policy{
						{{Subject: pushAddrA, Object: "refs/heads", Action: denyAction}},
						{{Subject: pushAddrA, Object: "refs/heads/master", Action: allowAction}},
					}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/heads/master", allowAction)
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("reference (refs/heads/master): not authorized to perform 'update' action"))
				})
			})

			When("action is denied on dir:refs/heads as subject:'all' on level 0 and "+
				"dir:refs/tags as subject is allowed on level 0 and "+
				"query subject is refs/tags/tag1", func() {
				It("should return nil", func() {
					policies := [][]*state.Policy{
						{
							{Subject: "all", Object: "refs/heads", Action: denyAction},
							{Subject: pushAddrA, Object: "refs/tags", Action: allowAction},
						}, {}, {},
					}
					enforcer = GetPolicyEnforcer(policies)
					err = CheckPolicy(enforcer, pushAddrA, "refs/tags/tag1", allowAction)
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
