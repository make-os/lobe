package txns

import (
	"github.com/make-os/kit/util"
	"github.com/make-os/kit/util/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/vmihailenco/msgpack"
)

// TxRepoCreate implements BaseTx, it describes a transaction that creates a
// repository for the signer
type TxRepoCreate struct {
	*TxCommon      `json:",flatten" msgpack:"-" mapstructure:"-"`
	*TxType        `json:",flatten" msgpack:"-" mapstructure:"-"`
	*TxValue       `json:",flatten" msgpack:"-" mapstructure:"-"`
	*TxDescription `json:",flatten" msgpack:"-" mapstructure:"-"`
	Name           string                 `json:"name" msgpack:"name" mapstructure:"name"`
	Config         map[string]interface{} `json:"config" msgpack:"config" mapstructure:"config"`
}

// NewBareTxRepoCreate returns an instance of TxRepoCreate with zero values
func NewBareTxRepoCreate() *TxRepoCreate {
	return &TxRepoCreate{
		TxCommon:      NewBareTxCommon(),
		TxType:        &TxType{Type: TxTypeRepoCreate},
		TxValue:       &TxValue{Value: "0"},
		TxDescription: &TxDescription{Description: ""},
		Name:          "",
		Config:        make(map[string]interface{}),
	}
}

// EncodeMsgpack implements msgpack.CustomEncoder
func (tx *TxRepoCreate) EncodeMsgpack(enc *msgpack.Encoder) error {
	return tx.EncodeMulti(enc,
		tx.Type,
		tx.Nonce,
		tx.Fee,
		tx.Sig,
		tx.Timestamp,
		tx.SenderPubKey,
		tx.Value,
		tx.Name,
		tx.Config,
		tx.Description)
}

// DecodeMsgpack implements msgpack.CustomDecoder
func (tx *TxRepoCreate) DecodeMsgpack(dec *msgpack.Decoder) error {
	return tx.DecodeMulti(dec,
		&tx.Type,
		&tx.Nonce,
		&tx.Fee,
		&tx.Sig,
		&tx.Timestamp,
		&tx.SenderPubKey,
		&tx.Value,
		&tx.Name,
		&tx.Config,
		&tx.Description)
}

// Bytes returns the serialized transaction
func (tx *TxRepoCreate) Bytes() []byte {
	return util.ToBytes(tx)
}

// GetBytesNoSig returns the serialized the transaction excluding the signature
func (tx *TxRepoCreate) GetBytesNoSig() []byte {
	sig := tx.Sig
	tx.Sig = nil
	bz := tx.Bytes()
	tx.Sig = sig
	return bz
}

// ComputeHash computes the hash of the transaction
func (tx *TxRepoCreate) ComputeHash() util.Bytes32 {
	return util.BytesToBytes32(tmhash.Sum(tx.Bytes()))
}

// GetHash returns the hash of the transaction
func (tx *TxRepoCreate) GetHash() util.HexBytes {
	return tx.ComputeHash().ToHexBytes()
}

// GetID returns the id of the transaction (also the hash)
func (tx *TxRepoCreate) GetID() string {
	return tx.ComputeHash().HexStr()
}

// GetEcoSize returns the size of the transaction for use in protocol economics
func (tx *TxRepoCreate) GetEcoSize() int64 {
	return tx.GetSize()
}

// GetSize returns the size of the tx object (excluding nothing)
func (tx *TxRepoCreate) GetSize() int64 {
	return int64(len(tx.Bytes()))
}

// Sign signs the transaction
func (tx *TxRepoCreate) Sign(privKey string) ([]byte, error) {
	return SignTransaction(tx, privKey)
}

// ToMap returns a map equivalent of the transaction
func (tx *TxRepoCreate) ToMap() map[string]interface{} {
	return util.ToJSONMap(tx)
}

// FromMap populates tx with a map generated by tx.ToMap.
func (tx *TxRepoCreate) FromMap(data map[string]interface{}) error {
	err := tx.TxCommon.FromMap(data)
	err = errors.CallIfNil(err, func() error { return tx.TxType.FromMap(data) })
	err = errors.CallIfNil(err, func() error { return tx.TxDescription.FromMap(data) })
	err = errors.CallIfNil(err, func() error { return tx.TxValue.FromMap(data) })
	err = errors.CallIfNil(err, func() error { return util.DecodeMap(data, &tx) })
	return err
}
