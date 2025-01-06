// Client gRPC interceptors
package interceptor

import (
	"context"
	"gophkeeper/pkg/constants"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Unary gRPC interceptor which adds auth token to metadata
func AddAuth(token *string, clientID uint32) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// pass request if token is empty
		if len(*token) == 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		// add access token to metadata
		md := metadata.New(map[string]string{
			constants.AccessTokenHeader: *token,
			constants.ClientIDHeader:    strconv.Itoa(int(clientID)),
		})

		mdCtx := metadata.NewOutgoingContext(ctx, md)
		return invoker(mdCtx, method, req, reply, cc, opts...)
	}
}

// Stream gRPC interceptor which adds auth token to metadata
func AddAuthStream(token *string, clientID uint64) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// pass request if token is empty
		if len(*token) == 0 {
			return streamer(ctx, desc, cc, method, opts...)
		}

		// add access token to metadata
		md := metadata.New(map[string]string{
			constants.AccessTokenHeader: *token,
			constants.ClientIDHeader:    strconv.Itoa(int(clientID)),
		})

		mdCtx := metadata.NewOutgoingContext(ctx, md)
		return streamer(mdCtx, desc, cc, method, opts...)
	}
}
