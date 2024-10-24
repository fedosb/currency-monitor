package errors

import "errors"

var (
	InvalidDateRangeError = errors.New("invalid date range")
	NotFoundError         = errors.New("not found")
)
