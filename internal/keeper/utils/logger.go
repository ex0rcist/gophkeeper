package utils

import (
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type LoggerDependencies struct {
	dig.In

	Level ZapLogLevel
}

func NewZapLogger(deps LoggerDependencies) (*zap.SugaredLogger, error) {
	lvl, err := zap.ParseAtomicLevel(string(deps.Level))
	if err != nil {
		return nil, err
	}

	// cfg := zap.NewDevelopmentConfig()
	cfg := zap.Config{
		Encoding:         "console",                            // log format ("console" or "json")
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel), // log level
		OutputPaths:      []string{"debug.log"},                // log path
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
	}

	cfg.Level = lvl
	cfg.EncoderConfig.CallerKey = ""

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

type ZapLogLevel string
