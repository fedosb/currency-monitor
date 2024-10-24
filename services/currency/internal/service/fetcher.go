package service

import (
	"context"
	"fmt"
	"github.com/fedosb/currency-monitor/services/currency/internal/entity"

	"github.com/fedosb/currency-monitor/services/currency/internal/dto"
)

type FetcherService struct {
	Gateway CurrencyAPIGateway

	RateService *RateService
}

type CurrencyAPIGateway interface {
	GetRubRates(ctx context.Context) (dto.RubRates, error)
}

func NewFetcherService(gateway CurrencyAPIGateway, rateSvc *RateService) *FetcherService {
	return &FetcherService{Gateway: gateway, RateService: rateSvc}
}

func (s *FetcherService) FetchAndUpdateRates(ctx context.Context) error {

	rubRates, err := s.Gateway.GetRubRates(ctx)
	if err != nil {
		return fmt.Errorf("get rub rates: %w", err)
	}

	const usd = "usd" // TODO: make dynamic

	rate := entity.Rate{
		Name: usd,
		Date: rubRates.Date,
		Rate: rubRates.Rates[usd],
	}

	_, err = s.RateService.Save(ctx, rate)
	if err != nil {
		return fmt.Errorf("create rate: %w", err)
	}

	return nil
}
