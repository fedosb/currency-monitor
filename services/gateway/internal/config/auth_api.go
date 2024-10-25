package config

type AuthApiConfig interface {
	GetURL() string
}

type authApiConfig struct {
	URL string `env-required:"true" env:"AUTH_API_URL"`
}

func (c authApiConfig) GetURL() string {
	return c.URL
}
