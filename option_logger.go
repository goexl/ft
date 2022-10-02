package ft

import (
	"github.com/goexl/simaqian"
)

var (
	_           = Logger
	_ newOption = (*optionLogger)(nil)
)

type optionLogger struct {
	logger simaqian.Logger
}

// Logger 应用
func Logger(logger simaqian.Logger) *optionLogger {
	return &optionLogger{
		logger: logger,
	}
}

func (h *optionLogger) applyNew(options *newOptions) {
	options.logger = h.logger
}
