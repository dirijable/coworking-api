package server

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type TimeoutConfig struct {
	ReadTimeout     time.Duration `yaml:"read"     env-default:"5s"`
	WriteTimeout    time.Duration `yaml:"write"    env-default:"5s"`
	IdleTimeout     time.Duration `yaml:"idle"     env-default:"30s"`
	ShutdownTimeout time.Duration `yaml:"shutdown" env-default:"10s"`
}

type HTTPServerConfig struct {
	TimeoutConfig `yaml:"timeout"`
	Addr          string `yaml:"address"  env-default:"0.0.0.0:8080"`
}

type Config struct {
	HTTPServerConfig `yaml:"server"`
}

func MustNewConfig(path string) Config {
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		err = fmt.Errorf("read conifg: %w", err)
		panic(err)
	}
	return config
}
