package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	errsinternal "github.com/fedosb/currency-monitor/services/gateway/internal/errors"
	codeutil "github.com/fedosb/currency-monitor/services/gateway/internal/utils/codes"
)

type jsonResponse struct {
	Data  any     `json:"data,omitempty"`
	Error *string `json:"error,omitempty"`
}

func respond(c *gin.Context, status int, data any) {
	c.JSON(status, jsonResponse{Data: data})
}

func respondError(c *gin.Context, err error) {
	if err == nil {
		c.JSON(http.StatusInternalServerError, jsonResponse{})
		return
	}

	var (
		msg  string
		code int
	)

	currencyError := &errsinternal.GatewayError{}
	authError := &errsinternal.AuthError{}
	switch {
	case errors.As(err, &currencyError):
		msg = currencyError.Error()
		code = codeutil.GRPCCodeToHTTPStatus(currencyError.Code)
	case errors.As(err, &authError):
		msg = authError.Message
		code = http.StatusUnauthorized
	default:
		code = http.StatusInternalServerError
		msg = errsinternal.DefaultError.Error()
	}

	c.JSON(code, jsonResponse{Error: &msg})
}
