package txns

import (
	"github.com/stretchr/objx"
	"github.com/vmihailenco/msgpack"
	"gitlab.com/makeos/mosdef/types"
	"gitlab.com/makeos/mosdef/util"
	"gitlab.com/makeos/mosdef/util/crypto"
)

// TxTicketUnbond implements BaseTx, it describes a transaction that unbonds a
// staked coin owned by the signer
type TxTicketUnbond struct {
	*TxType    `json:",flatten" msgpack:"-" mapstructure:"-"`
	*TxCommon  `json:",flatten" msgpack:"-" mapstructure:"-"`
	TicketHash util.HexBytes `json:"hash" msgpack:"hash" mapstructure:"hash"`
}

// NewBareTxTicketUnbond returns an instance of TxTicketUnbond with zero values
func NewBareTxTicketUnbond(ticketType types.TxCode) *TxTicketUnbond {
	return &TxTicketUnbond{
		TxType:     &TxType{Type: ticketType},
		TxCommon:   NewBareTxCommon(),
		TicketHash: []byte{},
	}
}

// EncodeMsgpack implements msgpack.CustomEncoder
func (tx *TxTicketUnbond) EncodeMsgpack(enc *msgpack.Encoder) error {
	return tx.EncodeMulti(enc,
		tx.Type,
		tx.Nonce,
		tx.Fee,
		tx.Sig,
		tx.Timestamp,
		tx.SenderPubKey,
		tx.TicketHash)
}

// DecodeMsgpack implements msgpack.CustomDecoder
func (tx *TxTicketUnbond) DecodeMsgpack(dec *msgpack.Decoder) error {
	return tx.DecodeMulti(dec,
		&tx.Type,
		&tx.Nonce,
		&tx.Fee,
		&tx.Sig,
		&tx.Timestamp,
		&tx.SenderPubKey,
		&tx.TicketHash)
}

// Bytes returns the serialized transaction
func (tx *TxTicketUnbond) Bytes() []byte {
	return util.ToBytes(tx)
}

// GetBytesNoSig returns the serialized the transaction excluding the signature
func (tx *TxTicketUnbond) GetBytesNoSig() []byte {
	sig := tx.Sig
	tx.Sig = nil
	bz := tx.Bytes()
	tx.Sig = sig
	return bz
}

// ComputeHash computes the hash of the transaction
func (tx *TxTicketUnbond) ComputeHash() util.Bytes32 {
	return util.BytesToBytes32(crypto.Blake2b256(tx.Bytes()))
}

// GetHash returns the hash of the transaction
func (tx *TxTicketUnbond) GetHash() util.HexBytes {
	return tx.ComputeHash().ToHexBytes()
}

// GetID returns the id of the transaction (also the hash)
func (tx *TxTicketUnbond) GetID() string {
	return tx.ComputeHash().HexStr()
}

// GetEcoSize returns the size of the transaction for use in protocol economics
func (tx *TxTicketUnbond) GetEcoSize() int64 {
	return tx.GetSize()
}

// GetSize returns the size of the tx object (excluding nothing)
func (tx *TxTicketUnbond) GetSize() int64 {
	return int64(len(tx.Bytes()))
}

// Sign signs the transaction
func (tx *TxTicketUnbond) Sign(privKey string) ([]byte, error) {
	return SignTransaction(tx, privKey)
}

// ToBasicMap returns a map equivalent of the transaction
func (tx *TxTicketUnbond) ToMap() map[string]interface{} {
	return util.ToBasicMap(tx)
}

// FromMap populates tx with a map generated by tx.ToMap.
func (tx *TxTicketUnbond) FromMap(data map[string]interface{}) error {
	err := tx.TxCommon.FromMap(data)
	err = util.CallOnNilErr(err, func() error { return tx.TxType.FromMap(data) })

	o := objx.New(data)
	if tickHashVal := o.Get("hash"); !tickHashVal.IsNil() && tickHashVal.IsStr() {
		bz, err := util.FromHex(tickHashVal.Str())
		if err != nil {
			return util.FieldError("hash", "unable to decode hex value")
		}
		o.Set("hash", bz)
	}

	err = util.CallOnNilErr(err, func() error { return util.DecodeMapStrict(data, &tx) })
	return err
}
