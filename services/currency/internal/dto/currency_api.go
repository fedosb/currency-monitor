package dto

import "time"

type RubRates struct {
	Date  time.Time
	Rates map[string]float64
}
