package main

import (
	"fmt"
	"github.com/fedosb/currency-monitor/services/currency/internal/db/postgres"

	"github.com/fedosb/currency-monitor/services/currency/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	db, err := postgres.New(cfg.DB)
	if err != nil {
		panic(err)
	}

	if cfg.DB.Migrate() {
		fmt.Println("Applying migrations...")
		if err := db.ApplyMigrations(); err != nil {
			panic(err)
		}
	}
}
