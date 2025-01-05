package server

import (
	"context"
	"fmt"

	"gophkeeper/internal/server/config"
	"gophkeeper/internal/server/grpcbackend"

	"gophkeeper/internal/server/storage"

	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

const shutdownTimeout = 60 * time.Second

// Main server app
type Server struct {
	config  *config.Config
	log     *zap.SugaredLogger
	deps    *dig.Container
	storage storage.ServerStorage
	// privateKey     security.PrivateKey

	grpcServer *grpcbackend.GRPCServer
}

type ServerDependencies struct {
	dig.In

	Config     *config.Config
	Storage    storage.ServerStorage
	GRPCServer *grpcbackend.GRPCServer
	Logger     *zap.SugaredLogger
}

// Create new Server
func New(deps ServerDependencies) *Server {
	server := &Server{
		config:  deps.Config,
		log:     deps.Logger,
		storage: deps.Storage,

		grpcServer: deps.GRPCServer,
	}

	return server
}

// Start all subservices
func (s *Server) Start() error {

	// privateKey, err := preparePrivateKey(config)
	// if err != nil {
	// 	return nil, err
	// }

	s.grpcServer.Start()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	s.log.Info("server ready")

	select {
	case sig := <-quit:
		s.log.Info("interrupt: signal " + sig.String())
	case err := <-s.grpcServer.Notify():
		s.log.Error(err, "Server -> Start() -> s.grpcServer.Notify")
	}

	s.shutdown()

	return nil
}

// Stringer for logging
func (s *Server) String() string {
	var sb strings.Builder

	sb.WriteString("Configuration:\n")
	sb.WriteString(fmt.Sprintf("\t\tListen: %s\n", s.config.Address))

	sb.WriteString("Storage:\n")
	sb.WriteString(s.storage.String())

	return sb.String()
}

// Shutdown server and it subservices. Will block upto shutdownTimeout
func (s *Server) shutdown() {
	s.log.Info("shutting down application now")

	stopped := make(chan struct{})
	stopCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	go func() {
		s.log.Info("shutting down GRPC API...")
		if err := s.grpcServer.Shutdown(stopCtx); err != nil {
			s.log.Error(err)
		}

		close(stopped)
	}()

	select {
	case <-stopped:
		s.log.Info("server shutdown successful")

	case <-stopCtx.Done():
		s.log.Info("shutdown timeout exceeded")
	}
}

// func preparePrivateKey(config *Config) (security.PrivateKey, error) {
// 	var (
// 		privateKey security.PrivateKey
// 		err        error
// 	)

// 	if len(config.PrivateKeyPath) != 0 {
// 		privateKey, err = security.NewPrivateKey(config.PrivateKeyPath)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return privateKey, err
// }

// Common interface for different subservices: http, grpc, profiling, etc
type ServerService interface {
	Start()
	Notify() <-chan error
	Shutdown(ctx context.Context) error
	String() string
}

var _ ServerService = (*grpcbackend.GRPCServer)(nil)
