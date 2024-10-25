package errors

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrDefault  = errors.New("something went wrong")
)
