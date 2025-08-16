package config

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/vctrl/currency-service/pkg/config"
)

type ServiceConfig struct {
	ServerPort string `mapstructure:"server_port"`
}

type APIConfig struct {
	BaseURL        string `mapstructure:"base_url"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

type WorkerConfig struct {
	Schedule     string `mapstructure:"schedule"`
	CurrencyPair struct {
		BaseCurrency   string `mapstructure:"base_currency"`
		TargetCurrency string `mapstructure:"target_currency"`
	} `mapstructure:"currency_pair"`
}

type AppConfig struct {
	Service  ServiceConfig         `mapstructure:"service"`
	API      APIConfig             `mapstructure:"api"`
	Database config.DatabaseConfig `mapstructure:"database"`
	Worker   WorkerConfig          `mapstructure:"worker"`
}

func LoadConfig(path string) (AppConfig, error) {
	var config AppConfig

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return config, nil
}
