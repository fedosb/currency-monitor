package transport

import (
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errsinternal "github.com/fedosb/currency-monitor/services/currency/internal/errors"
)

func handleError(err error) error {
	if err == nil {
		return nil
	}

	log.Println(err)

	var (
		code codes.Code
		msg  string
	)

	validationErr := &errsinternal.ValidationError{}
	switch {
	case errors.As(err, &validationErr):
		code = codes.InvalidArgument
		msg = validationErr.Error()
	case errors.Is(err, errsinternal.ErrNotFound):
		code = codes.NotFound
		msg = errsinternal.ErrNotFound.Error()
	default:
		code = codes.Internal
		msg = errsinternal.ErrDefault.Error()
	}

	return status.Error(code, msg)
}
