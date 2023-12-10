package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		WbApi WbApiConfig
	}

	WbApiConfig struct {
		Token  string `env:"WB_API_TOKEN,required"`
		ApiUrl string `env:"WB_API_URL,required"`
	}
)

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
