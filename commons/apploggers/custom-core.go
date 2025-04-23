package apploggers

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

// wrapper for zap actual core to perform any required pre / post processing for message
type customCore struct {
	delegate zapcore.Core
}

func NewCustomCore(delegate zapcore.Core) zapcore.Core {
	return &customCore{delegate: delegate}
}

func (z *customCore) Enabled(l zapcore.Level) bool {
	return z.delegate.Enabled(l)
}

func (z *customCore) With(f []zapcore.Field) zapcore.Core {
	newDelegate := z.delegate.With(f)
	return NewCustomCore(newDelegate)
}

func (z *customCore) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	res := z.delegate.Check(e, ce)
	if res != nil {
		// remove all newlines from message
		res.Message = strings.ReplaceAll(e.Message, "\n", "<LINEBREAK>")
	}
	return res
}

func (z *customCore) Write(e zapcore.Entry, f []zapcore.Field) error {
	return z.delegate.Write(e, f)
}

func (z *customCore) Sync() error {
	return z.delegate.Sync()
}
