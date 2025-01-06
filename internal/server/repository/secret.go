package repository

import (
	"context"

	"gophkeeper/pkg/models"
)

//go:generate mockgen -source secret.go -destination mocks/mock_secret.go -package repository
type SecretsRepository interface {
	GetSecret(ctx context.Context, secretID uint64, userID uint64) (*models.Secret, error)
	GetUserSecrets(ctx context.Context, userID uint64) (models.Secrets, error)
	Create(ctx context.Context, secret *models.Secret) (uint64, error)
	Update(ctx context.Context, secret *models.Secret) error
	Delete(ctx context.Context, secretID uint64, userID uint64) error
}
