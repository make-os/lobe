package rest

import (
	"fmt"
	"path"
	"strings"

	"github.com/imroc/req"
)

// M is actually a map[string]interface{}
type M map[string]interface{}

func joinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}

// Client is a REST API client
type Client struct {
	apiRoot string
}

// NewClient creates an instance of Client;
//
// ARGS:
// apiRoot is the URL path to the API server
func NewClient(apiRoot string) *Client {
	return &Client{apiRoot}
}

// GetCall makes a get call to the endpoint
func (c *Client) GetCall(endpoint string, params map[string]interface{}) (*req.Resp, error) {
	url := joinURL(c.apiRoot, endpoint)
	return req.Get(url, req.QueryParam(params))
}

// PostCall makes a get call to the endpoint
func (c *Client) PostCall(endpoint string, body map[string]interface{}) (*req.Resp, error) {
	url := joinURL(c.apiRoot, endpoint)
	return req.Post(url, req.BodyJSON(body))
}
