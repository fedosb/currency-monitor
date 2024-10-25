package config

type NetConfig interface {
	GetHTTPAddress() string
}

type netConfig struct {
	HTTPAddress string `env-required:"true" env:"HTTP_ADDRESS"`
}

func (c netConfig) GetHTTPAddress() string {
	return c.HTTPAddress
}
