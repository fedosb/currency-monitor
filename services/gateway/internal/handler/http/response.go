package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type jsonResponse struct {
	Data  any     `json:"data,omitempty"`
	Error *string `json:"error,omitempty"`
}

func respond(c *gin.Context, status int, data any) {
	c.JSON(status, jsonResponse{Data: data})
}

func respondError(c *gin.Context, err error) {
	msg := err.Error()
	c.JSON(http.StatusBadRequest, jsonResponse{Error: &msg})
}
