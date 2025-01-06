package grpchandlers

import (
	"context"
	"errors"
	"gophkeeper/internal/server/entities"
	"gophkeeper/internal/server/service"
	"gophkeeper/pkg/proto/keeper/grpcapi"

	"go.uber.org/dig"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// HealthServer verifies current health status of the service.
type HealthServer struct {
	grpcapi.UnimplementedHealthServer

	healthManager service.HealthManager
}

type HealthServerDependencies struct {
	dig.In

	HealthManager service.HealthManager
}

func NewHealthServer(deps HealthServerDependencies) *HealthServer {
	return &HealthServer{
		healthManager: deps.HealthManager,
	}
}

// Ping verifies connection to the database.
func (s HealthServer) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	err := s.healthManager.Ping(ctx)
	if err == nil {
		return new(emptypb.Empty), nil
	}

	if errors.Is(err, entities.ErrStorageUnpingable) {
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	return nil, status.Error(codes.Internal, err.Error())
}
