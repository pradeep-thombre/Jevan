package apploggers

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// function to create new zap logger
// logger will write to stdout only for now. hence commenting write to file
func NewZapLogger() *zap.Logger {
	encoder := getEncoder()
	core := zapcore.NewTee(
		NewCustomCore(zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)),
	)
	allLogger := zap.New(core, zap.AddCaller())
	return allLogger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
