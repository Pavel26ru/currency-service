package config

import "fmt"

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
