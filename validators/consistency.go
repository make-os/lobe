package validators

import (
	"crypto/rsa"
	"fmt"

	"github.com/makeos/mosdef/crypto"
	"github.com/makeos/mosdef/params"
	"github.com/makeos/mosdef/repo"
	"github.com/makeos/mosdef/types"
	"github.com/makeos/mosdef/util"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// CheckTxCoinTransferConsistency performs consistency checks on TxCoinTransfer
func CheckTxCoinTransferConsistency(
	tx *types.TxCoinTransfer,
	index int,
	logic types.Logic) error {

	bi, err := logic.SysKeeper().GetLastBlockInfo()
	if err != nil {
		return errors.Wrap(err, "failed to fetch current block info")
	}

	pubKey, _ := crypto.PubKeyFromBytes(tx.GetSenderPubKey().Bytes())
	if err = logic.Tx().CanExecCoinTransfer(tx.GetType(), pubKey, tx.Value, tx.Fee,
		tx.GetNonce(), uint64(bi.Height)); err != nil {
		return err
	}

	return nil
}

// CheckTxTicketPurchaseConsistency performs consistency checks on TxTicketPurchase
func CheckTxTicketPurchaseConsistency(
	tx *types.TxTicketPurchase,
	index int,
	logic types.Logic) error {

	bi, err := logic.SysKeeper().GetLastBlockInfo()
	if err != nil {
		return errors.Wrap(err, "failed to fetch current block info")
	}

	// When delegate is set, the delegate must have an active, non-delegated ticket
	if !tx.Delegate.IsEmpty() {
		r, err := logic.GetTicketManager().GetActiveTicketsByProposer(tx.Delegate, tx.Type, false)
		if err != nil {
			return errors.Wrap(err, "failed to get active delegate tickets")
		} else if len(r) == 0 {
			return feI(index, "delegate", "specified delegate is not active")
		}
	}

	// For validator ticket transaction, the value must not be lesser than
	// the current price per ticket
	if tx.Type == types.TxTypeValidatorTicket {
		curTicketPrice := logic.Sys().GetCurValidatorTicketPrice()
		if tx.Value.Decimal().LessThan(decimal.NewFromFloat(curTicketPrice)) {
			return feI(index, "value", fmt.Sprintf("value is lower than the"+
				" minimum ticket price (%f)", curTicketPrice))
		}
	}

	pubKey, _ := crypto.PubKeyFromBytes(tx.GetSenderPubKey().Bytes())
	if err = logic.Tx().CanExecCoinTransfer(tx.GetType(), pubKey, tx.Value, tx.Fee,
		tx.GetNonce(), uint64(bi.Height)); err != nil {
		return err
	}

	return nil
}

// CheckTxUnbondTicketConsistency performs consistency checks on TxTicketUnbond
func CheckTxUnbondTicketConsistency(
	tx *types.TxTicketUnbond,
	index int,
	logic types.Logic) error {

	bi, err := logic.SysKeeper().GetLastBlockInfo()
	if err != nil {
		return errors.Wrap(err, "failed to fetch current block info")
	}

	// The ticket must exist
	ticket := logic.GetTicketManager().GetByHash(tx.TicketHash)
	if ticket == nil {
		return feI(index, "hash", "ticket not found")
	}

	// Ensure the tx creator is the owner of the ticket.
	// For delegated ticket, compare the delegator address with the sender address
	authErr := feI(index, "hash", "sender not authorized to unbond this ticket")
	if ticket.Delegator == "" {
		if tx.SenderPubKey != ticket.ProposerPubKey {
			return authErr
		}
	} else if ticket.Delegator != tx.GetFrom().String() {
		return authErr
	}

	// Ensure the ticket is still active
	decayBy := ticket.DecayBy
	if decayBy != 0 && decayBy > uint64(bi.Height) {
		return feI(index, "hash", "ticket is already decaying")
	} else if decayBy != 0 && decayBy <= uint64(bi.Height) {
		return feI(index, "hash", "ticket has already decayed")
	}

	pubKey, _ := crypto.PubKeyFromBytes(tx.GetSenderPubKey().Bytes())
	if err = logic.Tx().CanExecCoinTransfer(tx.GetType(), pubKey, "0", tx.Fee,
		tx.GetNonce(), uint64(bi.Height)); err != nil {
		return err
	}

	return nil
}

// CheckTxRepoCreateConsistency performs consistency checks on TxRepoCreate
func CheckTxRepoCreateConsistency(
	tx *types.TxRepoCreate,
	index int,
	logic types.Logic) error {

	bi, err := logic.SysKeeper().GetLastBlockInfo()
	if err != nil {
		return errors.Wrap(err, "failed to fetch current block info")
	}

	repo := logic.RepoKeeper().GetRepo(tx.Name)
	if !repo.IsNil() {
		msg := "name is not available. choose another"
		return feI(index, "name", msg)
	}

	pubKey, _ := crypto.PubKeyFromBytes(tx.GetSenderPubKey().Bytes())
	if err = logic.Tx().CanExecCoinTransfer(tx.GetType(), pubKey, tx.Value, tx.Fee,
		tx.GetNonce(), uint64(bi.Height)); err != nil {
		return err
	}

	return nil
}

// CheckTxEpochSeedConsistency performs consistency checks on TxEpochSeed
func CheckTxEpochSeedConsistency(
	tx *types.TxEpochSeed,
	index int,
	logic types.Logic) error {

	return nil
}

// CheckTxSetDelegateCommissionConsistency performs consistency checks on TxSetDelegateCommission
func CheckTxSetDelegateCommissionConsistency(
	tx *types.TxSetDelegateCommission,
	index int,
	logic types.Logic) error {

	bi, err := logic.SysKeeper().GetLastBlockInfo()
	if err != nil {
		return errors.Wrap(err, "failed to fetch current block info")
	}

	pubKey, _ := crypto.PubKeyFromBytes(tx.GetSenderPubKey().Bytes())
	if err = logic.Tx().CanExecCoinTransfer(tx.GetType(), pubKey, "0", tx.Fee,
		tx.GetNonce(), uint64(bi.Height)); err != nil {
		return err
	}

	return nil
}

// CheckTxAddGPGPubKeyConsistency performs consistency checks on TxAddGPGPubKey
func CheckTxAddGPGPubKeyConsistency(
	tx *types.TxAddGPGPubKey,
	index int,
	logic types.Logic) error {

	bi, err := logic.SysKeeper().GetLastBlockInfo()
	if err != nil {
		return errors.Wrap(err, "failed to fetch current block info")
	}

	entity, _ := crypto.PGPEntityFromPubKey(tx.PublicKey)
	pk := entity.PrimaryKey.PublicKey.(*rsa.PublicKey)

	// Ensure bit length is not less than 256
	if pk.Size() < 256 {
		msg := "gpg public key bit length must be at least 2048 bits"
		return feI(index, "pubKey", msg)
	}

	// Check whether there is a matching gpg key already existing
	pkID := util.RSAPubKeyID(entity.PrimaryKey.PublicKey.(*rsa.PublicKey))
	gpgPubKey := logic.GPGPubKeyKeeper().GetGPGPubKey(pkID)
	if !gpgPubKey.IsNil() {
		return feI(index, "pubKey", "gpg public key already registered")
	}

	pubKey, _ := crypto.PubKeyFromBytes(tx.GetSenderPubKey().Bytes())
	if err = logic.Tx().CanExecCoinTransfer(tx.GetType(), pubKey, "0", tx.Fee,
		tx.GetNonce(), uint64(bi.Height)); err != nil {
		return err
	}

	return nil
}

// CheckTxPushConsistency performs consistency checks on TxPush
func CheckTxPushConsistency(
	tx *types.TxPush,
	index int,
	logic types.Logic,
	repoGetter func(name string) (types.BareRepo, error)) error {

	storers, err := logic.GetTicketManager().GetTopStorers(params.NumTopStorersLimit)
	if err != nil {
		return errors.Wrap(err, "failed to get top storers")
	}

	localRepo, err := repoGetter(tx.PushNote.GetRepoName())
	if err != nil {
		return errors.Wrap(err, "failed to get repo")
	}

	for _, pok := range tx.PushOKs {

		// Ensure that the signers of the PushOK are part of the storers
		spk, err := crypto.PubKeyFromBytes(pok.SenderPubKey.Bytes())
		if err != nil {
			return err
		}
		if !storers.Has(spk.MustBytes32()) {
			return feI(index, "endorsements.senderPubKey",
				"sender public key does not belong to an active storer")
		}

		// Ensure the references hash are valid
		for i, refHash := range pok.ReferencesHash {
			ref := tx.PushNote.References[i]
			curRefHash, err := localRepo.TreeRoot(ref.Name)
			if err != nil {
				return errors.Wrapf(err, "failed to get reference (%s) tree root hash", ref.Name)
			}
			if !refHash.Hash.Equal(curRefHash) {
				return feI(index, "endorsements.refsHash",
					fmt.Sprintf("wrong tree hash for reference (%s)", ref.Name))
			}
		}
	}

	// Check push note
	if err := repo.CheckPushNoteConsistency(tx.PushNote, logic); err != nil {
		return err
	}

	return nil
}

// CheckTxNSAcquireConsistency performs consistency checks on TxNamespaceAcquire
func CheckTxNSAcquireConsistency(
	tx *types.TxNamespaceAcquire,
	index int,
	logic types.Logic) error {

	bi, err := logic.SysKeeper().GetLastBlockInfo()
	if err != nil {
		return errors.Wrap(err, "failed to fetch current block info")
	}

	ns := logic.NamespaceKeeper().GetNamespace(tx.Name)
	if !ns.IsNil() && ns.GraceEndAt > uint64(bi.Height) {
		return feI(index, "name", "chosen name is not currently available")
	}

	pubKey, _ := crypto.PubKeyFromBytes(tx.GetSenderPubKey().Bytes())
	if err = logic.Tx().CanExecCoinTransfer(tx.GetType(), pubKey, tx.Value, tx.Fee,
		tx.GetNonce(), uint64(bi.Height)); err != nil {
		return err
	}

	return nil
}

// CheckTxNamespaceDomainUpdateConsistency performs consistency
// checks on TxNamespaceDomainUpdate
func CheckTxNamespaceDomainUpdateConsistency(
	tx *types.TxNamespaceDomainUpdate,
	index int,
	logic types.Logic) error {

	bi, err := logic.SysKeeper().GetLastBlockInfo()
	if err != nil {
		return errors.Wrap(err, "failed to fetch current block info")
	}

	pubKey, _ := crypto.PubKeyFromBytes(tx.GetSenderPubKey().Bytes())

	// Ensure the sender of the transaction is the owner of the namespace
	ns := logic.NamespaceKeeper().GetNamespace(tx.Name)
	if ns.IsNil() {
		return feI(index, "name", "namespace not found")
	}

	if ns.Owner != pubKey.Addr().String() {
		return feI(index, "senderPubKey", "sender not permitted to perform this operation")
	}

	if err = logic.Tx().CanExecCoinTransfer(tx.GetType(), pubKey, "0", tx.Fee,
		tx.GetNonce(), uint64(bi.Height)); err != nil {
		return err
	}

	return nil
}
