package ft

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/emmansun/gmsm/sm2"
	"github.com/emmansun/gmsm/sm4"
	"github.com/emmansun/gmsm/smx509"
	"github.com/go-resty/resty/v2"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/simaqian"
)

var _ = New

// Client 客户端
type Client struct {
	http       *resty.Client
	key        *sm2.PrivateKey
	privateHex string
	publicHex  string
	logger     simaqian.Logger
}

// New 创建客户端
func New(opts ...newOption) (client *Client, err error) {
	_options := defaultNewOptions()
	for _, opt := range opts {
		opt.applyNew(_options)
	}

	client = new(Client)
	if nil != _options.http {
		client.http = _options.http
	} else {
		client.http = resty.New()
	}
	if nil != _options.logger {
		client.logger = _options.logger
	} else {
		client.logger = simaqian.Must()
	}

	if client.key, err = sm2.GenerateKey(rand.Reader); nil != err {
		return
	}

	if pk, me := smx509.MarshalPKIXPublicKey(client.key.Public()); nil != me {
		err = me
	} else {
		client.publicHex = strings.ToUpper(hex.EncodeToString(pk))
	}
	if nil != err {
		return
	}

	if pk, me := smx509.MarshalSM2PrivateKey(client.key); nil != me {
		err = me
	} else {
		client.privateHex = strings.ToUpper(hex.EncodeToString(pk))
	}

	return
}

//go:inline
func (c *Client) request(api string, _req any, rsp any, opts ...option) (err error) {
	fr := new(req)
	fr.PublicKey = c.publicHex

	// 加密请求
	if bytes, je := json.Marshal(_req); nil != je {
		err = je
	} else {
		so := sm2.NewPlainEncrypterOpts(sm2.MarshalUncompressed, sm2.C1C2C3)
		fr.Data, err = sm2.Encrypt(rand.Reader, &c.key.PublicKey, bytes, so)
	}
	if nil != err {
		return
	}

	hr := c.http.R()
	hr.SetBody(fr)
	_options := apply(opts...)
	err = c.post(api, hr, rsp, _options)

	return
}

//go:inline
func (c *Client) sendfile(api string, file string, req any, rsp any, opts ...option) (err error) {
	_options := apply(opts...)
	hr := c.http.R()
	if form, formErr := gox.StructToForm(req); nil != formErr {
		err = formErr
	} else {
		form[`publicKey`] = c.publicHex
		form[`token`] = `eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhcHBJZCI6IjE1Njk2MTMyMzE4ODQ4ODYwMTciLCJleHBpcmVzVGltZSI6MTY2NDYzMDY4Mzc5OX0.nDZDWLutKRpbNwq3ltCo4uFQ3V456fxnyeWry8cbre8`
		hr.SetFormData(form)
	}
	if nil != err {
		return
	}

	// 设置上传文件路径
	if `` != file {
		hr.SetFile(`file`, file)
	}
	err = c.post(api, hr, rsp, _options)

	return
}

//go:inline
func (c *Client) post(api string, req *resty.Request, rsp any, _options *options) (err error) {
	if raw, reqErr := req.Post(fmt.Sprintf(`%s%s`, _options.host, api)); nil != reqErr {
		err = reqErr
		c.logger.Error(`发送到省大数据中心出错`, field.String(`api`, api), field.Error(err))
	} else if raw.IsError() {
		c.logger.Warn(`发送到省大数据中心出错`, field.String(`api`, api), field.String(`raw`, raw.String()))
	} else {
		err = c.decrypt(raw.Body(), rsp)
	}

	return
}

//go:inline
func (c *Client) decrypt(raw []byte, _rsp any) (err error) {
	__rsp := new(rsp)
	if err = json.Unmarshal(raw, __rsp); nil != err {
		return
	}

	// 解密
	var block cipher.Block
	if key, de := c.key.Decrypt(rand.Reader, []byte(__rsp.Key), sm2.NewPlainDecrypterOpts(sm2.C1C2C3)); nil != de {
		err = de
	} else {
		block, err = sm4.NewCipher(key)
	}
	if nil != err {
		return
	}

	var decrypted []byte
	block.Decrypt(decrypted, []byte(__rsp.Data))
	err = json.Unmarshal(decrypted, _rsp)

	return
}
