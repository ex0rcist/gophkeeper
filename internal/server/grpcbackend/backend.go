package grpcbackend

import (
	"gophkeeper/internal/server/config"
	grpchandlers "gophkeeper/internal/server/grpcbackend/handlers"
	"gophkeeper/internal/server/grpcbackend/interceptor"
	"gophkeeper/pkg/proto/keeper/grpcapi"
	"net"

	"go.uber.org/dig"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
)

type Backend struct {
	// privateKey    security.PrivateKey
	trustedSubnet *net.IPNet

	server *grpc.Server
}

type BackendDependencies struct {
	dig.In

	Logger        *zap.SugaredLogger
	Config        *config.Config
	HealthServer  *grpchandlers.HealthServer
	UsersServer   *grpchandlers.UsersServer
	SecretsServer *grpchandlers.SecretsServer
	// NotificationServer    Notificationerver
}

// Backend constructor
func NewBackend(deps BackendDependencies) *Backend {
	iceps := make([]grpc.UnaryServerInterceptor, 0, 2)
	iceps = append(iceps, interceptor.Authentication([]byte(deps.Config.SecretKey)))
	iceps = append(iceps, interceptor.Logger(deps.Logger))

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(iceps...),
	}

	grpcServer := grpc.NewServer(grpcOpts...)

	// Register servers
	grpcapi.RegisterHealthServer(grpcServer, deps.HealthServer)
	grpcapi.RegisterUsersServer(grpcServer, deps.UsersServer)
	grpcapi.RegisterSecretsServer(grpcServer, deps.SecretsServer)

	backend := &Backend{server: grpcServer}

	return backend
}
