package config

type CurrencyApiConfig interface {
	GetURL() string
}

type currencyApiConfig struct {
	URL string `env-required:"true" env:"CURRENCY_API_URL"`
}

func (c currencyApiConfig) GetURL() string {
	return c.URL
}
