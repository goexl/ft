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
		iv     []byte
		logger simaqian.Logger
	}
)

func defaultNewOptions() *newOptions {
	return &newOptions{
		iv: []byte(`UI8wC9fW6cFh9SOS`),
	}
}
