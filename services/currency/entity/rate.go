package entity

import "time"

type Rate struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Date      time.Time
	Rate      float64
}
