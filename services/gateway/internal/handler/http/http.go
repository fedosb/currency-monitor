package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fedosb/currency-monitor/services/gateway/internal/dto"
)

type Handler struct {
	router *gin.Engine

	authSvc     AuthService
	currencySvc CurrencyService
}

type AuthService interface {
	SignIn(ctx context.Context, req dto.SignInRequest) (dto.SignInResponse, error)
	ValidateToken(ctx context.Context, req dto.ValidateTokenRequest) error
}

type CurrencyService interface {
	GetRateByNameAndDate(ctx context.Context, req dto.GetByNameAndDateRequest) (dto.GetByNameAndDateResponse, error)
	GetRateByNameAndDateRange(ctx context.Context, req dto.GetByNameAndDateRangeRequest) (dto.GetByNameAndDateRangeResponse, error)
}

func NewHandler(authSvc AuthService, currencySvc CurrencyService) *Handler {
	router := gin.Default()

	router.Use(gin.Recovery())

	h := Handler{
		router:      router,
		authSvc:     authSvc,
		currencySvc: currencySvc,
	}

	api := router.Group("/api")
	{
		api.GET("/health", h.health)
		api.POST("/sign-in", h.signIn)
		authorized := api.Group("/", h.authMiddleware)
		{
			authorized.GET("/currency/:name", h.getByNameAndDate)
			authorized.GET("/currency/:name/range", h.getByNameAndRange)
		}
	}

	return &h
}

func (h *Handler) HTTPHandler() http.Handler {
	return http.Handler(h.router)
}
