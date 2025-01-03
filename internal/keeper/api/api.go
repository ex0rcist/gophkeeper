// Provides interface for API communication
package api

import (
	"context"
	"gophkeeper/pkg/models"
)

//go:generate mockgen -source client.go -destination mocks/mock_client.go -package api
type IApiClient interface {
	Register(ctx context.Context, login string, password string) (string, error)
	Login(ctx context.Context, login string, password string) (string, error)

	LoadSecret(ctx context.Context, ID uint64) (*models.Secret, error)

	SaveCredential(ctx context.Context, ID uint64, metadata, login, password string) error
	SaveText(ctx context.Context, ID uint64, metadata, text string) error
	SaveCard(ctx context.Context, ID uint64, metadata, number string, expMonth, expYear, cvv uint32) error

	UploadFile(ctx context.Context, metadata string, filePath string) error
	DownloadFile(ctx context.Context, ID uint64, fileName string) error
}
