package ft

import (
	"github.com/goexl/gox/rand"
)

//go:inline
func (c *Client) auth(req *request, data []byte, _options *options, opts ...option) (err error) {
	// 随机生成加密密钥
	key := rand.New().String().Length(16).Generate()
	if pk, pe := c.PublicKey(opts...); nil != pe {
		err = pe
	} else {
		req.Key, err = c.encryptKey(pk, key)
	}
	if nil != err {
		return
	}

	if req.Data, err = c.cbcEncrypt(data, key, _options); nil == err {
		req.Signature, err = c.sign(data)
	}

	return
}
