package interceptor

import (
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
)

// Stream auth interceptor checks provided in metadata token
func StreamAuthentication(secretKey []byte) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Check auth token and store user id in ctx
		ctx, err := authContext(secretKey, ss.Context())
		if err != nil {
			return err
		}

		// Wrap server stream
		wrappedStream := middleware.WrapServerStream(ss)
		wrappedStream.WrappedContext = ctx

		return handler(srv, wrappedStream)
	}
}
