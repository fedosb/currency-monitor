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
		return errsinternal.NewValidationError(entity.RateValidationNameRequiredMsg)
	}

	if r.Date.IsZero() {
		return errsinternal.NewValidationError(entity.RateValidationDateRequiredMsg)
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
		return errsinternal.NewValidationError(entity.RateValidationNameRequiredMsg)
	}

	if r.From.IsZero() || r.To.IsZero() {
		return errsinternal.NewValidationError(entity.RateValidationFromAndToRequiredMsg)
	}

	if r.From.After(r.To) {
		return errsinternal.NewValidationError(entity.RateValidationFromBeforeToMsg)
	}

	return nil
}

type GetByNameAndDateRangeResponse struct {
	Rates []entity.Rate
}
