package interceptor

import (
	"context"
	"gophkeeper/internal/server/auth"
	"gophkeeper/pkg/constants"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestAuthentication(t *testing.T) {
	secretKey := "test"

	authInterceptor := Authentication([]byte(secretKey))

	handler := func(ctx context.Context, req any) (any, error) {
		var auth bool
		uid := ctx.Value(constants.CtxUserIDKey)

		if uid != nil {
			auth = true
		}

		return auth, nil
	}

	t.Run("auth methods skip", func(t *testing.T) {
		res, err := authInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "..../RegisterV1",
		}, handler)

		assert.NoError(t, err)
		assert.False(t, res.(bool))
	})

	t.Run("valid auth", func(t *testing.T) {
		userID := uint64(111)
		token, err := auth.CreateToken(int(userID), time.Now().Add(time.Hour), []byte(secretKey))
		require.NoError(t, err)

		md := metadata.New(map[string]string{
			constants.AccessTokenHeader: token,
		})

		mdCtx := metadata.NewIncomingContext(context.Background(), md)

		res, err := authInterceptor(mdCtx, nil, &grpc.UnaryServerInfo{
			FullMethod: "SomeMethod",
		}, handler)

		assert.NoError(t, err)
		assert.True(t, res.(bool))
	})

	t.Run("failed auth", func(t *testing.T) {
		_, err := authInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "SomeMethod",
		}, handler)

		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = unable to extract metadata")
	})
}

func TestAuthContext(t *testing.T) {
	secretKey := []byte("test")

	t.Run("missing access token", func(t *testing.T) {
		md := metadata.New(nil)

		mdCtx := metadata.NewIncomingContext(context.Background(), md)

		_, err := authContext(secretKey, mdCtx)

		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = missing access token")
	})

	t.Run("failed to verify", func(t *testing.T) {
		md := metadata.New(map[string]string{
			constants.AccessTokenHeader: "invalid",
		})

		mdCtx := metadata.NewIncomingContext(context.Background(), md)

		_, err := authContext(secretKey, mdCtx)

		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = failed to verify token: token contains an invalid number of segments")
	})
}
