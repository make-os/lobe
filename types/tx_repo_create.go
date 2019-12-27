package types

import (
	"github.com/fatih/structs"
	"github.com/makeos/mosdef/util"
	"github.com/vmihailenco/msgpack"
)

// TxRepoCreate implements BaseTx, it describes a transaction that creates a
// repository for the signer
type TxRepoCreate struct {
	*TxCommon `json:"-" mapstructure:"-"`
	*TxType   `json:"-" msgpack:"-"`
	*TxValue  `json:"-" mapstructure:"-"`
	Name      string `json:"name" msgpack:"name"`
}

// NewBareTxRepoCreate returns an instance of TxRepoCreate with zero values
func NewBareTxRepoCreate() *TxRepoCreate {
	return &TxRepoCreate{
		TxCommon: NewBareTxCommon(),
		TxType:   &TxType{Type: TxTypeRepoCreate},
		TxValue:  &TxValue{Value: "0"},
		Name:     "",
	}
}

// EncodeMsgpack implements msgpack.CustomEncoder
func (tx *TxRepoCreate) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeMulti(
		tx.Type,
		tx.Nonce,
		tx.Fee,
		tx.Sig,
		tx.Timestamp,
		tx.SenderPubKey,
		tx.Value,
		tx.Name)
}

// DecodeMsgpack implements msgpack.CustomDecoder
func (tx *TxRepoCreate) DecodeMsgpack(dec *msgpack.Decoder) error {
	return dec.DecodeMulti(
		&tx.Type,
		&tx.Nonce,
		&tx.Fee,
		&tx.Sig,
		&tx.Timestamp,
		&tx.SenderPubKey,
		&tx.Value,
		&tx.Name)
}

// Bytes returns the serialized transaction
func (tx *TxRepoCreate) Bytes() []byte {
	return util.ObjectToBytes(tx)
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
func (tx *TxRepoCreate) ComputeHash() util.Hash {
	return util.BytesToHash(util.Blake2b256(tx.Bytes()))
}

// GetHash returns the hash of the transaction
func (tx *TxRepoCreate) GetHash() util.Hash {
	return tx.ComputeHash()
}

// GetID returns the id of the transaction (also the hash)
func (tx *TxRepoCreate) GetID() string {
	return tx.ComputeHash().HexStr()
}

// GetEcoSize returns the size of the transaction for use in protocol economics
func (tx *TxRepoCreate) GetEcoSize() int64 {
	fee := tx.Fee
	tx.Fee = ""
	bz := tx.Bytes()
	tx.Fee = fee
	return int64(len(bz))
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
	s := structs.New(tx)
	s.TagName = "json"
	return s.Map()
}