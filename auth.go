package ft

func (c *Client) PublicKey(opts ...option) (key string, err error) {
	_options := apply(opts...)
	if cached, ok := c.keys[_options.id]; ok {
		key = cached
	}
	if `` != key {
		return
	}

	req := new(publicKeyReq)
	req.AppId = _options.id
	rsp := new(publicKeyRsp)
	if err = c.request(getPublicKeyApi, req, rsp, opts...); nil == err {
		c.keys[_options.id] = rsp.Key
		key = rsp.Key
	}

	return
}

func (c *Client) Token(opts ...option) (token string, err error) {
	_options := apply(opts...)
	if cached, ok := c.tokens[_options.id]; ok && !cached.Expired() {
		token = cached.Token
	}
	if `` != token {
		return
	}

	_req := new(tokenReq)
	_req.AppId = _options.id
	_req.AppKey = _options.key
	_req.AppSecret = _options.secret
	_rsp := new(tokenRsp)
	if err = c.request(getTokenApi, _req, _rsp, opts...); nil == err {
		c.tokens[_options.id] = _rsp
		token = _rsp.Token
	}

	return
}
