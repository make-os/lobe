package purchaseticket

import (
	"fmt"

	"github.com/themakeos/lobe/crypto"
	"github.com/themakeos/lobe/params"
	"github.com/themakeos/lobe/types"
	"github.com/themakeos/lobe/types/core"
	"github.com/themakeos/lobe/types/state"
	"github.com/themakeos/lobe/types/txns"
	"github.com/themakeos/lobe/util"
)

// TicketPurchaseContract is a system contract for purchasing a ticket.
// TicketPurchaseContract implements SystemContract.
type TicketPurchaseContract struct {
	core.Logic
	tx          *txns.TxTicketPurchase
	chainHeight uint64
}

// NewContract creates an instance of TicketPurchaseContract
func NewContract() *TicketPurchaseContract {
	return &TicketPurchaseContract{}
}

func (c *TicketPurchaseContract) CanExec(typ types.TxCode) bool {
	return typ == txns.TxTypeValidatorTicket || typ == txns.TxTypeHostTicket
}

func (c *TicketPurchaseContract) Init(logic core.Logic, tx types.BaseTx, curChainHeight uint64) core.SystemContract {
	c.Logic = logic
	c.tx = tx.(*txns.TxTicketPurchase)
	c.chainHeight = curChainHeight
	return c
}

func (c *TicketPurchaseContract) Exec() error {

	spk, _ := crypto.PubKeyFromBytes(c.tx.SenderPubKey.Bytes())
	acctKeeper := c.AccountKeeper()
	fee := c.tx.Fee
	txType := c.tx.GetType()
	value := c.tx.Value

	// Get sender accounts
	sender := spk.Addr()
	senderAcct := acctKeeper.Get(sender)

	// Deduct the transaction fee and increment nonce
	senderBal := senderAcct.Balance.Decimal()
	senderAcct.Balance = util.String(senderBal.Sub(fee.Decimal()).String())
	senderAcct.Nonce = senderAcct.Nonce + 1

	// Register a stake entry
	unbondHeight := uint64(0)
	switch txType {
	case txns.TxTypeValidatorTicket:

		// Determine unbond height. The unbond height is height of the next block
		// (or proposed block) plus minimum ticket maturation duration, max ticket
		// active duration + thawing period.
		unbondHeight = c.chainHeight + 1 + uint64(params.MinTicketMatDur) +
			uint64(params.MaxTicketActiveDur) +
			uint64(params.NumBlocksInThawPeriod)
		senderAcct.Stakes.Add(state.StakeTypeValidator, value, unbondHeight)

	case txns.TxTypeHostTicket:
		senderAcct.Stakes.Add(state.StakeTypeHost, value, unbondHeight)

	default:
		return fmt.Errorf("unknown transaction type")
	}

	// Update the sender's account
	senderAcct.Clean(c.chainHeight)
	acctKeeper.Update(sender, senderAcct)

	return nil
}
