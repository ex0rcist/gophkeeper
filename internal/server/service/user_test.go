package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"gophkeeper/internal/server/entities"
	"gophkeeper/internal/server/utils"
	"gophkeeper/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUsersRepository is a mock implementation of UsersRepository.
type MockUsersRepository struct {
	mock.Mock
}

func (m *MockUsersRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUsersRepository) GetUserByID(ctx context.Context, ID int) (*models.User, error) {
	args := m.Called(ctx, ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUsersRepository) Create(ctx context.Context, user models.User) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func TestUsersService_RegisterUser(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockUsersRepository)

		service := NewUsersService(UsersManagerDependencies{
			Repo: mockRepo,
		})

		mockRepo.On("GetUserByLogin", ctx, "testuser").Return(nil, entities.ErrUserNotFound)
		mockRepo.On("Create", ctx, mock.Anything).Return(1, nil)

		user, err := service.RegisterUser(ctx, "testuser", "password")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Login)
		mockRepo.AssertCalled(t, "GetUserByLogin", ctx, "testuser")
		mockRepo.AssertCalled(t, "Create", ctx, mock.Anything)
	})

	t.Run("User Already Exists", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockUsersRepository)

		service := NewUsersService(UsersManagerDependencies{
			Repo: mockRepo,
		})

		mockRepo.On("GetUserByLogin", ctx, "testuser").Return(&models.User{Login: "testuser"}, nil)

		user, err := service.RegisterUser(ctx, "testuser", "password")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "already exists")
		mockRepo.AssertCalled(t, "GetUserByLogin", ctx, "testuser")
	})

	t.Run("Create User Failure", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockUsersRepository)

		service := NewUsersService(UsersManagerDependencies{
			Repo: mockRepo,
		})

		mockRepo.On("GetUserByLogin", ctx, "testuser").Return(nil, entities.ErrUserNotFound)
		mockRepo.On("Create", ctx, mock.Anything).Return(0, errors.New("create error"))

		user, err := service.RegisterUser(ctx, "testuser", "password")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "create error")
		mockRepo.AssertCalled(t, "GetUserByLogin", ctx, "testuser")
		mockRepo.AssertCalled(t, "Create", ctx, mock.Anything)
	})
}

func TestUsersService_LoginUser(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockUsersRepository)

		service := NewUsersService(UsersManagerDependencies{
			Repo: mockRepo,
		})

		pw, _ := utils.HashPassword("password")
		mockRepo.On("GetUserByLogin", ctx, "testuser").Return(&models.User{Login: "testuser", Password: pw}, nil)

		user, err := service.LoginUser(ctx, "testuser", "password")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Login)
		mockRepo.AssertCalled(t, "GetUserByLogin", ctx, "testuser")
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockUsersRepository)

		service := NewUsersService(UsersManagerDependencies{
			Repo: mockRepo,
		})

		mockRepo.On("GetUserByLogin", ctx, "testuser").Return(&models.User{Login: "testuser", Password: "$2a$12$EXAMPLE"}, nil)

		user, err := service.LoginUser(ctx, "testuser", "wrongpassword")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "bad auth credentials")
		mockRepo.AssertCalled(t, "GetUserByLogin", ctx, "testuser")
	})

	t.Run("User Not Found", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockUsersRepository)

		service := NewUsersService(UsersManagerDependencies{
			Repo: mockRepo,
		})

		mockRepo.On("GetUserByLogin", ctx, "testuser").Return(nil, sql.ErrNoRows)

		user, err := service.LoginUser(ctx, "testuser", "password")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "bad auth credentials")
		mockRepo.AssertCalled(t, "GetUserByLogin", ctx, "testuser")
	})

	t.Run("Repository Error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := new(MockUsersRepository)

		service := NewUsersService(UsersManagerDependencies{
			Repo: mockRepo,
		})

		mockRepo.On("GetUserByLogin", ctx, "testuser").Return(nil, errors.New("repo error"))

		user, err := service.LoginUser(ctx, "testuser", "password")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "repo error")
		mockRepo.AssertCalled(t, "GetUserByLogin", ctx, "testuser")
	})
}
