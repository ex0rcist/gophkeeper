package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gophkeeper/internal/server/entities"
	"gophkeeper/internal/server/repository"

	"gophkeeper/pkg/models"

	"go.uber.org/dig"
)

//go:generate mockgen -source secret.go -destination mocks/mock_secret.go -package service

var _ SecretsManager = SecretsService{}

// Interface for secrets service
type SecretsManager interface {
	GetSecret(ctx context.Context, ID uint64, userID uint64) (*models.Secret, error)
	GetUserSecrets(ctx context.Context, userID uint64) (models.Secrets, error)
	CreateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error)
	UpdateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error)
	DeleteSecret(ctx context.Context, ID uint64, userID uint64) error
}

type SecretsManagerDependencies struct {
	dig.In
	Repo repository.SecretsRepository
}

// Secrets service implementation
type SecretsService struct {
	repo repository.SecretsRepository
}

// Create new secret service
func NewSecretsService(deps SecretsManagerDependencies) *SecretsService {
	return &SecretsService{repo: deps.Repo}
}

// Returns decrypted secret
func (s SecretsService) GetSecret(ctx context.Context, secretID uint64, userID uint64) (*models.Secret, error) {
	secret, err := s.repo.GetSecret(ctx, secretID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entities.ErrorSecretNotFound(secretID)
	}

	if err != nil {
		return nil, err
	}

	return secret, nil
}

// Get user's secrets list
func (s SecretsService) GetUserSecrets(ctx context.Context, userID uint64) (models.Secrets, error) {
	secrets, err := s.repo.GetUserSecrets(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(secrets) == 0 {
		return nil, entities.ErrNoSecrets
	}

	return secrets, nil
}

// Try create secret
func (s SecretsService) CreateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error) {
	var err error

	secret.ID, err = s.repo.Create(ctx, secret)

	if err != nil {
		return nil, fmt.Errorf("failed to create secret: %w", err)
	}

	return secret, nil
}

// Try update secret
func (s SecretsService) UpdateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error) {
	var err error

	err = s.repo.Update(ctx, secret)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entities.ErrorSecretNotFound(secret.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to store secret: %w", err)
	}

	return secret, nil
}

// Delete secret
func (s SecretsService) DeleteSecret(ctx context.Context, secretID uint64, userID uint64) error {
	var err error

	err = s.repo.Delete(ctx, secretID, userID)
	return err
}
