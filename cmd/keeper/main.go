package main

import (
	"gophkeeper/internal/keeper"
	"gophkeeper/internal/keeper/config"
	"gophkeeper/internal/keeper/tui/top"
	"gophkeeper/internal/keeper/utils"
	"log"

	"go.uber.org/dig"
	//_ "github.com/charmbracelet/lipgloss"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {

	cfg := config.New()

	err := runApp(cfg)
	if err != nil {
		log.Fatal("failed to start app: ", err)
	}
}

func runApp(cfg *config.Config) error {

	top.Start(cfg)

	return nil
	// cont := buildDepContainer(cfg)

	// return cont.Invoke(func(keeper *keeper.Keeper) error {
	// 	return keeper.Start()
	// })
}

func buildDepContainer(cfg *config.Config) *dig.Container {
	container := dig.New()

	// NB: could use .Provide(config.New), but try to avoid
	// parsing config twice, using closure instead
	container.Provide(func() *config.Config { return cfg })

	// Server
	container.Provide(keeper.New)

	// Logging
	container.Provide(utils.NewZapLogger)
	container.Provide(func(cfg *config.Config) utils.ZapLogLevel { return utils.ZapLogLevel(cfg.LogLevel) })

	return container
}
