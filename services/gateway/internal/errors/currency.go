package errors

import "google.golang.org/grpc/codes"

type CurrencyError struct {
	Message string
	Code    codes.Code
}

func (e CurrencyError) Error() string {
	return e.Message
}

func NewCurrencyError(msg string, code codes.Code) *CurrencyError {
	return &CurrencyError{Message: msg, Code: code}
}
