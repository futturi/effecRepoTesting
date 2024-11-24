package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port               string `env:"PORT"`
	DbConnectionString string `env:"DB_CONNECTION_STRING"`
	ApiUrl             string `env:"API_URL"`
	ApiTimeout         int    `env:"API_TIMEOUT"`
	LogLevel           string `env:"LOG_LEVEL"`
}

func NewConfig() (Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig("config/.env", &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
