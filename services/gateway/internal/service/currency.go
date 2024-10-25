package service

import (
	"context"
	"fmt"
	"github.com/fedosb/currency-monitor/services/gateway/internal/dto"
	"time"
)

type CurrencyService struct {
	currencyGateway CurrencyGateway
}

type CurrencyGateway interface {
	GetRateByNameAndDate(ctx context.Context, name string, date time.Time) (dto.Rate, error)
	GetRateByNameAndDateRange(ctx context.Context, name string, from, to time.Time) ([]dto.Rate, error)
}

func NewCurrencyService(currencyGateway CurrencyGateway) *CurrencyService {
	return &CurrencyService{
		currencyGateway: currencyGateway,
	}
}

func (s *CurrencyService) GetRateByNameAndDate(
	ctx context.Context,
	req dto.GetByNameAndDateRequest,
) (dto.GetByNameAndDateResponse, error) {
	if err := req.Validate(); err != nil {
		return dto.GetByNameAndDateResponse{}, fmt.Errorf("validate request: %w", err)
	}

	rate, err := s.currencyGateway.GetRateByNameAndDate(ctx, req.Name, req.Date)
	if err != nil {
		return dto.GetByNameAndDateResponse{}, fmt.Errorf("get rate by name and date: %w", err)
	}

	return dto.GetByNameAndDateResponse{Rate: rate}, nil
}

func (s *CurrencyService) GetRateByNameAndDateRange(
	ctx context.Context,
	req dto.GetByNameAndDateRangeRequest,
) (dto.GetByNameAndDateRangeResponse, error) {
	if err := req.Validate(); err != nil {
		return dto.GetByNameAndDateRangeResponse{}, fmt.Errorf("validate request: %w", err)
	}

	rates, err := s.currencyGateway.GetRateByNameAndDateRange(ctx, req.Name, req.From, req.To)
	if err != nil {
		return dto.GetByNameAndDateRangeResponse{}, fmt.Errorf("get rate by name and date range: %w", err)
	}

	return dto.GetByNameAndDateRangeResponse{Rates: rates}, nil
}
