package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *Handler) authMiddleware(c *gin.Context) {
	fmt.Println("Auth middleware")
	c.Next()
}
