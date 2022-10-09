package ft

import (
	"github.com/go-resty/resty/v2"
	"github.com/goexl/simaqian"
)

var _ = NewNewOptions

type (
	newOption interface {
		applyNew(options *newOptions)
	}

	newOptions struct {
		http   *resty.Client
		logger simaqian.Logger
	}
)

// NewNewOptions 创建选项
func NewNewOptions(opts ...newOption) []newOption {
	return opts
}

func defaultNewOptions() *newOptions {
	return &newOptions{
		http:   resty.New(),
		logger: simaqian.Must(),
	}
}
