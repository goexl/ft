package ft

type req struct {
	// 公钥
	PublicKey string `json:"publicKey,omitempty"`
	// 签名后的数据
	Data []byte `json:"requestData,omitempty"`
	// 令牌
	Token string `json:"token,omitempty"`
}
