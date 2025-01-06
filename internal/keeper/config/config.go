package config

import (
	"fmt"
	"gophkeeper/internal/keeper/entities"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Username      string
	Password      entities.Secret
	CAPath        string
	Verbose       bool
	LogLevel      string
	ServerAddress entities.Address
	DownloadPath  string
	BuildDate     string
	BuildVersion  string
}

func New() *Config {
	viper.SetDefault("server-address", "127.0.0.1:50051")
	viper.SetDefault("verbose", false)

	viper.SetEnvPrefix("GOPH")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	cfg := &Config{
		Username:      viper.GetString("username"),
		Password:      entities.Secret(viper.GetString("password")),
		ServerAddress: entities.Address(viper.GetString("server-address")),
		CAPath:        viper.GetString("ca-path"),
		Verbose:       viper.GetBool("verbose"),
	}

	return cfg
}

func (c *Config) String() string {
	var sb strings.Builder

	sb.WriteString("Configuration:\n")
	sb.WriteString(fmt.Sprintf("\t\tUsername: %s\n", c.Username))
	sb.WriteString(fmt.Sprintf("\t\tPassword: %s\n", c.Password))
	sb.WriteString(fmt.Sprintf("\t\tKeeper address: %s\n", c.ServerAddress))
	sb.WriteString(fmt.Sprintf("\t\tCertificate authority path: %s\n", c.CAPath))
	sb.WriteString(fmt.Sprintf("\t\tVerbose: %t", c.Verbose))

	return sb.String()
}
