package http

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/fedosb/currency-monitor/services/gateway/internal/dto"
)

func (h *Handler) authMiddleware(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			_ = c.Error(err)
		}
	}()

	bearerToken := strings.Split(c.GetHeader("Authorization"), " ")
	if len(bearerToken) != 2 {
		err = errors.New("invalid token format, use Bearer <token>")
		respondError(c, err)
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
