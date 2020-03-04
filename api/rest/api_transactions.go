package rest

import (
	"encoding/json"
	"net/http"

	"gitlab.com/makeos/mosdef/util"
)

// TxSendPayload sends a signed transaction to the mempool
// Body (JSON): map that conforms to any valid types.BaseTx transaction
// Response <map>
// - hash <string>: The hash of the transaction
func (r *RESTApi) TxSendPayload(w http.ResponseWriter, req *http.Request) {
	var body = make(map[string]interface{})
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		util.WriteJSON(w, 400, util.RESTApiErrorMsg("malformed body", "", "0"))
		return
	}
	util.WriteJSON(w, 201, r.Modules().Tx.SendPayload(body))
}