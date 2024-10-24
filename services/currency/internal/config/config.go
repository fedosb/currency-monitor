package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB          DBConfig
	Net         NetConfig
	CurrencyAPI CurrencyApiConfig
}

type config struct {
	DB          dbConfig
	Net         netConfig
	CurrencyApi currencyApiConfig
}

func New() (Config, error) {
	cfg := config{}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("read env: %w", err)
	}

	return Config{
		DB:          cfg.DB,
		Net:         cfg.Net,
		CurrencyAPI: cfg.CurrencyApi,
	}, nil
}
