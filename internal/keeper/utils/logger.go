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
		Encoding:         "console",                            // Формат логов (можно использовать "console" или "json")
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel), // Уровень логирования
		OutputPaths:      []string{"debug.log"},                // Путь к файлу для логирования
		ErrorOutputPaths: []string{"stderr"},                   // Куда писать ошибки логгера
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),    // Конфигурация форматирования
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
