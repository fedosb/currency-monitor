package postgres

import (
	"fmt"

	"github.com/fedosb/currency-monitor/services/currency/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	postgresDriver = "postgres"
)

type DB struct {
	*sqlx.DB
}

func New(cfg config.DBConfig) (*DB, error) {

	db, err := sqlx.Connect(postgresDriver, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns(100))
	db.SetMaxIdleConns(cfg.MaxIdleConns(10))
	db.SetConnMaxLifetime(cfg.MaxConnLifetime(0))

	return &DB{db}, nil
}
