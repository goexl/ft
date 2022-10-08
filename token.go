package ft

import (
	"time"
)

type (
	// 52号文获取令牌请求
	tokenReq struct {
		baseReq

		// 应用编号
		AppId string `json:"appId,omitempty"`
		// 用户名
		AppKey string `json:"appKey,omitempty"`
		// 密码
		AppSecret string `json:"appSecret,omitempty"`
		// 加密密钥
		Key string `json:"key,omitempty"`
		// 签名
		Signature string `json:"signatureData,omitempty"`
	}

	// 52号文获取令牌响应
	tokenRsp struct {
		// 令牌
		Token string `json:"token"`
		// 过期时间
		Expires time.Time `json:"expiresTime"`
	}
)
