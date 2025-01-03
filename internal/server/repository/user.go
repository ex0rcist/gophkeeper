package repository

import (
	"context"
	"gophkeeper/pkg/models"
)

//go:generate mockgen -source user.go -destination mocks/mock_user.go -package repository
type UsersRepository interface {
	Create(ctx context.Context, user models.User) (int, error)
	GetUserByID(ctx context.Context, ID int) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
}
