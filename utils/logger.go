// utils implements utilities used all over the project
package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitializeLogger configures and initialized the logger
func InitializeLogger(loglevel string) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	switch loglevel {
	case "INFO":
		config.Level.SetLevel(zap.InfoLevel)
	case "ERROR":
		config.Level.SetLevel(zap.ErrorLevel)
	default:
		config.Level.SetLevel(zap.InfoLevel)
	}
	Logger, _ = config.Build()
	defer Logger.Sync()
}
