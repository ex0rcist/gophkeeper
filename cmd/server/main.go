package main

import (
	"log"

	"gophkeeper/internal/server"
	"gophkeeper/internal/server/config"
	"gophkeeper/internal/server/httpbackend"
	httphandler "gophkeeper/internal/server/httpbackend/handler"
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
	container.Provide(func() *config.Config { return cfg })

	// Server
	container.Provide(server.New)

	// Logging
	container.Provide(utils.NewZapLogger)
	container.Provide(func(cfg *config.Config) utils.ZapLogLevel { return utils.ZapLogLevel(cfg.LogLevel) })

	// Routing
	container.Provide(chi.NewRouter, dig.As(new(chi.Router)))

	// HTTP server and backend
	container.Provide(httpbackend.NewHTTPServer)
	container.Provide(httpbackend.NewBackend)
	container.Provide(httpbackend.NewHTTPServerAddress)

	// HTTP handlers
	container.Provide(httphandler.NewHealthHandler)
	container.Provide(httphandler.NewSecretsHandler)
	container.Provide(httphandler.NewUsersHandler)

	// services
	container.Provide(service.NewHealthService, dig.As(new(service.HealthManager)))
	container.Provide(service.NewSecretsService, dig.As(new(service.SecretsManager)))
	container.Provide(service.NewUsersService, dig.As(new(service.SecretsManager)))

	return container
}

func addAppSpecificDependencies(container *dig.Container, cfg *config.Config) *dig.Container {
	switch {
	case len(cfg.PostgresDSN) > 0:
		// Postgres storage
		container.Provide(pgStorage.NewPostgresDSN)
		container.Provide(pgStorage.NewPostgresConn)
		container.Provide(pgStorage.NewPostgresStorage, dig.As(new(storage.ServerStorage)))

		// Postgres repos
		container.Provide(pgRepo.NewSecretsRepository, dig.As(new(repository.SecretsRepository)))
		container.Provide(pgRepo.NewUsersRepository, dig.As(new(repository.SecretsRepository)))
	}

	return container
}
