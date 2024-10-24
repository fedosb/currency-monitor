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

type GetByNameAndDateResponse struct {
	Rate entity.Rate
}

type GetByNameAndDateRangeRequest struct {
	Name string
	From time.Time
	To   time.Time
}

func (r GetByNameAndDateRangeRequest) Validate() error {
	if r.From.After(r.To) {
		return errsinternal.InvalidDateRangeError
	}

	return nil
}

type GetByNameAndDateRangeResponse struct {
	Rates []entity.Rate
}
