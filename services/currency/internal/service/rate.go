package service

import (
	"context"
	"fmt"
	"time"

	"github.com/fedosb/currency-monitor/services/currency/internal/dto"
	"github.com/fedosb/currency-monitor/services/currency/internal/entity"
)

type RateService struct {
	repository RateRepository
}

type RateRepository interface {
	Save(ctx context.Context, rate entity.Rate) (entity.Rate, error)
	GetByNameAndDate(ctx context.Context, name string, date time.Time) (entity.Rate, error)
	GetByNameAndDateRange(ctx context.Context, name string, from, to time.Time) ([]entity.Rate, error)
}

func NewRateService(repository RateRepository) *RateService {
	return &RateService{repository: repository}
}

func (s *RateService) Save(ctx context.Context, rate entity.Rate) (entity.Rate, error) {
	rate, err := s.repository.Save(ctx, rate)
	if err != nil {
		return entity.Rate{}, fmt.Errorf("create in repository: %w", err)
	}

	return rate, nil
}

func (s *RateService) GetByNameAndDate(ctx context.Context, request dto.GetByNameAndDateRequest) (dto.GetByNameAndDateResponse, error) {
	rate, err := s.repository.GetByNameAndDate(ctx, request.Name, request.Date)
	if err != nil {
		return dto.GetByNameAndDateResponse{}, fmt.Errorf("get from repository: %w", err)
	}

	return dto.GetByNameAndDateResponse{Rate: rate}, nil
}

func (s *RateService) GetByNameAndDateRange(ctx context.Context, request dto.GetByNameAndDateRangeRequest) (dto.GetByNameAndDateRangeResponse, error) {
	if err := request.Validate(); err != nil {
		return dto.GetByNameAndDateRangeResponse{}, fmt.Errorf("validate request: %w", err)
	}

	rates, err := s.repository.GetByNameAndDateRange(ctx, request.Name, request.From, request.To)
	if err != nil {
		return dto.GetByNameAndDateRangeResponse{}, fmt.Errorf("get from repository: %w", err)
	}

	return dto.GetByNameAndDateRangeResponse{Rates: rates}, nil
}
