package ft

var _ = NewOptions

type (
	option interface {
		apply(options *options)
	}

	options struct {
		id     string
		key    string
		secret string
	}
)

// NewOptions 创建选项
func NewOptions(opts ...option) []option {
	return opts
}

//go:inline
func apply(opts ...option) (_options *options) {
	_options = defaultOptions()
	for _, opt := range opts {
		opt.apply(_options)
	}

	return
}

func defaultOptions() *options {
	return &options{}
}
