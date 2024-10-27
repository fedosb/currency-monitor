package http

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/fedosb/currency-monitor/services/gateway/internal/dto"
	errsinternal "github.com/fedosb/currency-monitor/services/gateway/internal/errors"
)

const AuthorizationHeader = "Authorization"

func (h *Handler) authMiddleware(c *gin.Context) {
	bearerToken := strings.Split(c.GetHeader(AuthorizationHeader), " ")
	if len(bearerToken) != 2 {
		err := errsinternal.NewAuthError(errsinternal.AuthErrorInvalidTokenFormatMsg)

		respondError(c, fmt.Errorf("invalid authorization header: %s", err))
		c.Abort()
		return
	}

	token := strings.TrimSpace(bearerToken[1])
	err := h.authSvc.ValidateToken(c.Request.Context(), dto.ValidateTokenRequest{Token: token})
	if err != nil {
		respondError(c, err)
		c.Abort()
	}

	c.Next()
}

func (h *Handler) logMiddleware(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path

	c.Next()

	// Stop timer
	latency := time.Now().Sub(start)

	clientIP := c.ClientIP()
	method := c.Request.Method
	statusCode := c.Writer.Status()
	errorMessage := c.Errors.String()

	log.Info().
		Int("status_code", statusCode).
		Dur("latency", latency).
		Str("client_ip", clientIP).
		Str("method", method).
		Str("path", path).
		Err(fmt.Errorf(errorMessage)).
		Send()
}
