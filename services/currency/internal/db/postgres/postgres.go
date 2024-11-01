package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/fedosb/currency-monitor/services/currency/internal/config"
)

const (
	postgresDriver = "postgres"
	migrationsPath = "migrations"
)

type DB struct {
	*sqlx.DB
}

func New(cfg config.DBConfig) (*DB, error) {

	db, err := sqlx.Connect(postgresDriver, cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}

	db.SetMaxOpenConns(cfg.GetMaxOpenConnections(100))
	db.SetMaxIdleConns(cfg.GetMaxIdleConnections(10))
	db.SetConnMaxLifetime(cfg.GetMaxConnectionLifetime(0))

	return &DB{db}, nil
}

func (db *DB) ApplyMigrations() error {
	goose.SetLogger(gooseLogger{})

	if err := goose.SetDialect(postgresDriver); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.Up(db.DB.DB, migrationsPath); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}

func (db *DB) RollbackMigrations() error {
	goose.SetLogger(gooseLogger{})

	if err := goose.SetDialect(postgresDriver); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.Down(db.DB.DB, migrationsPath); err != nil {
		return fmt.Errorf("rollback migrations: %w", err)
	}

	return nil
}
