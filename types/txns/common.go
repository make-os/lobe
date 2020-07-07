package txns

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/spf13/cast"
	"github.com/stretchr/objx"
	"github.com/vmihailenco/msgpack"
	msgpack2 "github.com/vmihailenco/msgpack/v4"
	"gitlab.com/makeos/mosdef/crypto"
	"gitlab.com/makeos/mosdef/types"
	"gitlab.com/makeos/mosdef/util"
)

// All Transaction type
const (
	TxTypeCoinTransfer                types.TxCode = iota + 1 // For native coin transfer to/between accounts
	TxTypeValidatorTicket                                     // For validator ticket purchase
	TxTypeSetDelegatorCommission                              // For setting delegator commission
	TxTypeHostTicket                                          // For purchasing host ticket
	TxTypeUnbondHostTicket                                    // For unbonding host ticket
	TxTypeRepoCreate                                          // For creating a repository
	TxTypeRegisterPushKey                                     // For adding a push key
	TxTypePush                                                // For pushing updates to a repository
	TxTypeNamespaceRegister                                   // For namespace purchase
	TxTypeNamespaceDomainUpdate                               // For setting namespace domains
	TxTypeRepoProposalUpsertOwner                             // For creating a proposal to add repo owner
	TxTypeRepoProposalVote                                    // For voting on a repo proposal
	TxTypeRepoProposalUpdate                                  // For creating a repo update proposal
	TxTypeRepoProposalSendFee                                 // For native coin transfer to repo as proposal fee
	TxTypeRepoProposalRegisterPushKey                         // For adding push keys to a repo
	TxTypeUpDelPushKey                                        // For updating or deleting a push key
	MergeRequestProposalAction                                // For identifying merge request proposal
)

// TxMeta stores arbitrary, self-contained state information for a transaction
type TxMeta struct {
	meta map[string]interface{}
}

// GetMeta returns the meta information of the transaction
func (m *TxMeta) GetMeta() map[string]interface{} {
	return m.meta
}

// TxType implements some of BaseTx, it includes type information about a transaction
type TxType struct {
	Type types.TxCode `json:"type" msgpack:"type" mapstructure:"type"`
}

// GetType returns the type of the transaction
func (tx *TxType) GetType() types.TxCode {
	return tx.Type
}

// Is checks if the tx is a given type
func (tx *TxType) Is(txType types.TxCode) bool {
	return tx.Type == txType
}

// ToBasicMap returns a map equivalent of the transaction, decode-able by FromMap.
func (tx *TxType) ToMap() map[string]interface{} {
	m := util.ToMap(tx, "mapstructure")
	m["type"] = int(tx.Type)
	return m
}

// FromMap populates tx with a map generated by ToBasicMap.
func (tx *TxType) FromMap(data map[string]interface{}) (err error) {
	o := objx.New(data)

	// Type: expects int, int64 or float64 types in map
	if typeVal := o.Get("type"); !typeVal.IsNil() {
		if typeVal.IsInt64() {
			tx.Type = types.TxCode(typeVal.Int64())
		} else if typeVal.IsInt() {
			tx.Type = types.TxCode(typeVal.Int())
		} else if typeVal.IsFloat64() {
			tx.Type = types.TxCode(typeVal.Float64())
		} else {
			return util.FieldError("type", fmt.Sprintf("invalid value type: has %T, "+
				"wants int", typeVal.Inter()))
		}
	}

	return nil
}

// TxCommon implements some of BaseTx, it includes some common fields and methods
type TxCommon struct {
	util.CodecUtil `json:"-" msgpack:"-" mapstructure:"-"`
	*TxMeta        `json:"-" msgpack:"-" mapstructure:"-"`
	Nonce          uint64           `json:"nonce" msgpack:"nonce" mapstructure:"nonce"`
	Fee            util.String      `json:"fee" msgpack:"fee" mapstructure:"fee"`
	Sig            []byte           `json:"sig" msgpack:"sig" mapstructure:"sig"`
	Timestamp      int64            `json:"timestamp" msgpack:"timestamp" mapstructure:"timestamp"`
	SenderPubKey   crypto.PublicKey `json:"senderPubKey" msgpack:"senderPubKey" mapstructure:"senderPubKey"`
}

func (tx *TxCommon) EncodeMsgpack(enc *msgpack.Encoder) error {
	return tx.EncodeMulti(enc,
		tx.Nonce,
		tx.Fee,
		tx.Sig,
		tx.Timestamp,
		tx.SenderPubKey.Bytes())
}

func (tx *TxCommon) DecodeMsgpack(dec *msgpack2.Decoder) error {
	return tx.DecodeMulti(dec,
		&tx.Nonce,
		&tx.Fee,
		&tx.Sig,
		&tx.Timestamp,
		&tx.SenderPubKey)
}

// NewBareTxCommon returns an instance of TxCommon with zero values
func NewBareTxCommon() *TxCommon {
	return &TxCommon{
		TxMeta:       &TxMeta{meta: make(map[string]interface{})},
		Nonce:        0,
		Fee:          "0",
		Timestamp:    0,
		SenderPubKey: crypto.EmptyPublicKey,
	}
}

// ToBasicMap returns a map equivalent of the transaction, decode-able by FromMap.
func (tx *TxCommon) ToMap() map[string]interface{} {
	m := util.ToMap(tx, "mapstructure")
	m["fee"] = tx.Fee.String()
	m["sig"] = util.ToHex(tx.Sig)
	m["nonce"] = fmt.Sprintf("%d", tx.Nonce)

	if !tx.SenderPubKey.IsEmpty() {
		m["senderPubKey"] = crypto.MustPubKeyFromBytes(tx.SenderPubKey.Bytes()).Base58()
	}

	return m
}

// FromMap populates tx with a map generated by ToBasicMap.
func (tx *TxCommon) FromMap(data map[string]interface{}) (err error) {
	o := objx.New(data)

	// Timestamp: expects int, float64, int64 or string types in map
	if tsVal := o.Get("timestamp"); !tsVal.IsNil() {
		if tsVal.IsInt() {
			tx.Timestamp = int64(tsVal.Int())
		} else if tsVal.IsInt64() {
			tx.Timestamp = tsVal.Int64()
		} else if tsVal.IsFloat64() {
			tx.Timestamp = int64(tsVal.Float64())
		} else if tsVal.IsStr() {
			tx.Timestamp, err = strconv.ParseInt(tsVal.Str(), 10, 64)
			if err != nil {
				return util.FieldError("timestamp", "must be numeric")
			}
		} else {
			return util.FieldError("timestamp", fmt.Sprintf("invalid value type: has %T, "+
				"wants string|int", tsVal.Inter()))
		}
	}

	// Fee: expects int64, float64 or string types in map
	if feeVal := o.Get("fee"); !feeVal.IsNil() {
		if feeVal.IsInt() || feeVal.IsInt64() || feeVal.IsFloat64() {
			tx.Fee = util.String(fmt.Sprintf("%v", feeVal.Inter()))
		} else if feeVal.IsStr() {
			tx.Fee = util.String(feeVal.Str())
		} else {
			return util.FieldError("fee", fmt.Sprintf("invalid value type: has %T, "+
				"wants string|int|float", feeVal.Inter()))
		}
	}

	// Nonce: expects int64 or string types in map
	if nonceVal := o.Get("nonce"); !nonceVal.IsNil() {
		if nonceVal.IsInt() {
			tx.Nonce = uint64(nonceVal.Int())
		} else if nonceVal.IsInt64() {
			tx.Nonce = uint64(nonceVal.Int64())
		} else if nonceVal.IsFloat64() {
			tx.Nonce = uint64(nonceVal.Float64())
		} else if nonceVal.IsStr() {
			tx.Nonce, err = strconv.ParseUint(nonceVal.Str(), 10, 64)
			if err != nil {
				return util.FieldError("nonce", "must be numeric")
			}
		} else {
			return util.FieldError("nonce", fmt.Sprintf("invalid value type: has %T, "+
				"wants string|int", nonceVal.Inter()))
		}
	}

	// Sig: expects string type, hex encoded
	if sigVal := o.Get("sig"); !sigVal.IsNil() {
		if sigVal.IsStr() {
			tx.Sig, err = util.FromHex(sigVal.Str())
			if err != nil {
				return util.FieldError("sig", "unable to decode hex value")
			}
		} else {
			msg := fmt.Sprintf("invalid value type: has %T, wants hex string", sigVal.Inter())
			return util.FieldError("sig", msg)
		}
	}

	// SenderPubKey: expects string type, base58 encoded or a slice of float64 values
	if spkVal := o.Get("senderPubKey"); !spkVal.IsNil() {
		if spkVal.IsStr() {
			pubKey, err := crypto.PubKeyFromBase58(spkVal.Str())
			if err != nil {
				return util.FieldError("senderPubKey", "unable to decode from base58")
			}
			tx.SenderPubKey = crypto.BytesToPublicKey(pubKey.MustBytes())
		} else {
			msg := fmt.Sprintf("invalid value type: has %T, wants base58 string", spkVal.Inter())
			return util.FieldError("senderPubKey", msg)
		}
	}

	return nil
}

// GetFee returns the transaction nonce
func (tx *TxCommon) GetFee() util.String {
	return tx.Fee
}

// SetFee returns the transaction nonce
func (tx *TxCommon) SetFee(fee util.String) {
	tx.Fee = fee
}

// GetNonce returns the transaction nonce
func (tx *TxCommon) GetNonce() uint64 {
	return tx.Nonce
}

// SetNonce set the transaction nonce
func (tx *TxCommon) SetNonce(n uint64) {
	tx.Nonce = n
}

// GetSignature returns the transaction signature
func (tx *TxCommon) GetSignature() []byte {
	return tx.Sig
}

// SetSignature sets the transaction signature
func (tx *TxCommon) SetSignature(s []byte) {
	tx.Sig = s
}

// GetTimestamp return the transaction creation unix timestamp
func (tx *TxCommon) GetTimestamp() int64 {
	return tx.Timestamp
}

// SetTimestamp set the transaction creation unix timestamp
func (tx *TxCommon) SetTimestamp(t int64) {
	tx.Timestamp = t
}

// GetSenderPubKey returns the transaction sender public key
func (tx *TxCommon) GetSenderPubKey() crypto.PublicKey {
	return tx.SenderPubKey
}

// SetSenderPubKey set the transaction sender public key
func (tx *TxCommon) SetSenderPubKey(pk []byte) {
	tx.SenderPubKey = crypto.BytesToPublicKey(pk)
}

// GetFrom returns the address of the transaction sender
// Panics if sender's public key is invalid
func (tx *TxCommon) GetFrom() util.Address {
	pk, err := crypto.PubKeyFromBytes(tx.SenderPubKey.Bytes())
	if err != nil {
		panic(err)
	}
	return pk.Addr()
}

// SignTransaction signs a transaction.
// Expects private key in base58Check encoding.
func SignTransaction(tx types.BaseTx, privKey string) ([]byte, error) {
	pKey, err := crypto.PrivKeyFromBase58(privKey)
	if err != nil {
		return nil, err
	}

	sig, err := pKey.Sign(tx.GetBytesNoSig())
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// TxRecipient describes a transaction receiver
type TxRecipient struct {
	To util.Address `json:"to" msgpack:"to" mapstructure:"to"`
}

// SetRecipient sets the recipient
func (tx *TxRecipient) SetRecipient(to util.Address) {
	tx.To = to
}

// ToBasicMap returns a map equivalent of the transaction, decode-able by FromMap.
func (tx *TxRecipient) ToMap() map[string]interface{} {
	return util.ToMap(tx)
}

// FromMap populates tx with a map generated by ToBasicMap.
func (tx *TxRecipient) FromMap(data map[string]interface{}) (err error) {
	o := objx.New(data)

	// To: expects string type in map
	if toVal := o.Get("to"); !toVal.IsNil() {
		switch v := toVal.Inter().(type) {
		case string:
			tx.To = util.Address(toVal.Str())
		case util.Address:
			tx.To = v
		default:
			return util.FieldError("to", fmt.Sprintf("invalid value type: has %T, wants string",
				toVal.Inter()))
		}
	}

	return nil
}

// TxValue describes a transaction value
type TxValue struct {
	Value util.String `json:"value" msgpack:"value" mapstructure:"value"`
}

// SetValue sets the value
func (tx *TxValue) SetValue(value util.String) {
	tx.Value = value
}

// ToBasicMap returns a map equivalent of the transaction, decode-able by FromMap.
func (tx *TxValue) ToMap() map[string]interface{} {
	m := util.ToMap(tx)
	m["value"] = tx.Value.String()
	return m
}

// FromMap populates tx with a map generated by ToBasicMap.
func (tx *TxValue) FromMap(data map[string]interface{}) (err error) {
	o := objx.New(data)

	// Value: expects int64, float64 or string types in map
	if valVal := o.Get("value"); !valVal.IsNil() {
		if valVal.IsInt64() || valVal.IsFloat64() {
			tx.Value = util.String(fmt.Sprintf("%v", valVal.Inter()))
		} else if valVal.IsStr() {
			tx.Value = util.String(valVal.Str())
		} else {
			return util.FieldError("value", fmt.Sprintf("invalid value type: has %T, "+
				"wants string|int64|float", valVal.Inter()))
		}
	}

	return nil
}

// TxProposalCommon describes proposal fields
type TxProposalCommon struct {
	RepoName string      `json:"name" msgpack:"name" mapstructure:"name"`
	Value    util.String `json:"value" msgpack:"value" mapstructure:"value"`
	ID       string      `json:"id,omitempty" msgpack:"id" mapstructure:"id"`
}

// ToBasicMap returns a map equivalent of the transaction, decode-able by FromMap.
func (tx *TxProposalCommon) ToMap() map[string]interface{} {
	return util.ToMap(tx)
}

// FromMap populates tx with a map generated by ToBasicMap.
func (tx *TxProposalCommon) FromMap(data map[string]interface{}) (err error) {
	o := objx.New(data)

	// RepoName: expects string type in map
	if repoNameVal := o.Get("name"); !repoNameVal.IsNil() {
		if repoNameVal.IsStr() {
			tx.RepoName = repoNameVal.Str()
		} else {
			return util.FieldError("name", fmt.Sprintf("invalid value type: has %T, "+
				"wants string", repoNameVal.Inter()))
		}
	}

	// ProposalID: expects string type in map
	if propIDVal := o.Get("id"); !propIDVal.IsNil() {
		switch propIDVal.Inter().(type) {
		case string:
			tx.ID = propIDVal.Str()
		case int, int64, float64, uint64:
			tx.ID = cast.ToString(propIDVal.Inter())
		default:
			return util.FieldError("id", fmt.Sprintf("invalid value type: has %T, "+
				"wants string|int", propIDVal.Inter()))
		}
	}

	// Value: expects int64, float64 or string types in map
	if valVal := o.Get("value"); !valVal.IsNil() {
		if valVal.IsInt64() || valVal.IsFloat64() {
			tx.Value = util.String(fmt.Sprintf("%v", valVal.Inter()))
		} else if valVal.IsStr() {
			tx.Value = util.String(valVal.Str())
		} else {
			return util.FieldError("value", fmt.Sprintf("invalid value type: has %T, "+
				"wants string|int64|float", valVal.Inter()))
		}
	}

	return nil
}

// DecodeTxFromMap decodes a user-provided map to a transaction object.
func DecodeTxFromMap(data map[string]interface{}) (types.BaseTx, error) {
	txType := &TxType{}
	if err := txType.FromMap(data); err != nil {
		return nil, err
	}

	txObj, err := getBareTxObject(txType.Type)
	if err != nil {
		return nil, err
	}

	return txObj, txObj.FromMap(data)
}

// DecodeTx decodes msgpack data to transactions.
func DecodeTx(txBz []byte) (types.BaseTx, error) {
	dec := msgpack.NewDecoder(bytes.NewBuffer(txBz))

	// Skip object version
	dec.Skip()

	// Decode transaction type
	txType, err := dec.DecodeInt()
	if err != nil {
		return nil, fmt.Errorf("failed to decode tx type")
	}

	// Get the appropriate object for the transaction type
	tx, err := getBareTxObject(types.TxCode(txType))
	if err != nil {
		return nil, err
	}

	// Decode and return any error
	return tx, util.ToObject(txBz, tx)
}

func getBareTxObject(txType types.TxCode) (types.BaseTx, error) {
	var tx interface{}
	switch txType {
	case TxTypeCoinTransfer:
		tx = NewBareTxCoinTransfer()
	case TxTypeValidatorTicket:
		tx = NewBareTxTicketPurchase(TxTypeValidatorTicket)
	case TxTypeHostTicket:
		tx = NewBareTxTicketPurchase(TxTypeHostTicket)
	case TxTypeSetDelegatorCommission:
		tx = NewBareTxSetDelegateCommission()
	case TxTypeUnbondHostTicket:
		tx = NewBareTxTicketUnbond(TxTypeUnbondHostTicket)
	case TxTypeRepoCreate:
		tx = NewBareTxRepoCreate()
	case TxTypeRegisterPushKey:
		tx = NewBareTxRegister()
	case TxTypePush:
		tx = NewBareTxPush()
	case TxTypeNamespaceRegister:
		tx = NewBareTxNamespaceRegister()
	case TxTypeNamespaceDomainUpdate:
		tx = NewBareTxNamespaceDomainUpdate()
	case TxTypeRepoProposalUpsertOwner:
		tx = NewBareRepoProposalUpsertOwner()
	case TxTypeRepoProposalVote:
		tx = NewBareRepoProposalVote()
	case TxTypeRepoProposalUpdate:
		tx = NewBareRepoProposalUpdate()
	case TxTypeRepoProposalSendFee:
		tx = NewBareRepoProposalFeeSend()
	case TxTypeRepoProposalRegisterPushKey:
		tx = NewBareRepoProposalRegister()
	case TxTypeUpDelPushKey:
		tx = NewBareTxUpDelPushKey()
	default:
		return nil, fmt.Errorf("unsupported tx type")
	}

	return tx.(types.BaseTx), nil
}
