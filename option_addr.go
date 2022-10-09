package ft

var (
	_        = Addr
	_ option = (*optionAddr)(nil)
)

type optionAddr struct {
	addr string
}

// Addr 地址
func Addr(addr string) *optionAddr {
	return &optionAddr{
		addr: addr,
	}
}

func (a *optionAddr) apply(options *options) {
	options.addr = a.addr
}
