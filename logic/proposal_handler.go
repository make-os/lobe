package logic

import (
	"fmt"
	"math"

	"github.com/makeos/mosdef/types"
)

// getProposalOutcome returns the current outcome of a proposal
// whose voters are only network stakeholders; If the proposal requires
// a proposee max join height, only stakeholders whose tickets became mature
// before the proposee max join height
func getProposalOutcome(
	tickmgr types.TicketManager,
	prop types.Proposal,
	repo *types.Repository) types.ProposalOutcome {

	var err error
	totalPower := float64(64)

	// When proposee are the owners of the repo, the power is the total
	// number of owners of the repository - one vote to one owner.
	if prop.GetProposeeType() == types.ProposeeOwner {

		totalPower = float64(len(repo.Owners))

		// If there is a max proposee max join height, eligible owners are
		// those who joined on or before the proposee max join height block height.
		if prop.GetProposeeMaxJoinHeight() > 0 {
			totalPower = 0
			repo.Owners.ForEach(func(o *types.RepoOwner, addr string) {
				if o.JoinedAt <= prop.GetProposeeMaxJoinHeight() {
					totalPower++
				}
			})
		}
	}

	// When proposees are network stakeholders, the total power is the total
	// value of tickets eligible for participation
	if prop.GetProposeeType() == types.ProposeeNetStakeholders {
		totalPower, err = tickmgr.ValueOfAllTickets(prop.GetProposeeMaxJoinHeight())
		if err != nil {
			panic(err)
		}
	}

	// Calculate total votes received, quorum and threshold.
	totalVotesReceived := prop.GetAccepted() + prop.GetRejected() + prop.GetRejectedWithVeto()
	quorum := math.Round(totalPower * (prop.GetQuorum() / 100))
	threshold := math.Round(totalVotesReceived * (prop.GetThreshold() / 100))
	vetoQuorum := math.Round(totalVotesReceived * (prop.GetVetoQuorum() / 100))

	// Check if quorum is satisfied
	if totalVotesReceived < quorum {
		return types.ProposalOutcomeQuorumNotMet
	}

	// Check if the "NoWithVeto" votes reached the veto quorum
	// Only checked if veto quorum is set above zero
	if vetoQuorum > 0 && vetoQuorum <= prop.GetRejectedWithVeto() {
		return types.ProposalOutcomeRejectedWithVeto
	}

	// When veto quorum is unset and there is at least 1 NoWithVeto vote
	if vetoQuorum == 0 && prop.GetRejectedWithVeto() > 0 {
		return types.ProposalOutcomeRejectedWithVeto
	}

	accepted := prop.GetAccepted() >= threshold
	rejected := prop.GetRejected() >= threshold

	// Check if the "Yes" votes reached the threshold
	if accepted && !rejected {
		return types.ProposalOutcomeAccepted
	}

	// Check if the "No" votes reached the threshold
	if rejected && !accepted {
		return types.ProposalOutcomeRejected
	}

	return types.ProposalOutcomeTie
}

// determineProposalOutcome determines the outcome of the proposal votes
func determineProposalOutcome(
	keepers types.Keepers,
	proposal types.Proposal,
	repo *types.Repository,
	chainHeight uint64) types.ProposalOutcome {
	return getProposalOutcome(keepers.GetTicketManager(), proposal, repo)
}

// maybeApplyProposal attempts to apply the action of a proposal
func maybeApplyProposal(
	keepers types.Keepers,
	proposal types.Proposal,
	repo *types.Repository,
	chainHeight uint64) (bool, error) {

	if proposal.IsFinalized() {
		return false, nil
	}

	var outcome types.ProposalOutcome
	isOwnersOnlyProposal := proposal.GetProposeeType() == types.ProposeeOwner

	// When allowed voters are only the repo owners and there is just one owner
	// whom is also the creator of the proposal, instantly apply the proposal.
	if isOwnersOnlyProposal && len(repo.Owners) == 1 &&
		repo.Owners.Has(proposal.GetCreator()) {
		proposal.SetOutcome(types.ProposalOutcomeAccepted)
		goto apply
	}

	// Don't apply the proposal if the proposal end height is in the future.
	if proposal.GetEndAt() > chainHeight+1 {
		return false, nil
	}

	// Here, the proposal has come to its end. We need to determine if the
	// outcome was an acceptance, if not we return false.
	outcome = determineProposalOutcome(keepers, proposal, repo, chainHeight)
	proposal.SetOutcome(outcome)
	if outcome != types.ProposalOutcomeAccepted {
		return false, nil
	}

apply:
	switch proposal.GetAction() {
	case types.ProposalActionAddOwner:
		return true, applyProposalAddOwner(proposal, repo, chainHeight)
	case types.ProposalActionRepoUpdate:
		return true, applyProposalRepoUpdate(proposal, repo, chainHeight)
	}

	return false, fmt.Errorf("unsupported proposal action")
}

// maybeApplyEndedProposals finds and applies proposals that will
// end at the given height.
func maybeApplyEndedProposals(
	keepers types.Keepers,
	nextChainHeight uint64) error {

	repoKeeper := keepers.RepoKeeper()

	// Find proposals ending at the given height
	endingProps := repoKeeper.GetProposalsEndingAt(nextChainHeight)

	// Attempt to apply and close the proposal
	for _, ep := range endingProps {
		repo := repoKeeper.GetRepo(ep.RepoName)
		if repo.IsNil() {
			return fmt.Errorf("repo not found") // should never happen
		}
		_, err := maybeApplyProposal(keepers, repo.Proposals.Get(ep.ProposalID),
			repo, nextChainHeight-1)
		if err != nil {
			return err
		}
		repoKeeper.Update(ep.RepoName, repo)
	}

	return nil
}
