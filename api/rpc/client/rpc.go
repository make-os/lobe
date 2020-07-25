package client

import (
	"github.com/themakeos/lobe/api/types"
	"github.com/themakeos/lobe/util"
)

// GetMethods gets all methods supported by the RPC server
func (c *RPCClient) GetMethods() (*types.GetMethodResponse, error) {
	resp, statusCode, err := c.call("rpc_methods", nil)
	if err != nil {
		return nil, makeStatusErrorFromCallErr(statusCode, err)
	}

	var r types.GetMethodResponse
	if err := util.DecodeMap(resp, &r); err != nil {
		return nil, util.ReqErr(500, ErrCodeDecodeFailed, "", err.Error())
	}

	return &r, nil
}