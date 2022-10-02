package ft

type (
	// PublicKeyReq 52号文公钥请求
	PublicKeyReq struct {
		// 应用编号
		AppId string `json:"appId"`
	}

	// PublicKeyRsp 52号文公钥响应
	PublicKeyRsp struct {
		// 公钥
		Key string `json:"publicKey"`
	}
)
