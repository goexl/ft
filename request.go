package ft

type request struct {
	baseReq

	// 加密密钥
	Key string `json:"key,omitempty"`
	// 签名
	Signature string `json:"signatureData,omitempty"`
	// 令牌
	Token string `json:"token,omitempty"`
}
