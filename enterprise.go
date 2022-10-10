package ft

func (c *Client) Enterprise(req *EnterpriseReq, opts ...option) error {
	return c.request(getPublicKeyApi, req, new(empty), opts...)
}
