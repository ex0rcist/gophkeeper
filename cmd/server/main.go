package main

import (
	"fmt"
	"log"
	"os"

	"gophkeeper/internal/server"
	"gophkeeper/internal/server/config"
	"gophkeeper/internal/server/grpcbackend"
	grpchandlers "gophkeeper/internal/server/grpcbackend/handlers"
	"gophkeeper/internal/server/repository"
	"gophkeeper/internal/server/service"
	"gophkeeper/internal/server/storage"
	"gophkeeper/internal/server/utils"

	pgRepo "gophkeeper/internal/server/repository/postgres"
	pgStorage "gophkeeper/internal/server/storage/postgres"

	"github.com/go-chi/chi/v5"
	"go.uber.org/dig"
)

func main() {
	cfg := config.New()

	if cfg.PostgresDSN == "" {
		fmt.Println("please provide GOPH_POSTGRES_DSN in ENV")
		os.Exit(1)
	}

	err := runApp(cfg)
	if err != nil {
		log.Fatal("failed to start app: ", err)
	}
}

// Single entry point for all app's dependencies
func runApp(cfg *config.Config) error {
	cont := buildDepContainer(cfg)
	cont = addAppSpecificDependencies(cont, cfg)

	return cont.Invoke(func(server *server.Server) error {
		return server.Start()
	})
}

func buildDepContainer(cfg *config.Config) *dig.Container {
	container := dig.New()

	// NB: could use .Provide(config.New), but try to avoid
	// parsing config twice, using closure instead
	_ = container.Provide(func() *config.Config { return cfg })

	// Server
	_ = container.Provide(server.New)

	// Logging
	_ = container.Provide(utils.NewZapLogger)
	_ = container.Provide(func(cfg *config.Config) utils.ZapLogLevel { return utils.ZapLogLevel(cfg.LogLevel) })

	// Routing
	_ = container.Provide(chi.NewRouter, dig.As(new(chi.Router)))

	// HTTP server and backend
	_ = container.Provide(grpcbackend.NewGRPCServer)
	_ = container.Provide(grpcbackend.NewBackend)
	_ = container.Provide(grpcbackend.NewGRPCServerAddress)

	// HTTP handlers
	_ = container.Provide(grpchandlers.NewUsersServer)
	_ = container.Provide(grpchandlers.NewHealthServer)
	_ = container.Provide(grpchandlers.NewSecretsServer)
	// TODO: notification server

	// services
	_ = container.Provide(service.NewHealthService, dig.As(new(service.HealthManager)))
	_ = container.Provide(service.NewSecretsService, dig.As(new(service.SecretsManager)))
	_ = container.Provide(service.NewUsersService, dig.As(new(service.UsersManager)))

	return container
}

func addAppSpecificDependencies(container *dig.Container, cfg *config.Config) *dig.Container {
	switch {
	case len(cfg.PostgresDSN) > 0:
		// Postgres storage
		_ = container.Provide(pgStorage.NewPostgresDSN)
		_ = container.Provide(pgStorage.NewPostgresConn)
		_ = container.Provide(pgStorage.NewPostgresStorage, dig.As(new(storage.ServerStorage)))

		// Postgres repos
		_ = container.Provide(pgRepo.NewUsersRepository, dig.As(new(repository.UsersRepository)))
		_ = container.Provide(pgRepo.NewSecretsRepository, dig.As(new(repository.SecretsRepository)))
	}

	return container
}
