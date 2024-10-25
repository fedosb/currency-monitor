package http

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/fedosb/currency-monitor/services/gateway/internal/dto"
	errsinternal "github.com/fedosb/currency-monitor/services/gateway/internal/errors"
)

const AuthorizationHeader = "Authorization"

func (h *Handler) authMiddleware(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			_ = c.Error(err)
		}
	}()

	bearerToken := strings.Split(c.GetHeader(AuthorizationHeader), " ")
	if len(bearerToken) != 2 {
		err = errsinternal.NewAuthError(errsinternal.AuthErrorInvalidTokenFormatMsg)

		respondError(c, fmt.Errorf("invalid authorization header: %s", err))
		c.Abort()
		return
	}

	token := strings.TrimSpace(bearerToken[1])
	err = h.authSvc.ValidateToken(c.Request.Context(), dto.ValidateTokenRequest{Token: token})
	if err != nil {
		respondError(c, err)
		c.Abort()
	}

	c.Next()
}
