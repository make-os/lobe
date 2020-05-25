package txns

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/stretchr/objx"
	"github.com/vmihailenco/msgpack"
	"gitlab.com/makeos/mosdef/util"
)

// TxUpDelPushKey implements BaseTx, it describes a transaction used to update
// or delete a registered push key
type TxUpDelPushKey struct {
	*TxCommon    `json:",flatten" msgpack:"-" mapstructure:"-"`
	*TxType      `json:",flatten" msgpack:"-" mapstructure:"-"`
	ID           string      `json:"id" msgpack:"id" mapstructure:"id"`
	AddScopes    []string    `json:"addScopes" msgpack:"addScopes" mapstructure:"addScopes"`
	RemoveScopes []int       `json:"removeScopes" msgpack:"removeScopes" mapstructure:"removeScopes"`
	FeeCap       util.String `json:"feeCap" msgpack:"feeCap" mapstructure:"feeCap"`
	Delete       bool        `json:"delete" msgpack:"delete" mapstructure:"delete"`
}

// NewBareTxUpDelPushKey returns an instance of TxUpDelPushKey with zero values
func NewBareTxUpDelPushKey() *TxUpDelPushKey {
	return &TxUpDelPushKey{
		TxType:   &TxType{Type: TxTypeUpDelPushKey},
		TxCommon: NewBareTxCommon(),
	}
}

// EncodeMsgpack implements msgpack.CustomEncoder
func (tx *TxUpDelPushKey) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeMulti(
		tx.Type,
		tx.Nonce,
		tx.Fee,
		tx.Sig,
		tx.Timestamp,
		tx.SenderPubKey,
		tx.ID,
		tx.AddScopes,
		tx.RemoveScopes,
		tx.FeeCap,
		tx.Delete)
}

// DecodeMsgpack implements msgpack.CustomDecoder
func (tx *TxUpDelPushKey) DecodeMsgpack(dec *msgpack.Decoder) error {
	return tx.DecodeMulti(dec,
		&tx.Type,
		&tx.Nonce,
		&tx.Fee,
		&tx.Sig,
		&tx.Timestamp,
		&tx.SenderPubKey,
		&tx.ID,
		&tx.AddScopes,
		&tx.RemoveScopes,
		&tx.FeeCap,
		&tx.Delete)
}

// Bytes returns the serialized transaction
func (tx *TxUpDelPushKey) Bytes() []byte {
	return util.ToBytes(tx)
}

// GetBytesNoSig returns the serialized the transaction excluding the signature
func (tx *TxUpDelPushKey) GetBytesNoSig() []byte {
	sig := tx.Sig
	tx.Sig = nil
	bz := tx.Bytes()
	tx.Sig = sig
	return bz
}

// ComputeHash computes the hash of the transaction
func (tx *TxUpDelPushKey) ComputeHash() util.Bytes32 {
	return util.BytesToBytes32(util.Blake2b256(tx.Bytes()))
}

// GetHash returns the hash of the transaction
func (tx *TxUpDelPushKey) GetHash() util.Bytes32 {
	return tx.ComputeHash()
}

// GetID returns the id of the transaction (also the hash)
func (tx *TxUpDelPushKey) GetID() string {
	return tx.ComputeHash().HexStr()
}

// GetEcoSize returns the size of the transaction for use in protocol economics
func (tx *TxUpDelPushKey) GetEcoSize() int64 {
	return tx.GetSize()
}

// GetSize returns the size of the tx object (excluding nothing)
func (tx *TxUpDelPushKey) GetSize() int64 {
	return int64(len(tx.Bytes()))
}

// Sign signs the transaction
func (tx *TxUpDelPushKey) Sign(privKey string) ([]byte, error) {
	return SignTransaction(tx, privKey)
}

// ToMap returns a map equivalent of the transaction
func (tx *TxUpDelPushKey) ToMap() map[string]interface{} {
	s := structs.New(tx)
	s.TagName = "json"
	return s.Map()
}

// FromMap populates fields from a map.
// Note: Default or zero values may be set for fields that aren't present in the
// map. Also, an error will be returned when unable to convert types in map to
// actual types in the object.
func (tx *TxUpDelPushKey) FromMap(data map[string]interface{}) error {
	err := tx.TxCommon.FromMap(data)
	err = util.CallOnNilErr(err, func() error { return tx.TxType.FromMap(data) })

	o := objx.New(data)

	// ID: expects string or slice of string types in map
	if id := o.Get("id"); !id.IsNil() {
		if id.IsStr() {
			tx.ID = id.Str()
		} else {
			return util.FieldError("id", fmt.Sprintf("invalid value type: has %T, "+
				"wants string|[]string", id.Inter()))
		}
	}

	// AddScopes: expects string or slice of string types in map
	if scopesVal := o.Get("addScopes"); !scopesVal.IsNil() {
		if scopesVal.IsStr() {
			tx.AddScopes = strings.Split(scopesVal.Str(), ",")
		} else if scopesVal.IsStrSlice() {
			tx.AddScopes = scopesVal.StrSlice()
		} else {
			return util.FieldError("addScopes", fmt.Sprintf("invalid value type: has %T, "+
				"wants string|[]string", scopesVal.Inter()))
		}
	}

	// RemoveScopes: expects int64 or slice of int types in map
	if scopesVal := o.Get("removeScopes"); !scopesVal.IsNil() {
		if scopesVal.IsInt64() {
			tx.RemoveScopes = []int{int(scopesVal.Int64())}
		} else if scopesVal.IsInt64Slice() {
			for _, v := range scopesVal.Int64Slice() {
				tx.RemoveScopes = append(tx.RemoveScopes, int(v))
			}
		} else {
			return util.FieldError("removeScopes", fmt.Sprintf("invalid value type: has %T, "+
				"wants int|[]int", scopesVal.Inter()))
		}
	}

	// FeeCap: expects int64, float64 or string types in map
	if feeCap := o.Get("feeCap"); !feeCap.IsNil() {
		if feeCap.IsInt64() || feeCap.IsFloat64() {
			tx.FeeCap = util.String(fmt.Sprintf("%v", feeCap.Inter()))
		} else if feeCap.IsStr() {
			tx.FeeCap = util.String(feeCap.Str())
		} else {
			return util.FieldError("feeCap", fmt.Sprintf("invalid value type: has %T, "+
				"wants string|int64|float", feeCap.Inter()))
		}
	}

	// Delete: expects bool type in map
	if del := o.Get("delete"); !del.IsNil() {
		if del.IsBool() {
			tx.Delete = del.Bool()
		} else {
			return util.FieldError("delete", fmt.Sprintf("invalid value type: has %T, "+
				"wants bool", del.Inter()))
		}
	}

	return err
}