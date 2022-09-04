package utils

import (
	"log"

	"github.com/JeremyCurmi/simpleBank/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	InfoLevel  = "info"
	WarnLevel  = "warn"
	DebugLevel = "debug"
)

func parseLogLevel(level string) zapcore.Level {
	switch level {
	case WarnLevel:
		return zapcore.WarnLevel
	case DebugLevel:
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

func NewLogger() *zap.Logger {
	logLevel := zap.NewAtomicLevel()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("could not create logger: %v", err)
	}

	defer func() {
		_ = logger.Sync()
	}()

	logLevel.SetLevel(parseLogLevel(*config.LogLevel))

	return logger
}