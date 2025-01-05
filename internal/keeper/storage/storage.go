package storage

import (
	"context"
	"gophkeeper/pkg/models"
)

// // Kinds of storage
// const (
// 	KindMemory   = "memory"
// 	KindFile     = "file"
// 	KindDatabase = "database"
// )

type Storage interface {
	Get(ctx context.Context, id uint64) (*models.Secret, error)
	GetAll(ctx context.Context) ([]*models.Secret, error)
	Create(ctx context.Context, secret *models.Secret) error
	Update(ctx context.Context, secret *models.Secret) error
	Delete(ctx context.Context, id uint64) error
	String() string
	Close(ctx context.Context) error
}
