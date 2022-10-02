package ft

func (c *Client) PublicKey(req *PublicKeyReq) (rsp *PublicKeyRsp, err error) {
	rsp = new(PublicKeyRsp)
	err = c.request(`/api/publicKey`, req, rsp)

	return
}
