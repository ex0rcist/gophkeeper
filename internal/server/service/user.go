package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go.uber.org/dig"
	"gophkeeper/internal/server/entities"
	"gophkeeper/internal/server/repository"
	"gophkeeper/internal/server/utils"
	"gophkeeper/pkg/models"
)

//go:generate mockgen -source user.go -destination mocks/mock_user.go -package service

var _ UsersManager = UsersService{}

// User service interface
type UsersManager interface {
	RegisterUser(ctx context.Context, login string, password string) (*models.User, error)
	LoginUser(ctx context.Context, login string, password string) (*models.User, error)
}

type UsersManagerDependencies struct {
	dig.In
	Repo repository.UsersRepository
}

// User service implementation
type UsersService struct {
	repo repository.UsersRepository
}

// Create new UserService
func NewUsersService(deps UsersManagerDependencies) *UsersService {
	return &UsersService{repo: deps.Repo}
}

// Register new User
func (s UsersService) RegisterUser(ctx context.Context, login string, password string) (*models.User, error) {
	var newUser models.User

	// ensure we have no same login
	// TODO: transaction?
	user, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	if user != nil {
		return nil, entities.ErrorUserAlreadyExists(login)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password hash: %w", err)
	}

	// create new user
	newUser = models.User{Login: login, Password: hashedPassword}

	var newUserID int
	newUserID, err = s.repo.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	newUser.ID = newUserID

	return &newUser, nil
}

// Login user
func (s UsersService) LoginUser(ctx context.Context, login string, password string) (*models.User, error) {
	user, err := s.repo.GetUserByLogin(ctx, login)
	if errors.Is(err, sql.ErrNoRows) {
		return user, entities.ErrBadCredentials
	}

	if err != nil {
		return user, fmt.Errorf("failed to authenticate user: %w", err)
	}

	if !utils.ComparePassword(user.Password, password) {
		return user, entities.ErrBadCredentials
	}

	return user, nil
}
