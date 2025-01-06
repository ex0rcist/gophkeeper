// Server config
package config

import (
	"fmt"
	"gophkeeper/internal/server/entities"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/dig"
)

// Config is main configuration of client application.
type Config struct {
	Address     string
	PostgresDSN entities.SecretConnURI
	LogLevel    string
	SecretKey   string // key to sign jwt
	EnableTLS   bool
}

// Shortcut to use with dig
type Dependency struct {
	dig.In
	Config *Config
}

// Create server config from ENV vars and cmd flags
func New() *Config {
	viper.SetDefault("address", "127.0.0.1:50051")
	viper.SetDefault("postgres-dsn", "postgres://cm:cm@localhost:5432/gophkeeper?sslmode=disable")

	viper.SetDefault("verbose", false)
	viper.SetDefault("log-level", "INFO")
	viper.SetDefault("secret-key", "123456") // TODO: remove default, add warning

	viper.SetEnvPrefix("GOPH")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	cfg := &Config{
		Address:     viper.GetString("address"),
		PostgresDSN: entities.SecretConnURI(viper.GetString("postgres-dsn")),
		LogLevel:    viper.GetString("log-level"),
		EnableTLS:   true,
	}

	return cfg
}

func (c *Config) String() string {
	var sb strings.Builder

	sb.WriteString("Configuration:\n")
	sb.WriteString(fmt.Sprintf("\t\tServer address: %s\n", c.Address))
	sb.WriteString(fmt.Sprintf("\t\tPostgres DSN: %s\n", c.Address))

	return sb.String()
}
