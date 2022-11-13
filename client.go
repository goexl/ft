package ft

import (
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/emmansun/gmsm/sm2"
	"github.com/emmansun/gmsm/sm3"
	"github.com/emmansun/gmsm/sm4"
	"github.com/go-resty/resty/v2"
	"github.com/goexl/exc"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/xiren"
)

var _ = New

// Client 客户端
type Client struct {
	key *sm2.PrivateKey
	hex string

	options *newOptions
	tokens  map[string]*tokenRsp
	keys    map[string]string
}

// New 创建客户端
func New(opts ...newOption) (client *Client, err error) {
	client = new(Client)
	client.tokens = make(map[string]*tokenRsp)
	client.keys = make(map[string]string)
	client.options = defaultNewOptions()
	for _, opt := range opts {
		opt.applyNew(client.options)
	}

	// 不验证证书有效性
	client.options.http.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	if client.key, err = sm2.GenerateKey(rand.Reader); nil == err {
		client.hex = client.publicKeyToHex()
	}

	return
}

func (c *Client) request(api string, req any, rsp any, opts ...option) (err error) {
	_options := apply(opts...)
	_req := new(request)
	_req.PublicKey = c.hex

	// 注入令牌
	if getTokenApi != api && getPublicKeyApi != api {
		_req.Token, err = c.Token(opts...)
	}
	if nil != err {
		return
	}

	// 注入公钥
	if data, je := json.Marshal(req); nil != je {
		err = je
	} else if getPublicKeyApi != api {
		err = c.auth(_req, data, _options, opts...)
	} else {
		_req.Data = string(data)
	}

	hr := c.options.http.R()
	hr.SetBody(_req)
	err = c.post(api, hr, rsp, _options)

	return
}

//go:inline
func (c *Client) sendfile(api string, file string, req any, rsp any, opts ...option) (err error) {
	hr := c.options.http.R()
	if form, fe := gox.StructToForm(req); nil != fe {
		err = fe
	} else if form[`token`], err = c.Token(opts...); nil == err {
		form[`publicKey`] = c.hex
		hr.SetFormData(form)
	}
	if nil != err {
		return
	}

	// 设置上传文件路径
	if `` != file {
		hr.SetFile(`file`, file)
	}
	err = c.post(api, hr, rsp, apply(opts...))

	return
}

//go:inline
func (c *Client) post(api string, req *resty.Request, rsp any, _options *options) (err error) {
	if err = xiren.Struct(_options); nil != err {
		return
	}

	fields := gox.Fields{
		field.String(`api`, api),
	}
	if hr, pe := req.Post(fmt.Sprintf(`%s%s`, _options.addr, api)); nil != pe {
		err = pe
		c.options.logger.Error(`发送数据出错`, fields.Connect(field.Error(err))...)
	} else if hr.IsError() {
		code := field.Int("code", hr.StatusCode())
		raw := field.String(`raw`, hr.String())
		err = exc.NewFields("大数据中心返回错误", fields.Connect(code).Connect(raw)...)
	} else {
		err = c.unmarshal(hr.Body(), rsp, _options)
	}

	return
}

//go:inline
func (c *Client) auth(req *request, data []byte, _options *options, opts ...option) (err error) {
	// 随机生成加密密钥
	key := gox.RandString(16)
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

//go:inline
func (c *Client) sign(data []byte) (sign string, err error) {
	sm := sm3.New()
	sm.Write(data)
	hash := sm.Sum(nil)
	if _sign, se := c.key.Sign(rand.Reader, hash, nil); nil != se {
		err = se
	} else {
		sign = hex.EncodeToString(_sign)
	}

	return
}

//go:inline
func (c *Client) unmarshal(raw []byte, rsp any, _options *options) (err error) {
	_rsp := new(response)
	if err = json.Unmarshal(raw, _rsp); nil != err || `` == _rsp.Data {
		return
	}

	// 解密
	if key, ke := c.decryptKey(_rsp.Key); nil != ke {
		err = ke
	} else if decrypted, de := c.cbcDecrypt(_rsp.Data, key, _options); nil != de {
		err = de
	} else if 130 == len(decrypted) {
		// 处理四川站返回的公钥不是JSON格式的问题
		decrypted = []byte(fmt.Sprintf(`{"publicKey": "%s"}`, decrypted))
		err = json.Unmarshal(decrypted, rsp)
	} else {
		err = json.Unmarshal(decrypted, rsp)
	}

	return
}

//go:inline
func (c *Client) decryptKey(key string) (decrypted []byte, err error) {
	if decoded, ke := hex.DecodeString(key); nil != ke {
		err = ke
	} else {
		decrypted, err = c.key.Decrypt(rand.Reader, decoded, sm2.NewPlainDecrypterOpts(sm2.C1C2C3))
	}

	return
}

//go:inline
func (c *Client) encryptKey(pk string, key string) (encrypted string, err error) {
	var _encrypted []byte
	if _pk, pe := c.hexToPublicKey(pk); nil != pe {
		err = pe
	} else {
		opts := sm2.NewPlainEncrypterOpts(sm2.MarshalHybrid, sm2.C1C2C3)
		_encrypted, err = sm2.Encrypt(rand.Reader, _pk, []byte(key), opts)
	}
	if nil == err {
		encrypted = hex.EncodeToString(_encrypted)
	}

	return
}

//go:inline
func (c *Client) cbcDecrypt(raw string, key []byte, _options *options) (decrypted []byte, err error) {
	var block cipher.Block
	if decoded, de := base64.StdEncoding.DecodeString(raw); nil != de {
		err = de
	} else {
		decrypted = make([]byte, len(decoded))
		copy(decrypted, decoded)
		block, err = sm4.NewCipher(key)
	}
	if nil != err {
		return
	}

	cbc := cipher.NewCBCDecrypter(block, _options.iv)
	cbc.CryptBlocks(decrypted, decrypted)

	pkcs := newPkcs5(sm4.BlockSize)
	decrypted, err = pkcs.Unpad(decrypted)

	return
}

//go:inline
func (c *Client) cbcEncrypt(raw []byte, key string, _options *options) (encrypted string, err error) {
	if block, be := sm4.NewCipher([]byte(key)); nil != be {
		err = be
	} else {
		pkcs := newPkcs5(sm4.BlockSize)
		pad := pkcs.Pad(raw)
		_encrypted := make([]byte, len(pad))
		copy(_encrypted, pad)

		cbc := cipher.NewCBCEncrypter(block, _options.iv)
		cbc.CryptBlocks(_encrypted, _encrypted)

		encrypted = base64.StdEncoding.EncodeToString(_encrypted)
	}

	return
}

//go:inline
func (c *Client) privateKeyToHex() string {
	return c.key.D.Text(16)
}

//go:inline
func (c *Client) publicKeyToHex() string {
	key := &c.key.PublicKey
	x := key.X.Bytes()
	y := key.Y.Bytes()
	if n := len(x); n < 32 {
		x = append(c.zeroByteSlice()[:32-n], x...)
	}
	if n := len(y); n < 32 {
		y = append(c.zeroByteSlice()[:32-n], y...)
	}

	var data []byte
	data = append(data, x...)
	data = append(data, y...)
	data = append([]byte{0x04}, data...)

	return strings.ToUpper(hex.EncodeToString(data))
}

func (c *Client) hexToPublicKey(_hex string) (key *ecdsa.PublicKey, err error) {
	var q []byte
	if q, err = hex.DecodeString(_hex); nil != err {
		return
	}

	if len(q) == 65 && q[0] == byte(0x04) {
		q = q[1:]
	}

	if 64 != len(q) {
		err = exc.NewMessage(`公钥未被压缩`)
	}
	if nil != err {
		return
	}

	key = new(ecdsa.PublicKey)
	key.Curve = sm2.P256()
	key.X = new(big.Int).SetBytes(q[:32])
	key.Y = new(big.Int).SetBytes(q[32:])

	return
}

//go:inline
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
