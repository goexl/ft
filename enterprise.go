package ft

func (c *Client) Enterprise(req *EnterpriseUploadReq, opts ...option) error {
	return c.request(uploadEnterpriseApi, req, new(empty), opts...)
}
