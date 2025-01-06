// Keeper config
package config

import (
	"gophkeeper/internal/keeper/entities"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress entities.Address
	Verbose       bool
	EnableTLS     bool
	LogLevel      string
	BuildDate     string
	BuildVersion  string
}

func New() *Config {
	viper.SetDefault("address", "127.0.0.1:50051")
	viper.SetDefault("verbose", false)

	viper.SetEnvPrefix("GOPH")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	cfg := &Config{
		ServerAddress: entities.Address(viper.GetString("address")),
		Verbose:       viper.GetBool("verbose"),
		EnableTLS:     true,
	}

	return cfg
}
