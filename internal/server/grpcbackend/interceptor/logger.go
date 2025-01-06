package interceptor

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Logs gRPC requests for unary server requests
func Logger(logger *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		t := time.Now()

		// Handle request
		res, err := handler(ctx, req)

		// Get headers
		requestID, clientIP := extractMetaData(ctx)

		// Count request duration in ms
		latency := time.Since(t)
		miliSeconds := fmt.Sprintf("%d ms", latency.Milliseconds())

		// Log body
		logParams := []interface{}{
			"rid", requestID,
			"method", info.FullMethod,
			"duration", miliSeconds,
			"remote_addr", clientIP,
		}

		// Log error or just log general info
		if err != nil {
			logParams = append(logParams, "error", err)
			logger.Errorln(logParams...)
		} else {
			logger.Infoln(logParams...)
		}

		return res, err
	}
}
