package keeper

import (
	"context"
	"errors"
	"gophkeeper/internal/keeper/api"
	"gophkeeper/internal/keeper/config"
	"gophkeeper/internal/keeper/tui/app"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

var (
	errNotInitialized = errors.New("application is not initialized")
)

const shutdownTimeout = 5 * time.Second

type Keeper struct {
	config *config.Config
	tuiApp *app.App
	log    *zap.SugaredLogger

	client api.IApiClient
}

type KeeperDependencies struct {
	dig.In

	App    *app.App
	Config *config.Config
	Client api.IApiClient
	Logger *zap.SugaredLogger
}

func NewKeeper(deps KeeperDependencies) (*Keeper, error) {
	return &Keeper{
		tuiApp: deps.App,
		config: deps.Config,
		client: deps.Client,
		log:    deps.Logger,
	}, nil
}

func (k Keeper) Start() error {
	k.tuiApp.Start()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case sig := <-quit:
		k.log.Info("interrupt: signal " + sig.String())
	case err := <-k.tuiApp.Notify():
		if errors.Is(err, app.ErrExitCmd) {
			quit <- syscall.SIGINT
			break
		}

		k.log.Error(err, "TUI app error")
	}

	k.shutdown()

	return nil
}

// Shutdown gracefully stops client application.
func (k Keeper) shutdown() {
	stopped := make(chan struct{})
	stopCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	go func() {
		// Shut down seervices here, if any
		// TODO: notification service

		close(stopped)
	}()

	select {
	case <-stopped:
		k.log.Info("keeper shutdown successful")

	case <-stopCtx.Done():
		k.log.Info("shutdown timeout exceeded")
	}
}
