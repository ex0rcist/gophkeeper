package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"gophkeeper/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSecretsRepository is a mock implementation of SecretsRepository.
type MockSecretsRepository struct {
	mock.Mock
}

func (m *MockSecretsRepository) GetSecret(ctx context.Context, ID uint64, userID uint64) (*models.Secret, error) {
	args := m.Called(ctx, ID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Secret), args.Error(1)
}

func (m *MockSecretsRepository) GetUserSecrets(ctx context.Context, userID uint64) (models.Secrets, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(models.Secrets), args.Error(1)
}

func (m *MockSecretsRepository) Create(ctx context.Context, secret *models.Secret) (uint64, error) {
	args := m.Called(ctx, secret)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockSecretsRepository) Update(ctx context.Context, secret *models.Secret) error {
	args := m.Called(ctx, secret)
	return args.Error(0)
}

func (m *MockSecretsRepository) Delete(ctx context.Context, ID uint64, userID uint64) error {
	args := m.Called(ctx, ID, userID)
	return args.Error(0)
}

func TestSecretsService_GetSecret(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSecretsRepository)

	service := NewSecretsService(SecretsManagerDependencies{
		Repo: mockRepo,
	})

	t.Run("Success", func(t *testing.T) {
		mockSecret := &models.Secret{ID: 1, UserID: 1, Title: "Test Secret"}
		mockRepo.On("GetSecret", ctx, uint64(1), uint64(1)).Return(mockSecret, nil)

		secret, err := service.GetSecret(ctx, 1, 1)

		assert.NoError(t, err)
		assert.Equal(t, mockSecret, secret)
		mockRepo.AssertCalled(t, "GetSecret", ctx, uint64(1), uint64(1))
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.On("GetSecret", ctx, uint64(2), uint64(1)).Return(nil, sql.ErrNoRows)

		secret, err := service.GetSecret(ctx, 2, 1)

		assert.Error(t, err)
		assert.Nil(t, secret)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestSecretsService_CreateSecret(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSecretsRepository)

	service := NewSecretsService(SecretsManagerDependencies{
		Repo: mockRepo,
	})

	t.Run("Success", func(t *testing.T) {
		mockSecret := &models.Secret{UserID: 1, Title: "Test Secret"}
		mockRepo.On("Create", ctx, mockSecret).Return(uint64(1), nil)

		createdSecret, err := service.CreateSecret(ctx, mockSecret)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), createdSecret.ID)
		mockRepo.AssertCalled(t, "Create", ctx, mockSecret)
	})

	t.Run("Failure", func(t *testing.T) {
		mockSecret := &models.Secret{UserID: 1, Title: "Test Secret"}
		mockRepo.On("Create", ctx, mockSecret).Return(uint64(0), errors.New("create error"))

		createdSecret, err := service.CreateSecret(ctx, mockSecret)

		assert.Error(t, err)
		assert.Nil(t, createdSecret)
		assert.Contains(t, err.Error(), "create error")
	})
}

func TestSecretsService_DeleteSecret(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSecretsRepository)

	service := NewSecretsService(SecretsManagerDependencies{
		Repo: mockRepo,
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Delete", ctx, uint64(1), uint64(1)).Return(nil)

		err := service.DeleteSecret(ctx, 1, 1)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Delete", ctx, uint64(1), uint64(1))
	})
}
