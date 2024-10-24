package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/fedosb/currency-monitor/services/currency/internal/db/postgres"
	"github.com/fedosb/currency-monitor/services/currency/internal/entity"
	errsinternal "github.com/fedosb/currency-monitor/services/currency/internal/errors"
)

type RateRepository struct {
	db *postgres.DB
}

func NewRateRepository(db *postgres.DB) *RateRepository {
	return &RateRepository{db: db}
}

type rate struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Name      string    `db:"name"`
	Date      time.Time `db:"date"`
	Rate      float64   `db:"rate"`
}

func (r rate) Entity() entity.Rate {
	return entity.Rate{
		ID:        r.ID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Name:      r.Name,
		Date:      r.Date,
		Rate:      r.Rate,
	}
}

func mapRate(r entity.Rate) rate {
	return rate{
		ID:        r.ID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Name:      r.Name,
		Date:      r.Date,
		Rate:      r.Rate,
	}
}

type rateList []rate

func (l rateList) Entities() []entity.Rate {
	rates := make([]entity.Rate, 0, len(l))

	for r := range l {
		rates = append(rates, l[r].Entity())
	}

	return rates
}

func (r *RateRepository) Save(ctx context.Context, rate entity.Rate) (entity.Rate, error) {
	rateModel := mapRate(rate)

	query, args, err := sqlx.Named(
		`INSERT INTO rates (name, date, rate)
				VALUES (:name, :date, :rate)
				ON CONFLICT (name, date)
    			DO UPDATE SET	
    				rate       = EXCLUDED.rate,
                  	updated_at = now()
				RETURNING *`,
		&rateModel,
	)
	if err != nil {
		return entity.Rate{}, fmt.Errorf("mapping query: %w", err)
	}

	query = r.db.Rebind(query)

	if err = r.db.QueryRowxContext(ctx, query, args...).StructScan(&rateModel); err != nil {
		return entity.Rate{}, fmt.Errorf("query execution: %w", err)
	}

	return rateModel.Entity(), nil
}

func (r *RateRepository) GetByNameAndDate(ctx context.Context, name string, date time.Time) (entity.Rate, error) {
	var rateModel rate

	err := r.db.GetContext(ctx,
		&rateModel,
		`SELECT * FROM rates WHERE name = $1 AND date = $2`,
		name, date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Rate{}, errsinternal.NotFoundError
		}

		return entity.Rate{}, fmt.Errorf("select rate: %w", err)
	}

	return rateModel.Entity(), nil
}

func (r *RateRepository) GetByNameAndDateRange(ctx context.Context, name string, from, to time.Time) ([]entity.Rate, error) {
	var rateModels rateList

	err := r.db.SelectContext(ctx,
		&rateModels,
		`SELECT * FROM rates WHERE name = $1 AND date BETWEEN $2 AND $3 ORDER BY id DESC`,
		name, from, to)
	if err != nil {
		return nil, fmt.Errorf("select rates: %w", err)
	}

	return rateModels.Entities(), nil
}
