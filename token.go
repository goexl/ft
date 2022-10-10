package ft

import (
	"time"
)

type (
	// 52号文获取令牌请求
	tokenReq struct {
		// 应用编号
		AppId string `json:"appId,omitempty"`
		// 用户名
		AppKey string `json:"appKey,omitempty"`
		// 密码
		AppSecret string `json:"appSecret,omitempty"`
	}

	// 52号文获取令牌响应
	tokenRsp struct {
		// 令牌
		Token string `json:"token"`
		// 过期时间
		Expires int64 `json:"expiresTime"`
	}
)

func (tr *tokenRsp) Expired() bool {
	return time.UnixMilli(tr.Expires).After(time.Now())
}
