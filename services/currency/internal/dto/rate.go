package dto

import (
	"time"

	"github.com/fedosb/currency-monitor/services/currency/internal/entity"
	errsinternal "github.com/fedosb/currency-monitor/services/currency/internal/errors"
)

type GetByNameAndDateRequest struct {
	Name string
	Date time.Time
}

func (r GetByNameAndDateRequest) Validate() error {
	if r.Name == "" {
		return errsinternal.NewValidationError("name is required")
	}

	if r.Date.IsZero() {
		return errsinternal.NewValidationError("date is required")
	}

	return nil
}

type GetByNameAndDateResponse struct {
	Rate entity.Rate
}

type GetByNameAndDateRangeRequest struct {
	Name string
	From time.Time
	To   time.Time
}

func (r GetByNameAndDateRangeRequest) Validate() error {
	if r.Name == "" {
		return errsinternal.NewValidationError("name is required")
	}

	if r.From.IsZero() || r.To.IsZero() {
		return errsinternal.NewValidationError("from and to are required")
	}

	if r.From.After(r.To) {
		return errsinternal.NewValidationError("from must be before to")
	}

	return nil
}

type GetByNameAndDateRangeResponse struct {
	Rates []entity.Rate
}
