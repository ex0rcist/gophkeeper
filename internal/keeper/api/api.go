// Provides interface for API communication
package api

import (
	"context"
	"gophkeeper/pkg/models"

	tea "github.com/charmbracelet/bubbletea"
)

//go:generate mockgen -source client.go -destination mocks/mock_client.go -package api
type IApiClient interface {
	Register(ctx context.Context, login string, password string) (string, error)
	Login(ctx context.Context, login string, password string) (string, error)

	LoadSecrets(ctx context.Context) ([]*models.Secret, error)
	LoadSecret(ctx context.Context, ID uint64) (*models.Secret, error)
	SaveSecret(ctx context.Context, secret *models.Secret) error
	DeleteSecret(ctx context.Context, ID uint64) error

	SetToken(token string)
	GetToken() string

	SetPassword(password string)
	GetPassword() string

	Notifications(p *tea.Program)
}
