package errors

import "google.golang.org/grpc/codes"

type GatewayError struct {
	Message string
	Code    codes.Code
}

func (e GatewayError) Error() string {
	return e.Message
}

func NewGatewayError(msg string, code codes.Code) *GatewayError {
	return &GatewayError{Message: msg, Code: code}
}
