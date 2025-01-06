package storage

import (
	"context"
	"fmt"
)

// Kinds of storage
const (
	KindPostgres = "postgres"
)

// Common interface for storages
type ServerStorage interface {
	fmt.Stringer

	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}
