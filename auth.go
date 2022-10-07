package ft

import (
	"time"
)

func (c *Client) PublicKey(req *PublicKeyReq, opts ...option) (rsp *PublicKeyRsp, err error) {
	rsp = new(PublicKeyRsp)
	err = c.request(`/api/publicKey`, req, rsp, opts...)

	return
}

func (c *Client) Token(req *TokenReq, opts ...option) (rsp *TokenRsp, err error) {
	c.tokens.
	if nil != c.token && c.token.Expires.Before(time.Now()) {
		rsp = c.token
	} else {
		rsp, err = c.token(req, opts...)
	}

	return
}

func (c *Client) token(req *TokenReq, opts ...option) (rsp *TokenRsp, err error) {
	rsp = new(TokenRsp)
	if err = c.request(`/api/getToken`, req, rsp, opts...); nil == err {
		c.token = rsp
	}

	return
}
