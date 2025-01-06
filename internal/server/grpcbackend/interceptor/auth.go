package interceptor

import (
	"context"
	"gophkeeper/internal/server/auth"
	"gophkeeper/pkg/constants"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Checks auth token passed from context, returns new context with user id embedded
func authContext(secretKey []byte, ctx context.Context) (context.Context, error) {
	// Get token from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to extract metadata")
	}

	values := md.Get(constants.AccessTokenHeader)

	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing access token")
	}

	tokenText := values[0]

	// Parse and verify JWT token
	tokenMap, err := auth.VerifyToken(tokenText, secretKey)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to verify token: %s", err.Error())
	}

	// Extract user from claims
	uid, ok := tokenMap["user_id"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no user id in claims")
	}

	// Cast to float from claims
	userID, ok := uid.(float64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid user id in claims")
	}

	// Store user ID in context
	ctx = context.WithValue(ctx, constants.CtxUserIDKey, uint64(userID))

	return ctx, nil
}

// Unary auth interceptor checks provided in metadata token
func Authentication(secretKey []byte) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		// Allow login and register methods
		if strings.Contains(info.FullMethod, "RegisterV1") || strings.Contains(info.FullMethod, "LoginV1") {
			return handler(ctx, req)
		}

		var err error

		ctx, err = authContext(secretKey, ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
