package grpchandlers

import (
	"context"
	"testing"

	"gophkeeper/internal/server/entities"
	"gophkeeper/pkg/constants"
	"gophkeeper/pkg/models"
	"gophkeeper/pkg/proto/keeper/grpcapi"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MockSecretsManager is a mock implementation of the SecretsManager interface.
type MockSecretsManager struct {
	mock.Mock
}

func (m *MockSecretsManager) CreateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error) {
	args := m.Called(ctx, secret)
	return args.Get(0).(*models.Secret), args.Error(1)
}

func (m *MockSecretsManager) UpdateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error) {
	args := m.Called(ctx, secret)
	return args.Get(0).(*models.Secret), args.Error(1)
}

func (m *MockSecretsManager) GetSecret(ctx context.Context, id, userID uint64) (*models.Secret, error) {
	args := m.Called(ctx, id, userID)
	return args.Get(0).(*models.Secret), args.Error(1)
}

func (m *MockSecretsManager) GetUserSecrets(ctx context.Context, userID uint64) (models.Secrets, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(models.Secrets), args.Error(1)
}

func (m *MockSecretsManager) DeleteSecret(ctx context.Context, id, userID uint64) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func TestSecretsServer_SaveUserSecretV1(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.CtxUserIDKey, uint64(1))

	secretProto := &grpcapi.Secret{
		Id:      0,
		Title:   "test",
		Payload: []byte("data"),
	}

	t.Run("Create new secret", func(t *testing.T) {
		mockSecretsManager := new(MockSecretsManager)
		secretsServer := NewSecretsServer(SecretsServerDependencies{
			SecretsManager: mockSecretsManager,
		})

		mockSecretsManager.On("CreateSecret", ctx, mock.Anything).Return(&models.Secret{ID: 1}, nil)

		request := &grpcapi.SaveUserSecretRequestV1{
			Secret: secretProto,
		}

		response, err := secretsServer.SaveUserSecretV1(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		mockSecretsManager.AssertCalled(t, "CreateSecret", ctx, mock.Anything)
	})

	t.Run("Update existing secret", func(t *testing.T) {
		mockSecretsManager := new(MockSecretsManager)
		secretsServer := NewSecretsServer(SecretsServerDependencies{
			SecretsManager: mockSecretsManager,
		})

		secretProto.Id = 1
		mockSecretsManager.On("UpdateSecret", ctx, mock.Anything).Return(&models.Secret{ID: 1}, nil)

		request := &grpcapi.SaveUserSecretRequestV1{
			Secret: secretProto,
		}

		response, err := secretsServer.SaveUserSecretV1(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		mockSecretsManager.AssertCalled(t, "UpdateSecret", ctx, mock.Anything)
	})
}

func TestSecretsServer_GetUserSecretV1(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.CtxUserIDKey, uint64(1))

	t.Run("Get existing secret", func(t *testing.T) {
		mockSecretsManager := new(MockSecretsManager)
		secretsServer := NewSecretsServer(SecretsServerDependencies{
			SecretsManager: mockSecretsManager,
		})

		mockSecretsManager.On("GetSecret", ctx, uint64(1), uint64(1)).Return(&models.Secret{
			ID:      1,
			Title:   "test",
			Payload: []byte{},
			UserID:  1,
		}, nil)

		request := &grpcapi.GetUserSecretRequestV1{
			Id: 1,
		}

		response, err := secretsServer.GetUserSecretV1(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "test", response.Secret.Title)
		mockSecretsManager.AssertCalled(t, "GetSecret", ctx, uint64(1), uint64(1))
	})

	t.Run("Secret not found", func(t *testing.T) {
		mockSecretsManager := new(MockSecretsManager)
		secretsServer := NewSecretsServer(SecretsServerDependencies{
			SecretsManager: mockSecretsManager,
		})

		mockSecretsManager.On("GetSecret", ctx, uint64(2), uint64(1)).Return((*models.Secret)(nil), entities.ErrSecretNotFound)

		request := &grpcapi.GetUserSecretRequestV1{
			Id: 2,
		}

		response, err := secretsServer.GetUserSecretV1(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.NotFound, status.Code(err))
		mockSecretsManager.AssertCalled(t, "GetSecret", ctx, uint64(2), uint64(1))
	})
}

func TestSecretsServer_GetUserSecretsV1(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.CtxUserIDKey, uint64(1))

	mockSecretsManager := new(MockSecretsManager)
	secretsServer := NewSecretsServer(SecretsServerDependencies{
		SecretsManager: mockSecretsManager,
	})

	t.Run("Get user secrets", func(t *testing.T) {
		mockSecretsManager.On("GetUserSecrets", ctx, uint64(1)).Return(models.Secrets{
			{
				ID:      1,
				Title:   "secret1",
				Payload: []byte("data1"),
				UserID:  1,
			},
			{
				ID:      2,
				Title:   "secret1",
				Payload: []byte("data2"),
				UserID:  1,
			},
		}, nil)

		response, err := secretsServer.GetUserSecretsV1(ctx, &emptypb.Empty{})

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Secrets, 2)
		mockSecretsManager.AssertCalled(t, "GetUserSecrets", ctx, uint64(1))
	})
}

func TestSecretsServer_DeleteUserSecretV1(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.CtxUserIDKey, uint64(1))

	mockSecretsManager := new(MockSecretsManager)
	secretsServer := NewSecretsServer(SecretsServerDependencies{
		SecretsManager: mockSecretsManager,
	})

	t.Run("Delete existing secret", func(t *testing.T) {
		mockSecretsManager.On("DeleteSecret", ctx, uint64(1), uint64(1)).Return(nil)

		request := &grpcapi.DeleteUserSecretRequestV1{
			Id: 1,
		}

		response, err := secretsServer.DeleteUserSecretV1(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		mockSecretsManager.AssertCalled(t, "DeleteSecret", ctx, uint64(1), uint64(1))
	})

	t.Run("Secret not found", func(t *testing.T) {
		mockSecretsManager.On("DeleteSecret", ctx, uint64(2), uint64(1)).Return(entities.ErrSecretNotFound)

		request := &grpcapi.DeleteUserSecretRequestV1{
			Id: 2,
		}

		response, err := secretsServer.DeleteUserSecretV1(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.NotFound, status.Code(err))
		mockSecretsManager.AssertCalled(t, "DeleteSecret", ctx, uint64(2), uint64(1))
	})
}
