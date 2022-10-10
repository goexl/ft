package ft

type (
	// 52号文公钥请求
	publicKeyReq struct {
		// 应用编号
		AppId string `json:"appId"`
	}

	// 52号文公钥响应
	publicKeyRsp struct {
		// 公钥
		Key string `json:"publicKey"`
	}
)
