package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	Address  string        `yaml:"address" env:"API_ADDRESS" env-default:"localhost:8080"`
	Timeout  time.Duration `yaml:"timeout" env:"API_TIMEOUT" env-default:"5s"`
}

func MustLoad(cfgPath string) Config {
	var cfg Config
	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
