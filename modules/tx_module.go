package modules

import (
	"encoding/hex"
	"fmt"
	types3 "gitlab.com/makeos/mosdef/logic/types"
	types2 "gitlab.com/makeos/mosdef/services/types"
	"gitlab.com/makeos/mosdef/types/msgs"
	"strings"

	"github.com/c-bata/go-prompt"
	"gitlab.com/makeos/mosdef/types"
	"gitlab.com/makeos/mosdef/util"
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
)

// TxModule provides transaction functionalities to JS environment
type TxModule struct {
	vm      *otto.Otto
	keepers types3.Keepers
	service types2.Service
}

// NewTxModule creates an instance of TxModule
func NewTxModule(vm *otto.Otto, service types2.Service, keepers types3.Keepers) *TxModule {
	return &TxModule{vm: vm, service: service, keepers: keepers}
}

// txCoinFuncs are functions accessible using the `tx.coin` namespace
func (m *TxModule) txCoinFuncs() []*types.ModulesAggregatorFunc {
	return []*types.ModulesAggregatorFunc{
		&types.ModulesAggregatorFunc{
			Name:        "send",
			Value:       m.sendTx,
			Description: "Send coins to another account",
		},
	}
}

// funcs are functions accessible using the `tx` namespace
func (m *TxModule) funcs() []*types.ModulesAggregatorFunc {
	return []*types.ModulesAggregatorFunc{
		&types.ModulesAggregatorFunc{
			Name:        "get",
			Value:       m.Get,
			Description: "Get a transactions by its hash",
		},
		&types.ModulesAggregatorFunc{
			Name:        "sendPayload",
			Value:       m.SendPayload,
			Description: "Send a signed transaction payload to the network",
		},
	}
}

func (m *TxModule) globals() []*types.ModulesAggregatorFunc {
	return []*types.ModulesAggregatorFunc{}
}

// Configure configures the JS context and return
// any number of console prompt suggestions
func (m *TxModule) Configure() []prompt.Suggest {
	suggestions := []prompt.Suggest{}

	// Add the main tx namespace
	txMap := map[string]interface{}{}
	util.VMSet(m.vm, types.NamespaceTx, txMap)

	// Add 'coin' namespaced functions
	coinMap := map[string]interface{}{}
	txMap[types.NamespaceCoin] = coinMap
	for _, f := range m.txCoinFuncs() {
		coinMap[f.Name] = f.Value
		funcFullName := fmt.Sprintf("%s.%s.%s", types.NamespaceTx, types.NamespaceCoin, f.Name)
		suggestions = append(suggestions, prompt.Suggest{Text: funcFullName,
			Description: f.Description})
	}

	// Add other funcs to `tx` namespace
	for _, f := range m.funcs() {
		txMap[f.Name] = f.Value
		funcFullName := fmt.Sprintf("%s.%s", types.NamespaceTx, f.Name)
		suggestions = append(suggestions, prompt.Suggest{Text: funcFullName,
			Description: f.Description})
	}

	// Add global functions
	for _, f := range m.globals() {
		m.vm.Set(f.Name, f.Value)
		suggestions = append(suggestions, prompt.Suggest{Text: f.Name,
			Description: f.Description})
	}

	return suggestions
}

// sendTx sends the native coin from a source account
// to a destination account. It returns an object containing
// the hash of the transaction. It panics when an error occurs.
//
// params {
// 		nonce: number,
//		fee: string,
// 		value: string,
//		to: string
//		timestamp: number
// }
// options[0]: key
// options[1]: payloadOnly - When true, returns the payload only, without sending the tx.
func (m *TxModule) sendTx(params map[string]interface{}, options ...interface{}) interface{} {
	var err error

	var tx = msgs.NewBareTxCoinTransfer()
	if err = tx.FromMap(params); err != nil {
		panic(err)
	}

	payloadOnly := finalizeTx(tx, m.service, options...)
	if payloadOnly {
		return EncodeForJS(tx.ToMap())
	}

	hash, err := m.service.SendTx(tx)
	if err != nil {
		panic(errors.Wrap(err, "failed to send transaction"))
	}

	return EncodeForJS(map[string]interface{}{
		"hash": hash,
	})
}

// get fetches a tx by its hash
func (m *TxModule) Get(hash string) interface{} {

	if strings.ToLower(hash[:2]) == "0x" {
		hash = hash[2:]
	}

	// decode the hash from hex to byte
	bz, err := hex.DecodeString(hash)
	if err != nil {
		panic(errors.Wrap(err, "invalid transaction hash"))
	}

	tx, err := m.keepers.TxKeeper().GetTx(bz)
	if err != nil {
		panic(err)
	}

	return EncodeForJS(tx)
}

// sendPayload sends an already signed transaction object to the network
func (m *TxModule) SendPayload(txData map[string]interface{}) interface{} {
	tx, err := msgs.DecodeTxFromMap(txData)
	if err != nil {
		panic(err)
	}

	hash, err := m.service.SendTx(tx)
	if err != nil {
		panic(errors.Wrap(err, "failed to send transaction"))
	}

	return EncodeForJS(map[string]interface{}{
		"hash": hash,
	})
}
