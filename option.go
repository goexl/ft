package ft

import (
	"github.com/goexl/simaqian"
)

type (
	option interface {
		apply(options *options)
	}

	options struct {
		host string

		id     string
		key    string
		secret string
		iv     []byte

		logger simaqian.Logger
	}
)

//go:inline
func apply(opts ...option) (_options *options) {
	_options = defaultOptions()
	for _, opt := range opts {
		opt.apply(_options)
	}

	return
}

func defaultOptions() *options {
	return &options{
		host: `https://202.61.91.57:8092`,
		iv:   []byte(`UI8wC9fW6cFh9SOS`),
	}
}
