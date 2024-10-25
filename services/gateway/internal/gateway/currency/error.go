package currency

import (
	"google.golang.org/grpc/status"

	errsinternal "github.com/fedosb/currency-monitor/services/gateway/internal/errors"
)

func wrapError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	return errsinternal.NewGatewayError(st.Message(), st.Code())
}
