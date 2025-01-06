package grpcbackend

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gophkeeper/cert"
	"gophkeeper/internal/server/config"
	grpchandlers "gophkeeper/internal/server/grpcbackend/handlers"
	"gophkeeper/internal/server/grpcbackend/interceptor"
	"gophkeeper/pkg/proto/keeper/grpcapi"

	"go.uber.org/dig"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
)

type Backend struct {
	server *grpc.Server
}

type BackendDependencies struct {
	dig.In

	Logger             *zap.SugaredLogger
	Config             *config.Config
	HealthServer       *grpchandlers.HealthServer
	UsersServer        *grpchandlers.UsersServer
	SecretsServer      *grpchandlers.SecretsServer
	NotificationServer *grpchandlers.NotificationServer
}

// Backend constructor
func NewBackend(deps BackendDependencies) (*Backend, error) {
	iceps := make([]grpc.UnaryServerInterceptor, 0, 2)
	iceps = append(iceps, interceptor.Authentication([]byte(deps.Config.SecretKey)))
	iceps = append(iceps, interceptor.Logger(deps.Logger))

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(iceps...),
	}

	// TLS config
	if deps.Config.EnableTLS {
		tlsCreds, err := loadTLSConfig("ca-cert.pem", "server-cert.pem", "server-key.pem")
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS config: %w", err)
		}

		// Append TLS credentials to server options
		grpcOpts = append(grpcOpts, grpc.Creds(tlsCreds))
	}

	// Stream interceptor
	grpcOpts = append(
		grpcOpts,
		grpc.StreamInterceptor(
			interceptor.StreamAuthentication([]byte(deps.Config.SecretKey)),
		),
	)

	grpcServer := grpc.NewServer(grpcOpts...)

	// Register servers
	grpcapi.RegisterHealthServer(grpcServer, deps.HealthServer)
	grpcapi.RegisterUsersServer(grpcServer, deps.UsersServer)
	grpcapi.RegisterSecretsServer(grpcServer, deps.SecretsServer)
	grpcapi.RegisterNotificationServer(grpcServer, deps.NotificationServer)

	backend := &Backend{server: grpcServer}

	return backend, nil
}

func loadTLSConfig(caCertFile, serverCertFile, serverKeyFile string) (credentials.TransportCredentials, error) {
	// Read CA cert
	caPem, err := cert.Cert.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	// Read server cert
	serverCertPEM, err := cert.Cert.ReadFile(serverCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read server cert: %w", err)
	}

	// Read server key
	serverKeyPEM, err := cert.Cert.ReadFile(serverKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read server key: %w", err)
	}

	// Create key pair
	serverCert, err := tls.X509KeyPair(serverCertPEM, serverKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to load x509 key pair: %w", err)
	}

	// Create cert pool and append CA's cert
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		return nil, fmt.Errorf("failed to append CA cert to cert pool: %w", err)
	}

	// Create config
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}
