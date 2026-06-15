package logger

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type lConfig struct {
	Level   string   `yaml:"level" env-default:"info"`
	Format  string   `yaml:"format" env-default:"text"`
	Folder  string   `yaml:"folder" env-default:"output/logger"`
	Outputs []string `yaml:"outputs" env-default:"console"`
}
type Config struct {
	lConfig `yaml:"logger"`
}

func MustNewConfig(configPath string) Config {
	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		err = fmt.Errorf("read logger config: %w", err)
		panic(err)
	}
	return config
}
