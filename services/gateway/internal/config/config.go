package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AuthAPI         AuthApiConfig
	CurrencyService CurrencyServiceConfig
	Net             NetConfig
}

type config struct {
	AuthAPI         authApiConfig
	CurrencyService currencyServiceConfig
	Net             netConfig
}

func New() (Config, error) {
	cfg := config{}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("read env: %w", err)
	}

	fmt.Printf("Config: %+v\n", cfg)

	return Config{
		AuthAPI:         cfg.AuthAPI,
		CurrencyService: cfg.CurrencyService,
		Net:             cfg.Net,
	}, nil
}
