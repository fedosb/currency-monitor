package config

import (
	"fmt"
	"time"
)

type DBConfig interface {
	DSN() string
	MaxOpenConns(defaultValue int) int
	MaxIdleConns(defaultValue int) int
	MaxConnLifetime(defaultValue time.Duration) time.Duration
	Migrate() bool
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

func (c dbConfig) DSN() string {
	return fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s sslmode=%s`,
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Database,
		c.SSL,
	)
}

func (c dbConfig) MaxOpenConns(defaultValue int) int {
	if c.MaxOpenConnections == 0 {
		return defaultValue
	}

	return c.MaxIdleConnections
}

func (c dbConfig) MaxIdleConns(defaultValue int) int {
	if c.MaxIdleConnections == 0 {
		return defaultValue
	}

	return c.MaxIdleConnections
}

func (c dbConfig) MaxConnLifetime(defaultValue time.Duration) time.Duration {
	if c.MaxConnectionLifetime == 0 {
		return defaultValue
	}

	return c.MaxConnectionLifetime
}

func (c dbConfig) Migrate() bool {
	return c.ApplyMigrations
}
