package ft

var (
	_        = Host
	_ option = (*optionHost)(nil)
)

type optionHost struct {
	host string
}

// Host 应用
func Host(host string) *optionHost {
	return &optionHost{
		host: host,
	}
}

func (h *optionHost) apply(options *options) {
	options.host = h.host
}
