package txns

import (
	"github.com/themakeos/lobe/util"
	"github.com/themakeos/lobe/util/crypto"
	"github.com/vmihailenco/msgpack"
)

// TxRepoProposalUpsertOwner implements BaseTx, it describes a repository proposal
// transaction for adding a new owner to a repository
type TxRepoProposalUpsertOwner struct {
	*TxCommon         `json:",flatten" msgpack:"-" mapstructure:"-"`
	*TxType           `json:",flatten" msgpack:"-" mapstructure:"-"`
	*TxProposalCommon `json:",flatten" msgpack:"-" mapstructure:"-"`
	Addresses         []string `json:"addresses" msgpack:"addresses" mapstructure:"addresses"`
	Veto              bool     `json:"veto" msgpack:"veto" mapstructure:"veto"`
}

// NewBareRepoProposalUpsertOwner returns an instance of TxRepoProposalUpsertOwner with zero values
func NewBareRepoProposalUpsertOwner() *TxRepoProposalUpsertOwner {
	return &TxRepoProposalUpsertOwner{
		TxCommon:         NewBareTxCommon(),
		TxType:           &TxType{Type: TxTypeRepoProposalUpsertOwner},
		TxProposalCommon: &TxProposalCommon{Value: "0", RepoName: "", ID: ""},
		Addresses:        []string{},
		Veto:             false,
	}
}

// EncodeMsgpack implements msgpack.CustomEncoder
func (tx *TxRepoProposalUpsertOwner) EncodeMsgpack(enc *msgpack.Encoder) error {
	return tx.EncodeMulti(enc,
		tx.Type,
		tx.Nonce,
		tx.Value,
		tx.Fee,
		tx.Sig,
		tx.Timestamp,
		tx.SenderPubKey,
		tx.RepoName,
		tx.ID,
		tx.Addresses,
		tx.Veto)
}

// DecodeMsgpack implements msgpack.CustomDecoder
func (tx *TxRepoProposalUpsertOwner) DecodeMsgpack(dec *msgpack.Decoder) error {
	return tx.DecodeMulti(dec,
		&tx.Type,
		&tx.Nonce,
		&tx.Value,
		&tx.Fee,
		&tx.Sig,
		&tx.Timestamp,
		&tx.SenderPubKey,
		&tx.RepoName,
		&tx.ID,
		&tx.Addresses,
		&tx.Veto)
}

// Bytes returns the serialized transaction
func (tx *TxRepoProposalUpsertOwner) Bytes() []byte {
	return util.ToBytes(tx)
}

// GetBytesNoSig returns the serialized the transaction excluding the signature
func (tx *TxRepoProposalUpsertOwner) GetBytesNoSig() []byte {
	sig := tx.Sig
	tx.Sig = nil
	bz := tx.Bytes()
	tx.Sig = sig
	return bz
}

// ComputeHash computes the hash of the transaction
func (tx *TxRepoProposalUpsertOwner) ComputeHash() util.Bytes32 {
	return util.BytesToBytes32(crypto.Blake2b256(tx.Bytes()))
}

// GetHash returns the hash of the transaction
func (tx *TxRepoProposalUpsertOwner) GetHash() util.HexBytes {
	return tx.ComputeHash().ToHexBytes()
}

// GetID returns the id of the transaction (also the hash)
func (tx *TxRepoProposalUpsertOwner) GetID() string {
	return tx.ComputeHash().HexStr()
}

// GetEcoSize returns the size of the transaction for use in protocol economics
func (tx *TxRepoProposalUpsertOwner) GetEcoSize() int64 {
	return tx.GetSize()
}

// GetSize returns the size of the tx object (excluding nothing)
func (tx *TxRepoProposalUpsertOwner) GetSize() int64 {
	return int64(len(tx.Bytes()))
}

// Sign signs the transaction
func (tx *TxRepoProposalUpsertOwner) Sign(privKey string) ([]byte, error) {
	return SignTransaction(tx, privKey)
}

// ToBasicMap returns a map equivalent of the transaction
func (tx *TxRepoProposalUpsertOwner) ToMap() map[string]interface{} {
	return util.ToBasicMap(tx)
}

// FromMap populates tx with a map generated by tx.ToMap.
func (tx *TxRepoProposalUpsertOwner) FromMap(data map[string]interface{}) error {
	err := tx.TxCommon.FromMap(data)
	err = util.CallOnNilErr(err, func() error { return tx.TxType.FromMap(data) })
	err = util.CallOnNilErr(err, func() error { return tx.TxProposalCommon.FromMap(data) })
	err = util.CallOnNilErr(err, func() error { return util.DecodeMap(data, &tx) })
	return err
}
