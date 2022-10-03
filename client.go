package ft

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/emmansun/gmsm/sm2"
	"github.com/go-resty/resty/v2"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/simaqian"
)

var _ = New

// Client 客户端
type Client struct {
	http      *resty.Client
	key       *sm2.PrivateKey
	publicHex string
	logger    simaqian.Logger
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
	// 不验证证书有效性
	client.http.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	if nil != _options.logger {
		client.logger = _options.logger
	} else {
		client.logger = simaqian.Must()
	}

	if client.key, err = sm2.GenerateKey(rand.Reader); nil != err {
		return
	}
	client.publicHex = client.keyToHex(&client.key.PublicKey)

	return
}

//go:inline
func (c *Client) request(api string, _req any, rsp any, opts ...option) (err error) {
	fr := new(req)
	fr.PublicKey = c.publicHex

	// 加密请求
	// var encrypted []byte
	if bytes, je := json.Marshal(_req); nil != je {
		err = je
	} else {
		// so := sm2.NewPlainEncrypterOpts(sm2.MarshalUncompressed, sm2.C1C2C3)
		// encrypted, err = sm2.Encrypt(rand.Reader, &c.key.PublicKey, bytes, so)
		fr.Data = string(bytes)
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
		err = c.decrypt(raw.Body(), rsp, _options)
	}

	return
}

//go:inline
func (c *Client) decrypt(raw []byte, _rsp any, _options *options) (err error) {
	__rsp := new(rsp)
	if err = json.Unmarshal(raw, __rsp); nil != err {
		return
	}

	// 解密
	var decryptedKey []byte
	if keyBytes, ke := hex.DecodeString(__rsp.Key); nil != ke {
		err = ke
	} else {
		// decryptedKey, err = sm2.Decrypt(c.key, keyBytes, sm2.C1C2C3)
		decryptedKey, err = c.key.Decrypt(rand.Reader, keyBytes, sm2.NewPlainDecrypterOpts(sm2.C1C2C3))
	}
	if nil != err {
		return
	}

	var decrypted []byte
	if decoded, de := base64.StdEncoding.DecodeString(__rsp.Data); nil != de {
		err = de
	} else {
		// decrypted, err = sm4.Sm4Cbc(decryptedKey, decoded, false)
		fmt.Println(decoded, decryptedKey)
	}
	if nil != err {
		return
	}

	if decoded, de := hex.DecodeString(string(decrypted)); nil != de {
		err = de
	} else {
		err = json.Unmarshal(decoded, _rsp)
	}

	return
}

func (c *Client) keyToHex(key *ecdsa.PublicKey) string {
	x := key.X.Bytes()
	y := key.Y.Bytes()
	if n := len(x); n < 32 {
		x = append(c.zeroByteSlice()[:32-n], x...)
	}
	if n := len(y); n < 32 {
		y = append(c.zeroByteSlice()[:32-n], y...)
	}

	var bytes []byte
	bytes = append(bytes, x...)
	bytes = append(bytes, y...)
	bytes = append([]byte{0x04}, bytes...)

	return strings.ToUpper(hex.EncodeToString(bytes))
}

func (c *Client) zeroByteSlice() []byte {
	return []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}
}
