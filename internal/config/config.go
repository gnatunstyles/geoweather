package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ApiKey string `env:"API_GEO_KEY"`
}

func New() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return &Config{}, fmt.Errorf("error during config:%s", err)
	}
	log.Info().Msgf("config: %v", cfg)

	return cfg, nil
}
