package ft

type rsp struct {
	// 授权码
	Key string `json:"key"`
	// 签名数据
	SignatureData string `json:"signatureData"`
	// 数据
	Data string `json:"data"`
}
