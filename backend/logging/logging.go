package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      false,
		Encoding:         "json", // console or json
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		panic("[PANIC] Failed to create a logger: " + err.Error())
	}
	return logger
}
