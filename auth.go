package ft

import (
	"encoding/json"

	"github.com/goexl/gox"
)

func (c *Client) PublicKey(opts ...option) (key string, err error) {
	_options := apply(opts...)
	if cached, ok := c.keys[_options.id]; ok {
		key = cached
	}
	if `` != key {
		return
	}

	_req := new(publicKeyReq)
	_req.AppId = _options.id
	_req.PublicKey = c.hex

	if bytes, je := json.Marshal(_req); nil != je {
		err = je
	} else {
		_req.Data = string(bytes)
	}
	if nil != err {
		return
	}

	_rsp := new(publicKeyRsp)
	hr := c.options.http.R()
	hr.SetBody(_req)
	if err = c.post(`/api/publicKey`, hr, _rsp, _options); nil == err {
		c.keys[_options.id] = _rsp.Key
		key = _rsp.Key
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
	_req.PublicKey = c.hex

	// 随机生成加密密钥
	key := gox.RandString(16)
	if pk, pe := c.PublicKey(opts...); nil != pe {
		err = pe
	} else {
		_req.Key, err = c.encryptKey(pk, key)
	}
	if nil != err {
		return
	}

	if bytes, me := json.Marshal(_req); nil != me {
		err = me
	} else if _req.Data, err = c.cbcEncrypt(bytes, key, _options); nil == err {
		_req.Signature, err = c.sign(bytes)
	}
	if nil != err {
		return
	}

	_rsp := new(tokenRsp)
	hr := c.options.http.R()
	hr.SetBody(_req)
	if err = c.post(`/api/getToken`, hr, _rsp, _options); nil == err {
		c.tokens[_options.id] = _rsp
		token = _rsp.Token
	}

	return
}
