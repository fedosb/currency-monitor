package service

import (
	"context"
	"fmt"
	"github.com/fedosb/currency-monitor/services/currency/internal/config"
	"github.com/fedosb/currency-monitor/services/currency/internal/entity"
	"log"
	"time"

	"github.com/fedosb/currency-monitor/services/currency/internal/dto"
)

type FetchJobRepository interface {
	Create(ctx context.Context, job entity.FetchJob) (entity.FetchJob, error)
	GetPlanned(ctx context.Context) ([]entity.FetchJob, error)
	GetFailed(ctx context.Context) ([]entity.FetchJob, error)
	Update(ctx context.Context, job entity.FetchJob) (entity.FetchJob, error)
}

type FetcherService struct {
	gateway       CurrencyAPIGateway
	jobRepository FetchJobRepository

	rateService *RateService

	jobRescheduleInterval       time.Duration
	failedJobRescheduleInterval time.Duration
	failedJobMaxRetries         int
	saveRatesMaxWorkers         int
}

type CurrencyAPIGateway interface {
	GetRubRates(ctx context.Context) (dto.RubRates, error)
}

func NewFetcherService(
	gateway CurrencyAPIGateway,
	jobRepo FetchJobRepository,
	rateSvc *RateService,
	cfg config.FetcherConfig,
) *FetcherService {
	return &FetcherService{
		gateway:                     gateway,
		jobRepository:               jobRepo,
		rateService:                 rateSvc,
		jobRescheduleInterval:       cfg.GetJobRescheduleInterval(),
		failedJobRescheduleInterval: cfg.GetFailedJobRescheduleInterval(),
		failedJobMaxRetries:         cfg.GetFailedJobMaxRetries(),
		saveRatesMaxWorkers:         cfg.GetSaveRatesMaxWorkers(),
	}
}

func (s *FetcherService) FetchAndUpdateRates(ctx context.Context) error {

	jobs, err := s.jobRepository.GetPlanned(ctx)
	if err != nil {
		return fmt.Errorf("get planned jobs: %w", err)
	}

	if len(jobs) == 0 {
		return nil
	}

	rubRates, err := s.gateway.GetRubRates(ctx)
	if err != nil {
		return fmt.Errorf("get rub rates: %w", err)
	}

	// TODO: parallelize processing
	for _, job := range jobs {
		var err error
		defer func() {
			if err != nil {
				s.failJob(ctx, job, err.Error())
			}
		}()

		if err = s.processJob(ctx, job, rubRates); err != nil {
			log.Println("Failed to process job:", err)
		}
	}

	return nil
}

func (s *FetcherService) processJob(ctx context.Context, job entity.FetchJob, rubRates dto.RubRates) error {
	rate := entity.Rate{
		Name: job.Name,
		Date: rubRates.Date,
		Rate: rubRates.Rates[job.Name],
	}

	_, err := s.rateService.Save(ctx, rate)
	if err != nil {
		return fmt.Errorf("create rate: %w", err)
	}

	job.Succeed()
	job.Reschedule(s.jobRescheduleInterval)

	_, err = s.jobRepository.Update(ctx, job)
	if err != nil {
		return fmt.Errorf("update job: %w", err)
	}

	return nil
}

func (s *FetcherService) failJob(ctx context.Context, job entity.FetchJob, reason string) {
	job.Fail(reason)

	_, err := s.jobRepository.Update(ctx, job)
	if err != nil {
		log.Println("Error failing job:", err)
	}
}

func (s *FetcherService) RequeueFailedJobs(ctx context.Context) error {
	jobs, err := s.jobRepository.GetFailed(ctx)
	if err != nil {
		return fmt.Errorf("get failed jobs: %w", err)
	}

	for _, job := range jobs {
		if *job.Retries >= s.failedJobMaxRetries {
			log.Printf("Job %d (%s) reached max retries\n", job.ID, job.Name)
			job.Reschedule(s.jobRescheduleInterval)
			job.Retries = nil
		} else {
			job.Reschedule(s.failedJobRescheduleInterval)
			*job.Retries++
		}

		job.Status = entity.JobStatusReady

		_, err := s.jobRepository.Update(ctx, job)
		if err != nil {
			return fmt.Errorf("update job: %w", err)
		}
	}

	return nil
}
