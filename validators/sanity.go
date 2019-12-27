package validators

import (
	"fmt"

	v "github.com/go-ozzo/ozzo-validation"
	"github.com/makeos/mosdef/params"
	"github.com/makeos/mosdef/types"
	"github.com/shopspring/decimal"
)

func checkRecipient(tx *types.TxRecipient, index int) error {
	if err := v.Validate(tx.To,
		v.Required.Error(feI(index, "to", "recipient address is required").Error()),
		v.By(validAddrRule(feI(index, "to", "recipient address is not valid"))),
	); err != nil {
		return err
	}
	return nil
}

func checkValue(tx *types.TxValue, index int) error {
	if err := v.Validate(tx.Value, v.Required.Error(feI(index, "value",
		"value is required").Error()), v.By(validValueRule("value", index)),
	); err != nil {
		return err
	}
	return nil
}

func checkType(tx *types.TxType, expected int, index int) error {
	if tx.Type != expected {
		return feI(index, "type", "type is invalid")
	}
	return nil
}

func checkCommon(tx types.BaseTx, index int) error {

	if err := v.Validate(tx.GetNonce(),
		v.Required.Error(feI(index, "nonce", "nonce is required").Error())); err != nil {
		return err
	}

	if err := v.Validate(tx.GetFee(),
		v.Required.Error(feI(index, "fee", "fee is required").Error()),
		v.By(validValueRule("fee", index)),
	); err != nil {
		return err
	}

	// Fee must be at least equal to the base fee
	txSize := decimal.NewFromFloat(float64(tx.GetEcoSize()))
	baseFee := params.FeePerByte.Mul(txSize)
	if tx.GetFee().Decimal().LessThan(baseFee) {
		return types.FieldErrorWithIndex(index, "fee",
			fmt.Sprintf("fee cannot be lower than the base price of %s", baseFee.StringFixed(4)))
	}

	if err := v.Validate(tx.GetTimestamp(),
		v.Required.Error(feI(index, "timestamp", "timestamp is required").Error()),
		v.By(validTimestampRule("timestamp", index)),
	); err != nil {
		return err
	}

	if err := v.Validate(tx.GetSenderPubKey(),
		v.Required.Error(feI(index, "senderPubKey", "sender public key is required").Error()),
		v.By(validPubKeyRule(feI(index, "senderPubKey", "sender public key is not valid"))),
	); err != nil {
		return err
	}

	if err := v.Validate(tx.GetSignature(),
		v.Required.Error(feI(index, "sig", "signature is required").Error()),
	); err != nil {
		return err
	}

	if sigErr := checkSignature(tx, index); len(sigErr) > 0 {
		return sigErr[0]
	}

	return nil
}

// CheckTxCoinTransfer performs sanity checks on TxCoinTransfer
func CheckTxCoinTransfer(tx *types.TxCoinTransfer, index int) error {

	if err := checkType(tx.TxType, types.TxTypeCoinTransfer, index); err != nil {
		return err
	}

	if err := checkRecipient(tx.TxRecipient, index); err != nil {
		return err
	}

	if err := checkValue(tx.TxValue, index); err != nil {
		return err
	}

	if err := checkCommon(tx, index); err != nil {
		return err
	}

	return nil
}

// CheckTxTicketPurchase performs sanity checks on TxTicketPurchase
func CheckTxTicketPurchase(tx *types.TxTicketPurchase, index int) error {

	if tx.Type != types.TxTypeValidatorTicket && tx.Type != types.TxTypeStorerTicket {
		return feI(index, "type", "type is invalid")
	}

	if err := checkValue(tx.TxValue, index); err != nil {
		return err
	}

	if tx.GetType() == types.TxTypeStorerTicket {
		if tx.Value.Decimal().LessThan(params.MinStorerStake) {
			return feI(index, "value", fmt.Sprintf("value is lower than minimum storer stake"))
		}
	}

	if tx.Delegate != "" {
		if err := v.Validate(tx.Delegate,
			v.By(validPubKeyRule(feI(index, "delegate", "requires a valid public key"))),
		); err != nil {
			return err
		}
	}

	if err := checkCommon(tx, index); err != nil {
		return err
	}

	return nil
}

// CheckTxUnbondTicket performs sanity checks on TxTicketUnbond
func CheckTxUnbondTicket(tx *types.TxTicketUnbond, index int) error {

	if err := checkType(tx.TxType, types.TxTypeStorerTicket, index); err != nil {
		return err
	}

	if err := v.Validate(tx.TicketHash,
		v.Required.Error(feI(index, "ticket", "ticket id is required").Error()),
	); err != nil {
		return err
	}

	if err := checkCommon(tx, index); err != nil {
		return err
	}

	return nil
}

// CheckTxRepoCreate performs sanity checks on TxRepoCreate
func CheckTxRepoCreate(tx *types.TxRepoCreate, index int) error {

	if err := checkType(tx.TxType, types.TxTypeRepoCreate, index); err != nil {
		return err
	}

	if err := checkValue(tx.TxValue, index); err != nil {
		return err
	}

	if err := v.Validate(tx.Name,
		v.Required.Error(feI(index, "name", "requires a unique name").Error()),
	); err != nil {
		return err
	}

	if err := checkCommon(tx, index); err != nil {
		return err
	}

	return nil
}

// CheckTxEpochSecret performs sanity checks on TxEpochSecret
func CheckTxEpochSecret(tx *types.TxEpochSecret, index int) error {

	if err := checkType(tx.TxType, types.TxTypeEpochSecret, index); err != nil {
		return err
	}

	if err := v.Validate(tx.Secret,
		v.Required.Error(feI(index, "secret", "secret is required").Error()),
		v.By(validSecretRule("secret", index)),
	); err != nil {
		return err
	}

	if err := v.Validate(tx.PreviousSecret,
		v.Required.Error(feI(index, "previousSecret", "previous secret is required").Error()),
		v.By(validSecretRule("previousSecret", index)),
	); err != nil {
		return err
	}

	if err := v.Validate(tx.SecretRound,
		v.Required.Error(feI(index, "secretRound", "secret round is required").Error()),
	); err != nil {
		return err
	}

	return nil
}

// CheckTxAddGPGPubKey performs sanity checks on TxAddGPGPubKey
func CheckTxAddGPGPubKey(tx *types.TxAddGPGPubKey, index int) error {

	if err := checkType(tx.TxType, types.TxTypeAddGPGPubKey, index); err != nil {
		return err
	}

	if err := v.Validate(tx.PublicKey,
		v.Required.Error(feI(index, "pubKey", "public key is required").Error()),
		v.By(validGPGPubKeyRule("pubKey", index)),
	); err != nil {
		return err
	}

	if err := checkCommon(tx, index); err != nil {
		return err
	}

	return nil
}

// CheckTxSetDelegateCommission performs sanity checks on TxSetDelegateCommission
func CheckTxSetDelegateCommission(tx *types.TxSetDelegateCommission, index int) error {

	if err := checkType(tx.TxType, types.TxTypeSetDelegatorCommission, index); err != nil {
		return err
	}

	if err := v.Validate(tx.Commission,
		v.Required.Error(feI(index, "commission", "commission rate is required").Error()),
	); err != nil {
		return err
	}

	if tx.Commission.Decimal().LessThan(params.MinDelegatorCommission) {
		minPct := params.MinDelegatorCommission.String()
		return feI(index, "commission", "rate cannot be below the minimum ("+minPct+"%%)")
	}

	if tx.Commission.Decimal().GreaterThan(decimal.NewFromFloat(100)) {
		return types.FieldErrorWithIndex(index, "commission", "commission rate cannot exceed 100%%")
	}

	if err := checkCommon(tx, index); err != nil {
		return err
	}

	return nil
}