package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// Unary interceptor which adds standart timeout to all client requests
func Timeout(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		timedCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return invoker(timedCtx, method, req, reply, cc, opts...)
	}
}
