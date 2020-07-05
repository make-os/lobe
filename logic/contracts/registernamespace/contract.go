package registernamespace

import (
	"gitlab.com/makeos/mosdef/crypto"
	"gitlab.com/makeos/mosdef/params"
	"gitlab.com/makeos/mosdef/types"
	"gitlab.com/makeos/mosdef/types/core"
	"gitlab.com/makeos/mosdef/types/txns"
	"gitlab.com/makeos/mosdef/util"
)

// RegisterNamespaceContract is a system contract to register a namespace.
// RegisterNamespaceContract implements SystemContract.
type RegisterNamespaceContract struct {
	core.Logic
	tx          *txns.TxNamespaceRegister
	chainHeight uint64
}

// NewContract creates a new instance of RegisterNamespaceContract
func NewContract() *RegisterNamespaceContract {
	return &RegisterNamespaceContract{}
}

func (c *RegisterNamespaceContract) CanExec(typ types.TxCode) bool {
	return typ == txns.TxTypeNamespaceRegister
}

// Init initialize the contract
func (c *RegisterNamespaceContract) Init(logic core.Logic, tx types.BaseTx, curChainHeight uint64) core.SystemContract {
	c.Logic = logic
	c.tx = tx.(*txns.TxNamespaceRegister)
	c.chainHeight = curChainHeight
	return c
}

// Exec executes the contract
func (c *RegisterNamespaceContract) Exec() error {
	spk := crypto.MustPubKeyFromBytes(c.tx.SenderPubKey.Bytes())

	// Get the current namespace object and re-populate it
	ns := c.NamespaceKeeper().Get(c.tx.Name)
	ns.Owner = spk.Addr().String()
	ns.ExpiresAt.Set(c.chainHeight + uint64(params.NamespaceTTL))
	ns.GraceEndAt.Set(ns.ExpiresAt.UInt64() + uint64(params.NamespaceGraceDur))
	ns.Domains = c.tx.Domains
	ns.Owner = c.tx.TransferTo

	// Get the account of the sender
	acctKeeper := c.AccountKeeper()
	senderAcct := acctKeeper.Get(spk.Addr())

	// Deduct the fee + value
	senderAcctBal := senderAcct.Balance.Decimal()
	spendAmt := c.tx.Value.Decimal().Add(c.tx.Fee.Decimal())
	senderAcct.Balance = util.String(senderAcctBal.Sub(spendAmt).String())

	// Increment sender nonce, clean up and update
	senderAcct.Nonce = senderAcct.Nonce + 1
	senderAcct.Clean(c.chainHeight)
	acctKeeper.Update(spk.Addr(), senderAcct)

	// Send the value to the treasury
	treasuryAcct := acctKeeper.Get(util.Address(params.TreasuryAddress), c.chainHeight)
	treasuryBal := treasuryAcct.Balance.Decimal()
	treasuryAcct.Balance = util.String(treasuryBal.Add(c.tx.Value.Decimal()).String())
	treasuryAcct.Clean(c.chainHeight)
	acctKeeper.Update(util.Address(params.TreasuryAddress), treasuryAcct)

	// Update the namespace
	c.NamespaceKeeper().Update(c.tx.Name, ns)

	return nil
}