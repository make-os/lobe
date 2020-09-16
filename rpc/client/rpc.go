package client

import (
	"github.com/make-os/lobe/rpc"
	"github.com/make-os/lobe/util"
)

// RPCAPI provides access to the rpc server-related methods
type RPCAPI struct {
	c *RPCClient
}

// GetMethods gets all methods supported by the RPC server
func (c *RPCAPI) GetMethods() ([]rpc.MethodInfo, error) {
	resp, statusCode, err := c.c.call("rpc_methods", nil)
	if err != nil {
		return nil, makeStatusErrorFromCallErr(statusCode, err)
	}

	var r []rpc.MethodInfo
	if err := util.DecodeMap(resp["methods"], &r); err != nil {
		return nil, util.ReqErr(500, ErrCodeDecodeFailed, "", err.Error())
	}

	return r, nil
}