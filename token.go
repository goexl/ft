package ft

import (
	"time"
)

type (
	// TokenReq 52号文获取令牌请求
	TokenReq struct {
		// 应用编号
		AppId string `json:"appId"`
		// 用户名
		AppKey string `json:"appKey"`
		// 密码
		AppSecret string `json:"appSecret"`
	}

	// TokenRsp 52号文获取令牌响应
	TokenRsp struct {
		// 令牌
		Token string `json:"token"`
		// 过期时间
		Expires time.Time `json:"expiresTime"`
	}
)
