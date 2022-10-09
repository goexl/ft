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

//go:inline
func (c *Client) request(api string, _req any, rsp any) (err error) {
	fr := new(req)
	fr.PublicKey = c.hex

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

	hr := c.options.http.R()
	hr.SetBody(fr)
	err = c.post(api, hr, rsp)

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
	err = c.post(api, hr, rsp)

	return
}

//go:inline
func (c *Client) post(api string, req *resty.Request, rsp any) (err error) {
	if raw, reqErr := req.Post(fmt.Sprintf(`%s%s`, c.options.addr, api)); nil != reqErr {
		err = reqErr
		c.options.logger.Error(`发送到省大数据中心出错`, field.String(`api`, api), field.Error(err))
	} else if raw.IsError() {
		c.options.logger.Warn(`发送到省大数据中心出错`, field.String(`api`, api), field.String(`raw`, raw.String()))
	} else {
		err = c.unmarshal(raw.Body(), rsp)
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
func (c *Client) unmarshal(raw []byte, _rsp any) (err error) {
	__rsp := new(rsp)
	if err = json.Unmarshal(raw, __rsp); nil != err {
		return
	}

	// 解密
	if key, ke := c.decryptKey(__rsp.Key); nil != ke {
		err = ke
	} else if decrypted, de := c.cbcDecrypt(__rsp.Data, key); nil != de {
		err = de
	} else {
		err = json.Unmarshal(decrypted, _rsp)
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
func (c *Client) cbcDecrypt(raw string, key []byte) (decrypted []byte, err error) {
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

	cbc := cipher.NewCBCDecrypter(block, c.options.iv)
	cbc.CryptBlocks(decrypted, decrypted)

	pkcs := newPkcs5(sm4.BlockSize)
	decrypted, err = pkcs.Unpad(decrypted)

	return
}

//go:inline
func (c *Client) cbcEncrypt(raw []byte, key string) (encrypted string, err error) {
	if block, be := sm4.NewCipher([]byte(key)); nil != be {
		err = be
	} else {
		pkcs := newPkcs5(sm4.BlockSize)
		pad := pkcs.Pad(raw)
		_encrypted := make([]byte, len(pad))
		copy(_encrypted, pad)

		cbc := cipher.NewCBCEncrypter(block, c.options.iv)
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
