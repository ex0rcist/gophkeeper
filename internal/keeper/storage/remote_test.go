package storage

import (
	"context"
	"testing"
	"time"

	"gophkeeper/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockApiClient struct {
	mock.Mock
}

func (m *MockApiClient) Register(ctx context.Context, login string, password string) (string, error) {
	args := m.Called(ctx, login, password)
	return args.String(0), args.Error(1)
}

func (m *MockApiClient) Login(ctx context.Context, login string, password string) (string, error) {
	args := m.Called(ctx, login, password)
	return args.String(0), args.Error(1)
}

func (m *MockApiClient) LoadSecrets(ctx context.Context) ([]*models.Secret, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Secret), args.Error(1)
}

func (m *MockApiClient) LoadSecret(ctx context.Context, ID uint64) (*models.Secret, error) {
	args := m.Called(ctx, ID)
	return args.Get(0).(*models.Secret), args.Error(1)
}

func (m *MockApiClient) SaveSecret(ctx context.Context, secret *models.Secret) error {
	args := m.Called(ctx, secret)
	return args.Error(0)
}

func (m *MockApiClient) DeleteSecret(ctx context.Context, ID uint64) error {
	args := m.Called(ctx, ID)
	return args.Error(0)
}

func (m *MockApiClient) SetToken(token string) {
	m.Called(token)
}

func (m *MockApiClient) GetToken() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockApiClient) SetPassword(password string) {
	m.Called(password)
}

func (m *MockApiClient) GetPassword() string {
	args := m.Called()
	return args.String(0)
}

func TestRemoteStorage(t *testing.T) {
	mockClient := new(MockApiClient)
	encrypter := &MockEncrypter{}
	password := "testpassword"

	mockClient.On("GetPassword").Return(password)
	store, err := NewRemoteStorage(mockClient, encrypter)
	assert.NoError(t, err)

	secret := &models.Secret{
		ID:         1,
		Title:      "Test Secret",
		Metadata:   "metadata",
		SecretType: "credential",
		Payload:    []byte("payload"),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	t.Run("Create Secret", func(t *testing.T) {
		mockClient.On("SaveSecret", mock.Anything, secret).Return(nil)

		err := store.Create(context.Background(), secret)
		assert.NoError(t, err)
		mockClient.AssertCalled(t, "SaveSecret", mock.Anything, secret)
	})

	t.Run("Get Secret", func(t *testing.T) {
		mockClient.On("LoadSecret", mock.Anything, secret.ID).Return(secret, nil)

		result, err := store.Get(context.Background(), secret.ID)
		assert.NoError(t, err)
		assert.Equal(t, secret.Title, result.Title)
		mockClient.AssertCalled(t, "LoadSecret", mock.Anything, secret.ID)
	})

	t.Run("Update Secret", func(t *testing.T) {
		updatedSecret := &models.Secret{
			ID:         1,
			Title:      "Updated Secret",
			Metadata:   "updated metadata",
			SecretType: "credential",
			Payload:    []byte("updated payload"),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mockClient.On("SaveSecret", mock.Anything, updatedSecret).Return(nil)

		err := store.Update(context.Background(), updatedSecret)
		assert.NoError(t, err)
		mockClient.AssertCalled(t, "SaveSecret", mock.Anything, updatedSecret)
	})

	t.Run("Delete Secret", func(t *testing.T) {
		mockClient.On("DeleteSecret", mock.Anything, secret.ID).Return(nil)

		err := store.Delete(context.Background(), secret.ID)
		assert.NoError(t, err)
		mockClient.AssertCalled(t, "DeleteSecret", mock.Anything, secret.ID)
	})

	t.Run("Get All Secrets", func(t *testing.T) {
		secrets := []*models.Secret{
			{ID: 1, Title: "Secret 1"},
			{ID: 2, Title: "Secret 2"},
		}
		mockClient.On("LoadSecrets", mock.Anything).Return(secrets, nil)

		result, err := store.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Len(t, result, len(secrets))
		mockClient.AssertCalled(t, "LoadSecrets", mock.Anything)
	})
}
