package ft

type baseReq struct {
	// 公钥
	PublicKey string `json:"publicKey,omitempty"`
	// 签名后的数据
	Data string `json:"requestData,omitempty"`
}
