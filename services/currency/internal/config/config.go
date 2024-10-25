package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB          DBConfig
	Net         NetConfig
	CurrencyAPI CurrencyApiConfig
	Fetcher     FetcherConfig
}

type config struct {
	DB          dbConfig
	Net         netConfig
	CurrencyApi currencyApiConfig
	Fetcher     fetcherConfig
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
		Fetcher:     cfg.Fetcher,
	}, nil
}
