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
	GetSecret(ctx context.Context, ID uint64) (*models.Secret, error)
	GetUserSecrets(ctx context.Context, userID uint64) (models.Secrets, error)
	CreateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error)
	UpdateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error)
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
func (s SecretsService) GetSecret(ctx context.Context, secretID uint64) (*models.Secret, error) {
	secret, err := s.repo.GetSecret(ctx, secretID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entities.ErrorSecretNotFound(secretID)
	}

	if err != nil {
		return nil, err
	}

	// Decrypt secret
	// decrypted, err := crypto.Decrypt(password, secret.Payload)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to decrypt secret: %w", err)
	// }

	// Unmarshal corresponding struct
	// switch secret.SType {
	// case string(models.CredSecret):
	// 	secret.Creds = &models.Credentials{}
	// 	if err := json.Unmarshal(decrypted, secret.Creds); err != nil {
	// 		return nil, fmt.Errorf("failed to extract credentials: %w", err)
	// 	}
	// case string(models.TextSecret):
	// 	secret.Text = &models.Text{}
	// 	if err := json.Unmarshal(decrypted, secret.Text); err != nil {
	// 		return nil, fmt.Errorf("failed to extract text: %w", err)
	// 	}
	// case string(models.CardSecret):
	// 	secret.Card = &models.Card{}
	// 	if err := json.Unmarshal(decrypted, secret.Card); err != nil {
	// 		return nil, fmt.Errorf("failed to extract card: %w", err)
	// 	}
	// case string(models.BlobSecret):
	// 	secret.Blob = &models.Blob{}
	// 	if err := json.Unmarshal(decrypted, secret.Blob); err != nil {
	// 		return nil, fmt.Errorf("failed to extract blob: %w", err)
	// 	}
	// }

	return secret, nil
}

// func validateType(sType string) bool {
// 	allowedTypes := []string{
// 		string(pkgModel.CredSecret),
// 		string(pkgModel.TextSecret),
// 		string(pkgModel.BlobSecret),
// 		string(pkgModel.CardSecret),
// 	}

// 	return slices.Contains(allowedTypes, sType)
// }

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
	var (
		err                error
		payload, encrypted []byte
	)

	// if ok := validateType(secret.SType); !ok {
	// 	return 0, model.ErrWrongSecretType
	// }

	payload = []byte("dump")
	encrypted = payload // TODO

	// marshal corresponding data
	// switch secret.SType {
	// case string(pkgModel.CredSecret):
	// 	payload, err = json.Marshal(secret.Creds)
	// case string(pkgModel.TextSecret):
	// 	payload, err = json.Marshal(secret.Text)
	// case string(pkgModel.CardSecret):
	// 	payload, err = json.Marshal(secret.Card)
	// case string(pkgModel.BlobSecret):
	// 	payload, err = json.Marshal(secret.Blob)
	// }

	if err != nil {
		return nil, fmt.Errorf("failed to save secret data: %w", err)
	}

	// validate card number using Luhn's algo
	// if secret.SType == string(pkgModel.CardSecret) {
	// 	if ok := utils.CheckLuhn(secret.Card.Number); !ok {
	// 		return 0, model.ErrNumberInvaliod
	// 	}
	// }

	// Encrypt secret data
	// encrypted, err = crypto.Encrypt(password, payload)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to encrypt data: %w", err)
	// }

	secret.Payload = encrypted
	secret.ID, err = s.repo.Create(ctx, secret)

	if err != nil {
		return nil, fmt.Errorf("failed to store secret: %w", err)
	}

	return secret, nil
}

// Try update secret
func (s SecretsService) UpdateSecret(ctx context.Context, secret *models.Secret) (*models.Secret, error) {
	var (
		err                error
		payload, encrypted []byte
	)

	// if ok := validateType(secret.SType); !ok {
	// 	return 0, model.ErrWrongSecretType
	// }

	payload = []byte("dump")
	encrypted = payload // TODO

	// marshal corresponding data
	// switch secret.SType {
	// case string(pkgModel.CredSecret):
	// 	payload, err = json.Marshal(secret.Creds)
	// case string(pkgModel.TextSecret):
	// 	payload, err = json.Marshal(secret.Text)
	// case string(pkgModel.CardSecret):
	// 	payload, err = json.Marshal(secret.Card)
	// case string(pkgModel.BlobSecret):
	// 	payload, err = json.Marshal(secret.Blob)
	// }

	if err != nil {
		return nil, fmt.Errorf("failed to save secret data: %w", err)
	}

	// validate card number using Luhn's algo
	// if secret.SType == string(pkgModel.CardSecret) {
	// 	if ok := utils.CheckLuhn(secret.Card.Number); !ok {
	// 		return 0, model.ErrNumberInvaliod
	// 	}
	// }

	// Encrypt secret data
	// encrypted, err = crypto.Encrypt(password, payload)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to encrypt data: %w", err)
	// }

	secret.Payload = encrypted

	err = s.repo.Update(ctx, secret)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entities.ErrorSecretNotFound(secret.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to store secret: %w", err)
	}

	return secret, nil
}
