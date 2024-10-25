package dto

import (
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
	Name string
	Date time.Time
}

type GetByNameAndDateResponse struct {
	Rate Rate `json:"rate"`
}

type GetByNameAndDateRangeRequest struct {
	Name string
	From time.Time
	To   time.Time
}

type GetByNameAndDateRangeResponse struct {
	Rates []Rate `json:"rates"`
}
