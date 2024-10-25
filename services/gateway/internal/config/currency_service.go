package config

type CurrencyServiceConfig interface {
	GetAddress() string
}

type currencyServiceConfig struct {
	Address string `env-required:"true" env:"CURRENCY_SERVICE_ADDRESS"`
}

func (c currencyServiceConfig) GetAddress() string {
	return c.Address
}
