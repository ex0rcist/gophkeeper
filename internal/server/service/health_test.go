package service

import (
	"context"
	"errors"
	"testing"

	"gophkeeper/internal/server/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockServerStorage is a mock implementation of ServerStorage.
type MockServerStorage struct {
	mock.Mock
}

func (m *MockServerStorage) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockServerStorage) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockServerStorage) String() string {
	return "MockServerStorage"
}

func TestHealthService_Ping(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockStorage := new(MockServerStorage)
		mockStorage.On("Ping", mock.Anything).Return(nil)

		healthService := NewHealthService(HealthManagerDependencies{
			Storage: mockStorage,
		})

		ctx := context.Background()
		err := healthService.Ping(ctx)

		assert.NoError(t, err)
		mockStorage.AssertCalled(t, "Ping", mock.Anything)
	})

	t.Run("Storage is not Pingable", func(t *testing.T) {
		mockStorage := &MockServerStorage{}
		mockStorage.On("Ping", mock.Anything).Return(entities.ErrStorageUnpingable)

		healthService := NewHealthService(HealthManagerDependencies{
			Storage: mockStorage,
		})

		ctx := context.Background()
		err := healthService.Ping(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), entities.ErrStorageUnpingable.Error())
	})

	t.Run("Ping fails", func(t *testing.T) {
		mockStorage := new(MockServerStorage)
		mockStorage.On("Ping", mock.Anything).Return(errors.New("ping error"))

		healthService := NewHealthService(HealthManagerDependencies{
			Storage: mockStorage,
		})

		ctx := context.Background()
		err := healthService.Ping(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ping error")
		mockStorage.AssertCalled(t, "Ping", mock.Anything)
	})
}
