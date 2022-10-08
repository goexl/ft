package ft

type req struct {
	baseReq

	// 令牌
	Token string `json:"token,omitempty"`
}
