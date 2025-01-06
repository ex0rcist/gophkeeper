package service

import (
	"context"
	"fmt"
	"gophkeeper/internal/server/entities"
	"gophkeeper/internal/server/storage"
	"time"

	"go.uber.org/dig"
)

const defaultTimeout = 5 * time.Second

var _ HealthManager = HealthService{}

type HealthManager interface {
	Ping(ctx context.Context) error
}

type HealthManagerDependencies struct {
	dig.In
	Storage storage.ServerStorage
}

type HealthService struct {
	storage storage.ServerStorage
}

// Service constructor
func NewHealthService(deps HealthManagerDependencies) *HealthService {
	return &HealthService{storage: deps.Storage}
}

// Interface to check if storage supports healthcheck
type PingableStorage interface {
	Ping(ctx context.Context) error
}

// Ping-pong
func (s HealthService) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	strg, ok := s.storage.(PingableStorage)
	if !ok {
		return fmt.Errorf("storage ping failed: %w", entities.ErrStorageUnpingable)
	}

	if err := strg.Ping(ctx); err != nil {
		return fmt.Errorf("storage check failed: %w", err)
	}

	return nil
}
