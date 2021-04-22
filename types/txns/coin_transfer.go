package txns

import (
	"github.com/make-os/kit/crypto/ed25519"
	"github.com/make-os/kit/types"
	"github.com/make-os/kit/util"
	"github.com/make-os/kit/util/errors"
	"github.com/make-os/kit/util/identifier"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/vmihailenco/msgpack"
)

// TxCoinTransfer implements BaseTx, it describes a transaction that transfers
// the native coin from one account to another.
type TxCoinTransfer struct {
	*TxType      `json:",flatten" msgpack:"-" mapstructure:"-"`
	*TxCommon    `json:",flatten" msgpack:"-" mapstructure:"-"`
	*TxRecipient `json:",flatten" msgpack:"-" mapstructure:"-"`
	*TxValue     `json:",flatten" msgpack:"-" mapstructure:"-"`
}

// NewBareTxCoinTransfer returns an instance of TxCoinTransfer with zero values
func NewBareTxCoinTransfer() *TxCoinTransfer {
	return &TxCoinTransfer{
		TxType:      &TxType{Type: TxTypeCoinTransfer},
		TxCommon:    NewBareTxCommon(),
		TxRecipient: &TxRecipient{To: ""},
		TxValue:     &TxValue{Value: "0"},
	}
}

// NewCoinTransferTx creates and populates a coin transfer transaction
func NewCoinTransferTx(
	nonce uint64,
	to identifier.Address,
	senderKey *ed25519.Key,
	value util.String,
	fee util.String,
	timestamp int64) (baseTx types.BaseTx) {

	tx := NewBareTxCoinTransfer()
	tx.SetRecipient(to)
	tx.SetValue(value)
	baseTx = tx

	baseTx.SetTimestamp(timestamp)
	baseTx.SetFee(fee)
	baseTx.SetNonce(nonce)
	baseTx.SetSenderPubKey(senderKey.PubKey().MustBytes())
	sig, err := baseTx.Sign(senderKey.PrivKey().Base58())
	if err != nil {
		panic(err)
	}
	baseTx.SetSignature(sig)
	return
}

// EncodeMsgpack implements msgpack.CustomEncoder
func (tx *TxCoinTransfer) EncodeMsgpack(enc *msgpack.Encoder) error {
	return tx.EncodeMulti(enc,
		tx.Type,
		tx.Nonce,
		tx.Fee,
		tx.Sig,
		tx.Timestamp,
		tx.SenderPubKey,
		tx.To,
		tx.Value)
}

// DecodeMsgpack implements msgpack.CustomDecoder
func (tx *TxCoinTransfer) DecodeMsgpack(dec *msgpack.Decoder) error {
	return tx.DecodeMulti(dec,
		&tx.Type,
		&tx.Nonce,
		&tx.Fee,
		&tx.Sig,
		&tx.Timestamp,
		&tx.SenderPubKey,
		&tx.To,
		&tx.Value)
}

// Bytes returns the serialized transaction
func (tx *TxCoinTransfer) Bytes() []byte {
	return util.ToBytes(tx)
}

// GetBytesNoSig returns the serialized the transaction excluding the signature
func (tx *TxCoinTransfer) GetBytesNoSig() []byte {
	sig := tx.Sig
	tx.Sig = nil
	bz := tx.Bytes()
	tx.Sig = sig
	return bz
}

// ComputeHash computes the hash of the transaction
func (tx *TxCoinTransfer) ComputeHash() util.Bytes32 {
	return util.BytesToBytes32(tmhash.Sum(tx.Bytes()))
}

// GetHash returns the hash of the transaction
func (tx *TxCoinTransfer) GetHash() util.HexBytes {
	return tx.ComputeHash().ToHexBytes()
}

// GetID returns the id of the transaction (also the hash)
func (tx *TxCoinTransfer) GetID() string {
	return tx.ComputeHash().HexStr()
}

// GetEcoSize returns the size of the transaction for use in protocol economics
func (tx *TxCoinTransfer) GetEcoSize() int64 {
	return tx.GetSize()
}

// GetSize returns the size of the tx object (excluding nothing)
func (tx *TxCoinTransfer) GetSize() int64 {
	return int64(len(tx.Bytes()))
}

// Sign signs the transaction
func (tx *TxCoinTransfer) Sign(privKey string) ([]byte, error) {
	return SignTransaction(tx, privKey)
}

// ToMap returns a map equivalent of the transaction
func (tx *TxCoinTransfer) ToMap() map[string]interface{} {
	return util.ToJSONMap(tx)
}

// FromMap populates tx with a map generated by tx.ToMap.
func (tx *TxCoinTransfer) FromMap(data map[string]interface{}) error {
	err := tx.TxCommon.FromMap(data)
	err = errors.CallOnNilErr(err, func() error { return tx.TxType.FromMap(data) })
	err = errors.CallOnNilErr(err, func() error { return tx.TxRecipient.FromMap(data) })
	err = errors.CallOnNilErr(err, func() error { return tx.TxValue.FromMap(data) })
	return err
}
