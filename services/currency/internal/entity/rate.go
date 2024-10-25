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

const (
	RateValidationNameRequiredMsg      = "correct name is required"
	RateValidationDateRequiredMsg      = "correct date is required"
	RateValidationFromAndToRequiredMsg = "correct from and to are required"
	RateValidationFromBeforeToMsg      = "from must be before to"
)
