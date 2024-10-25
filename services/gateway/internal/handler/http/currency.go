package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fedosb/currency-monitor/services/gateway/internal/dto"
)

func (h *Handler) getByNameAndDate(c *gin.Context) {
	name := c.Param("name")
	date, _ := time.Parse(time.DateOnly, c.Query("date"))

	rate, err := h.currencySvc.GetRateByNameAndDate(c.Request.Context(), dto.GetByNameAndDateRequest{
		Name: name,
		Date: date,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respond(c, http.StatusOK, rate)
}

func (h *Handler) getByNameAndRange(c *gin.Context) {
	name := c.Param("name")
	from, _ := time.Parse(time.DateOnly, c.Query("from"))
	to, _ := time.Parse(time.DateOnly, c.Query("to"))

	rate, err := h.currencySvc.GetRateByNameAndDateRange(c.Request.Context(), dto.GetByNameAndDateRangeRequest{
		Name: name,
		From: from,
		To:   to,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respond(c, http.StatusOK, rate)
}
