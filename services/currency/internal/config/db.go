package config

import (
	"fmt"
	"time"
)

type DBConfig interface {
	GetDSN() string
	GetMaxOpenConnections(defaultValue int) int
	GetMaxIdleConnections(defaultValue int) int
	GetMaxConnectionLifetime(defaultValue time.Duration) time.Duration
	GetApplyMigrations() bool
}

type dbConfig struct {
	Host                  string        `env-required:"true" env:"DB_HOST"`
	Port                  int           `env-required:"true" env:"DB_PORT"`
	User                  string        `env-required:"true" env:"DB_USERNAME"`
	Password              string        `env-required:"true" env:"DB_PASSWORD"`
	Database              string        `env-required:"true" env:"DB_DATABASE"`
	SSL                   string        `env-required:"true" env:"DB_SSL"`
	MaxOpenConnections    int           `env:"DB_MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections    int           `env:"DB_MAX_IDLE_CONNECTIONS"`
	MaxConnectionLifetime time.Duration `env:"DB_MAX_CONNECTION_LIFETIME"`
	ApplyMigrations       bool          `env:"DB_APPLY_MIGRATIONS"`
}

func (c dbConfig) GetDSN() string {
	return fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s sslmode=%s`,
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Database,
		c.SSL,
	)
}

func (c dbConfig) GetMaxOpenConnections(defaultValue int) int {
	if c.MaxOpenConnections == 0 {
		return defaultValue
	}

	return c.MaxIdleConnections
}

func (c dbConfig) GetMaxIdleConnections(defaultValue int) int {
	if c.MaxIdleConnections == 0 {
		return defaultValue
	}

	return c.MaxIdleConnections
}

func (c dbConfig) GetMaxConnectionLifetime(defaultValue time.Duration) time.Duration {
	if c.MaxConnectionLifetime == 0 {
		return defaultValue
	}

	return c.MaxConnectionLifetime
}

func (c dbConfig) GetApplyMigrations() bool {
	return c.ApplyMigrations
}
