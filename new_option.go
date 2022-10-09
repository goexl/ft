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
		iv     []byte
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
		iv:     []byte(`UI8wC9fW6cFh9SOS`),
	}
}
