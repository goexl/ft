package ft

var (
	_           = Iv
	_ newOption = (*optionIv)(nil)
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

func (i *optionIv) applyNew(options *newOptions) {
	options.iv = []byte(i.iv)
}
