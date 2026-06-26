package pgpool

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	User     string `envconfig:"USER" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	Host     string `envconfig:"HOST" required:"true"`
	Port     string `envconfig:"PORT" default:"5433"`
	DBName   string `envconfig:"DB_NAME" required:"true"`
	SSLMode  string `envconfig:"SSL_MODE" default:"disable"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("PG", &config); err != nil {
		return Config{}, fmt.Errorf("parse db config error: %w", err)
	}
	return config, nil
}

func MustNewConfig() Config {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return config
}
