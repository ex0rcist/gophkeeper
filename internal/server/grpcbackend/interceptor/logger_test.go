package interceptor

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func TestLogger(t *testing.T) {
	var buffer bytes.Buffer

	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	writer := bufio.NewWriter(&buffer)

	memLogger := zap.New(
		zapcore.NewCore(encoder, zapcore.AddSync(writer), zapcore.DebugLevel),
	)

	logger := Logger(memLogger.Sugar())

	handler := func(ctx context.Context, req any) (any, error) {
		return nil, nil
	}

	handlerError := func(ctx context.Context, req any) (any, error) {
		return nil, errors.New("test error")
	}

	t.Run("log success", func(t *testing.T) {
		defer buffer.Reset()

		_, err := logger(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "TestMethod",
		}, handler)

		writer.Flush()

		assert.NoError(t, err)
		assert.Contains(t, buffer.String(), "TestMethod")
	})

	t.Run("log error", func(t *testing.T) {
		defer buffer.Reset()

		_, err := logger(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "TestMethod",
		}, handlerError)

		writer.Flush()

		assert.Error(t, err)
		assert.Contains(t, buffer.String(), "test error")
	})
}
