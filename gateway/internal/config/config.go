package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Name           string `mapstructure:"name"`
	MigrationsPath string `mapstructure:"migrations_path"`
}

func (dc DatabaseConfig) ToDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dc.User, dc.Password, dc.Host, dc.Port, dc.Name,
	)
}

type Config struct {
	Server         ServerConfig   `mapstructure:"server"`
	Auth           AuthConfig     `mapstructure:"auth"`
	GRPC           GRPCConfig     `mapstructure:"grpc"`
	DatabaseConfig DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type AuthConfig struct {
	BaseURL string `mapstructure:"base_url"`
}

type GRPCConfig struct {
	CurrencyServiceURL string `mapstructure:"currency_service_url"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return cfg, nil
}
