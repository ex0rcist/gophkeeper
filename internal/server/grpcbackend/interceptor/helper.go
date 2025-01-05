package interceptor

import (
	"context"
	"gophkeeper/internal/server/utils"

	"google.golang.org/grpc/metadata"
)

func extractMetaData(ctx context.Context) (string, string) {
	var requestID, clientIP string

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return requestID, clientIP
	}

	if values := md.Get("X-Request-ID"); len(values) > 0 {
		requestID = values[0]
	} else {
		requestID = utils.GenerateRequestID()
	}

	if values := md.Get("X-Real-IP"); len(values) > 0 {
		clientIP = values[0]
	}

	return requestID, clientIP
}
