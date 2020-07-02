package txns

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/stretchr/objx"
	"github.com/vmihailenco/msgpack"
	"gitlab.com/makeos/mosdef/util"
	"gitlab.com/makeos/mosdef/util/crypto"
)

// TxSetDelegateCommission implements BaseTx, it describes a transaction that
// sets the signers delegate commission rate.
type TxSetDelegateCommission struct {
	*TxType    `json:",flatten" msgpack:"-" mapstructure:"-"`
	*TxCommon  `json:",flatten" msgpack:"-" mapstructure:"-"`
	Commission util.String `json:"commission" msgpack:"commission" mapstructure:"commission"`
}

// NewBareTxSetDelegateCommission returns an instance of TxSetDelegateCommission with zero values
func NewBareTxSetDelegateCommission() *TxSetDelegateCommission {
	return &TxSetDelegateCommission{
		TxType:     &TxType{Type: TxTypeSetDelegatorCommission},
		TxCommon:   NewBareTxCommon(),
		Commission: "",
	}
}

// EncodeMsgpack implements msgpack.CustomEncoder
func (tx *TxSetDelegateCommission) EncodeMsgpack(enc *msgpack.Encoder) error {
	return tx.EncodeMulti(enc,
		tx.Type,
		tx.Nonce,
		tx.Fee,
		tx.Sig,
		tx.Timestamp,
		tx.SenderPubKey,
		tx.Commission)
}

// DecodeMsgpack implements msgpack.CustomDecoder
func (tx *TxSetDelegateCommission) DecodeMsgpack(dec *msgpack.Decoder) error {
	return tx.DecodeMulti(dec,
		&tx.Type,
		&tx.Nonce,
		&tx.Fee,
		&tx.Sig,
		&tx.Timestamp,
		&tx.SenderPubKey,
		&tx.Commission)
}

// Bytes returns the serialized transaction
func (tx *TxSetDelegateCommission) Bytes() []byte {
	return util.ToBytes(tx)
}

// GetBytesNoSig returns the serialized the transaction excluding the signature
func (tx *TxSetDelegateCommission) GetBytesNoSig() []byte {
	sig := tx.Sig
	tx.Sig = nil
	bz := tx.Bytes()
	tx.Sig = sig
	return bz
}

// ComputeHash computes the hash of the transaction
func (tx *TxSetDelegateCommission) ComputeHash() util.Bytes32 {
	return util.BytesToBytes32(crypto.Blake2b256(tx.Bytes()))
}

// GetHash returns the hash of the transaction
func (tx *TxSetDelegateCommission) GetHash() util.HexBytes {
	return tx.ComputeHash().ToHexBytes()
}

// GetID returns the id of the transaction (also the hash)
func (tx *TxSetDelegateCommission) GetID() string {
	return tx.ComputeHash().HexStr()
}

// GetEcoSize returns the size of the transaction for use in protocol economics
func (tx *TxSetDelegateCommission) GetEcoSize() int64 {
	return tx.GetSize()
}

// GetSize returns the size of the tx object (excluding nothing)
func (tx *TxSetDelegateCommission) GetSize() int64 {
	return int64(len(tx.Bytes()))
}

// Sign signs the transaction
func (tx *TxSetDelegateCommission) Sign(privKey string) ([]byte, error) {
	return SignTransaction(tx, privKey)
}

// ToMap returns a map equivalent of the transaction
func (tx *TxSetDelegateCommission) ToMap() map[string]interface{} {
	m := util.StructToMap(tx, "mapstructure")
	mergo.Map(&m, tx.TxType.ToMap())
	mergo.Map(&m, tx.TxCommon.ToMap())
	return m
}

// FromMap populates tx with a map generated by ToMap.
func (tx *TxSetDelegateCommission) FromMap(data map[string]interface{}) error {
	err := tx.TxCommon.FromMap(data)
	err = util.CallOnNilErr(err, func() error { return tx.TxType.FromMap(data) })

	o := objx.New(data)

	// Commission: expects int64, float64 or string types in map
	if commissionVal := o.Get("commission"); !commissionVal.IsNil() {
		if commissionVal.IsInt64() || commissionVal.IsFloat64() {
			tx.Commission = util.String(fmt.Sprintf("%v", commissionVal.Inter()))
		} else if commissionVal.IsStr() {
			tx.Commission = util.String(commissionVal.Str())
		} else {
			return util.FieldError("commission", fmt.Sprintf("invalid value type: has %T, "+
				"wants string|int64|float", commissionVal.Inter()))
		}
	}

	return err
}
