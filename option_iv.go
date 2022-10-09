package ft

var (
	_        = Iv
	_ option = (*optionIv)(nil)
)

type optionIv struct {
	iv string
}

// Iv 插值
func Iv(iv string) *optionIv {
	return &optionIv{
		iv: iv,
	}
}

func (i *optionIv) apply(options *options) {
	options.iv = []byte(i.iv)
}
