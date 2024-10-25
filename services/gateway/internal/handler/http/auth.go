package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fedosb/currency-monitor/services/gateway/internal/dto"
)

func (h *Handler) signIn(c *gin.Context) {
	var req dto.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, err)
		return
	}

	token, err := h.authSvc.SignIn(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}

	respond(c, http.StatusOK, token)
}
