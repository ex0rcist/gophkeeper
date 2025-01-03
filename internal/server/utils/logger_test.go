package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestZapLoggerInit(t *testing.T) {
	deps := LoggerDependencies{Level: "debug"}
	zapLogger, err := NewZapLogger(deps)

	assert.NotNil(t, zapLogger)
	assert.NoError(t, err)

	assert.Equal(t, zapLogger.Level(), zap.DebugLevel)
}

func TestZapLoggerInitInvalidLevel(t *testing.T) {
	deps := LoggerDependencies{Level: "invalid"}
	zapLogger, err := NewZapLogger(deps)

	assert.Nil(t, zapLogger)
	assert.Error(t, err)
}
