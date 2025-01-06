package grpchandlers

import (
	"context"
	"errors"
	"testing"

	"gophkeeper/internal/server/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MockHealthManager is a mock implementation of the HealthManager interface.
type MockHealthManager struct {
	mock.Mock
}

func (m *MockHealthManager) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestHealthServer_Ping(t *testing.T) {
	ctx := context.Background()

	t.Run("Ping succeeds", func(t *testing.T) {
		mockHealthManager := new(MockHealthManager)
		healthServer := NewHealthServer(HealthServerDependencies{
			HealthManager: mockHealthManager,
		})

		mockHealthManager.On("Ping", ctx).Return(nil)

		response, err := healthServer.Ping(ctx, &emptypb.Empty{})

		assert.NoError(t, err)
		assert.NotNil(t, response)
		mockHealthManager.AssertCalled(t, "Ping", ctx)
	})

	t.Run("Ping fails with ErrStorageUnpingable", func(t *testing.T) {
		mockHealthManager := new(MockHealthManager)
		healthServer := NewHealthServer(HealthServerDependencies{
			HealthManager: mockHealthManager,
		})

		mockHealthManager.On("Ping", ctx).Return(entities.ErrStorageUnpingable)

		response, err := healthServer.Ping(ctx, &emptypb.Empty{})

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Unimplemented, status.Code(err))
		assert.Contains(t, err.Error(), entities.ErrStorageUnpingable.Error())
		mockHealthManager.AssertCalled(t, "Ping", ctx)
	})

	t.Run("Ping fails with generic error", func(t *testing.T) {
		mockHealthManager := new(MockHealthManager)
		healthServer := NewHealthServer(HealthServerDependencies{
			HealthManager: mockHealthManager,
		})

		mockHealthManager.On("Ping", ctx).Return(errors.New("generic error"))

		response, err := healthServer.Ping(ctx, &emptypb.Empty{})

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Contains(t, err.Error(), "generic error")
		mockHealthManager.AssertCalled(t, "Ping", ctx)
	})
}
