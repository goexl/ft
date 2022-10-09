package ft

var (
	_           = Addr
	_ newOption = (*optionAddr)(nil)
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

func (a *optionAddr) applyNew(options *newOptions) {
	options.addr = a.addr
}
