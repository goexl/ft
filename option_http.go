package ft

import (
	"github.com/go-resty/resty/v2"
)

var (
	_           = Http
	_ newOption = (*optionHttp)(nil)
)

type optionHttp struct {
	http *resty.Client
}

// Http Http客户端
func Http(http *resty.Client) *optionHttp {
	return &optionHttp{
		http: http,
	}
}

func (h *optionHttp) applyNew(options *newOptions) {
	options.http = h.http
}
