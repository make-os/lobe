package depositproposalfee

import (
	"gitlab.com/makeos/mosdef/crypto"
	"gitlab.com/makeos/mosdef/logic/contracts/common"
	"gitlab.com/makeos/mosdef/types"
	"gitlab.com/makeos/mosdef/types/core"
	"gitlab.com/makeos/mosdef/types/txns"
)

// DepositProposalFeeContract is a system contract for adding deposit fee to a proposal.
// DepositProposalFeeContract implements SystemContract.
type DepositProposalFeeContract struct {
	core.Logic
	tx          *txns.TxRepoProposalSendFee
	chainHeight uint64
	contracts   []core.SystemContract
}

// NewContract creates a new instance of DepositProposalFeeContract
func NewContract() *DepositProposalFeeContract {
	return &DepositProposalFeeContract{}
}

func (c *DepositProposalFeeContract) CanExec(typ types.TxCode) bool {
	return typ == txns.TxTypeRepoProposalSendFee
}

// Init initialize the contract
func (c *DepositProposalFeeContract) Init(logic core.Logic, tx types.BaseTx, curChainHeight uint64) core.SystemContract {
	c.Logic = logic
	c.tx = tx.(*txns.TxRepoProposalSendFee)
	c.chainHeight = curChainHeight
	return c
}

// Exec executes the contract
func (c *DepositProposalFeeContract) Exec() error {
	spk, _ := crypto.PubKeyFromBytes(c.tx.SenderPubKey.Bytes())

	// Get the repo and proposal
	repoKeeper := c.RepoKeeper()
	repo := repoKeeper.Get(c.tx.RepoName)
	prop := repo.Proposals.Get(c.tx.ProposalID)

	// Register proposal fee if set.
	// If the sender already deposited, update their deposit.
	if c.tx.Value != "0" {
		addr := spk.Addr().String()
		if !prop.Fees.Has(addr) {
			prop.Fees.Add(addr, c.tx.Value.String())
		} else {
			existingFee := prop.Fees.Get(addr)
			updFee := existingFee.Decimal().Add(c.tx.Value.Decimal())
			prop.Fees.Add(addr, updFee.String())
		}
	}

	// Deduct network fee + proposal fee from sender
	totalFee := c.tx.Fee.Decimal().Add(c.tx.Value.Decimal())
	common.DebitAccount(c, spk, totalFee, c.chainHeight)

	repoKeeper.Update(c.tx.RepoName, repo)

	return nil
}