package grpchandlers

import (
	"context"
	"testing"

	"gophkeeper/internal/server/config"
	"gophkeeper/internal/server/entities"
	"gophkeeper/pkg/models"
	"gophkeeper/pkg/proto/keeper/grpcapi"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockUsersManager is a mock implementation of the UsersManager interface.
type MockUsersManager struct {
	mock.Mock
}

func (m *MockUsersManager) RegisterUser(ctx context.Context, login, password string) (*models.User, error) {
	args := m.Called(ctx, login, password)
	user, _ := args.Get(0).(*models.User)

	return user, args.Error(1)
}

func (m *MockUsersManager) LoginUser(ctx context.Context, login, password string) (*models.User, error) {
	args := m.Called(ctx, login, password)
	user, _ := args.Get(0).(*models.User)

	return user, args.Error(1)
}

func TestUsersServer_RegisterV1(t *testing.T) {
	ctx := context.Background()

	t.Run("Register new user", func(t *testing.T) {
		mockUsersManager := new(MockUsersManager)
		mockConfig := &config.Config{SecretKey: "test-secret-key"}

		usersServer := NewUsersServer(UsersServerDependencies{
			Config:       mockConfig,
			UsersManager: mockUsersManager,
		})

		mockUsersManager.On("RegisterUser", ctx, "testuser", "password").Return(&models.User{ID: 1}, nil)

		response, err := usersServer.RegisterV1(ctx, &grpcapi.RegisterRequestV1{
			Login:    "testuser",
			Password: "password",
		})

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.AccessToken)
		mockUsersManager.AssertCalled(t, "RegisterUser", ctx, "testuser", "password")
	})

	t.Run("User already exists", func(t *testing.T) {
		mockUsersManager := new(MockUsersManager)
		mockConfig := &config.Config{SecretKey: "test-secret-key"}

		usersServer := NewUsersServer(UsersServerDependencies{
			Config:       mockConfig,
			UsersManager: mockUsersManager,
		})

		mockUsersManager.On("RegisterUser", ctx, "testuser", "password").Return(nil, entities.ErrUserAlreadyExists)

		response, err := usersServer.RegisterV1(ctx, &grpcapi.RegisterRequestV1{
			Login:    "testuser",
			Password: "password",
		})

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.AlreadyExists, status.Code(err))
		mockUsersManager.AssertCalled(t, "RegisterUser", ctx, "testuser", "password")
	})
}

func TestUsersServer_LoginV1(t *testing.T) {
	ctx := context.Background()

	t.Run("Login successful", func(t *testing.T) {
		mockUsersManager := new(MockUsersManager)
		mockConfig := &config.Config{SecretKey: "test-secret-key"}

		usersServer := NewUsersServer(UsersServerDependencies{
			Config:       mockConfig,
			UsersManager: mockUsersManager,
		})

		mockUsersManager.On("LoginUser", ctx, "testuser", "password").Return(&models.User{ID: 1}, nil)

		response, err := usersServer.LoginV1(ctx, &grpcapi.LoginRequestV1{
			Login:    "testuser",
			Password: "password",
		})

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.AccessToken)
		mockUsersManager.AssertCalled(t, "LoginUser", ctx, "testuser", "password")
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		mockUsersManager := new(MockUsersManager)
		mockConfig := &config.Config{SecretKey: "test-secret-key"}

		usersServer := NewUsersServer(UsersServerDependencies{
			Config:       mockConfig,
			UsersManager: mockUsersManager,
		})

		mockUsersManager.On("LoginUser", ctx, "testuser", "wrongpassword").Return(nil, entities.ErrBadCredentials)

		response, err := usersServer.LoginV1(ctx, &grpcapi.LoginRequestV1{
			Login:    "testuser",
			Password: "wrongpassword",
		})

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
		mockUsersManager.AssertCalled(t, "LoginUser", ctx, "testuser", "wrongpassword")
	})
}
