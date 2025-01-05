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

	client      api.IApiClient
	accessToken string // ???
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

	// app, err := client.NewKeeperClient(cfg)
	// if err != nil {
	// 	log.Fatal("failed to initialize app: ", err)
	// }

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Run notification monitor
	// go app.Client.Notifications(p)

	select {
	case sig := <-quit:
		k.log.Info("interrupt: signal " + sig.String())
	case err := <-k.tuiApp.Notify():
		if errors.Is(err, app.ExitCmdErr) {
			quit <- syscall.SIGINT
			break
		}

		k.log.Error(err, "TUI app error")
	}

	k.shutdown()

	return nil
}

// // WithContext injects App into provided context.
// func (a *App) WithContext(ctx context.Context) context.Context {
// 	return context.WithValue(ctx, appKeyName, a)
// }

// FromContext extracts App from provided context.
// func FromContext(ctx context.Context) (*App, error) {
// 	if val := ctx.Value(appKeyName); val != nil {
// 		return val.(*App), nil
// 	}

// 	return nil, errNotInitialized
// }

// Shutdown gracefully stops client application.
func (k Keeper) shutdown() {
	// if err := a.conn.Close(); err != nil {
	// 	a.Log.Warn().Err(err).Msg("app - Shutdown - a.conn.Close")
	// }

	stopped := make(chan struct{})
	stopCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	go func() {
		// shut down seervices here ...

		close(stopped)
	}()

	select {
	case <-stopped:
		k.log.Info("server shutdown successful")

	case <-stopCtx.Done():
		k.log.Info("shutdown timeout exceeded")
	}
}
