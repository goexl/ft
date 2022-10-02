package ft

import (
	"github.com/go-resty/resty/v2"
	"github.com/goexl/simaqian"
)

type (
	newOption interface {
		applyNew(options *newOptions)
	}

	newOptions struct {
		http   *resty.Client
		logger simaqian.Logger
	}
)

func defaultNewOptions() *newOptions {
	return new(newOptions)
}
