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

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	cfg.EncoderConfig.CallerKey = ""

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

type ZapLogLevel string
