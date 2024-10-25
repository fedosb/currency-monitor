package dto

import (
	"errors"
	"time"
)

type Rate struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Date      time.Time `json:"date"`
	Rate      float64   `json:"rate"`
}

type GetByNameAndDateRequest struct {
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

func (r GetByNameAndDateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("empty name")
	}

	if r.Date.IsZero() {
		return errors.New("empty date")
	}

	return nil
}

type GetByNameAndDateResponse struct {
	Rate Rate `json:"rate"`
}

type GetByNameAndDateRangeRequest struct {
	Name string    `json:"name"`
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

func (r GetByNameAndDateRangeRequest) Validate() error {
	if r.Name == "" {
		return errors.New("empty name")
	}

	if r.From.IsZero() {
		return errors.New("empty from")
	}

	if r.To.IsZero() {
		return errors.New("empty to")
	}

	return nil
}

type GetByNameAndDateRangeResponse struct {
	Rates []Rate `json:"rates"`
}
