package ft

var (
	_        = App
	_ option = (*optionApp)(nil)
)

type optionApp struct {
	id     string
	key    string
	secret string
}

// App 应用
func App(id string, key string, secret string) *optionApp {
	return &optionApp{
		id:     id,
		key:    key,
		secret: secret,
	}
}

func (a *optionApp) apply(options *options) {
	options.id = a.id
	options.key = a.key
	options.secret = a.secret
}
