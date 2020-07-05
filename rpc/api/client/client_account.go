package client

import (
	"gitlab.com/makeos/mosdef/types/api"
	"gitlab.com/makeos/mosdef/types/state"
	"gitlab.com/makeos/mosdef/util"
)

// GetAccount gets an account corresponding to a given address
func (c *RPCClient) GetAccount(address string, blockHeight ...uint64) (*api.GetAccountResponse, *util.StatusError) {

	var height uint64
	if len(blockHeight) > 0 {
		height = blockHeight[0]
	}

	resp, statusCode, err := c.call("user_get", util.Map{"address": address, "height": height})
	if err != nil {
		return nil, makeStatusErrorFromCallErr(statusCode, err)
	}

	r := &api.GetAccountResponse{Account: state.BareAccount()}
	if err = r.Account.FromMap(resp); err != nil {
		return nil, util.StatusErr(500, ErrCodeDecodeFailed, "", err.Error())
	}

	return r, nil
}