package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/fedosb/currency-monitor/services/currency/internal/config"
	"github.com/fedosb/currency-monitor/services/currency/internal/dto"
	"github.com/fedosb/currency-monitor/services/currency/internal/entity"
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
	processJobMaxWorkers        int
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
		processJobMaxWorkers:        cfg.GetProcessJobMaxWorkers(),
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

	// This application only supports RUB rates for now
	rubRates, err := s.gateway.GetRubRates(ctx)
	if err != nil {
		return fmt.Errorf("get rub rates: %w", err)
	}

	var (
		wg  = sync.WaitGroup{}
		sem = semaphore.NewWeighted(int64(s.processJobMaxWorkers))
	)

	for _, job := range jobs {
		err := sem.Acquire(ctx, 1)
		if err != nil {
			return fmt.Errorf("acquiring semaphore: %w", err)
		}

		wg.Add(1)
		go func(job entity.FetchJob) {
			defer sem.Release(1)
			defer wg.Done()

			err := s.processJob(ctx, job, rubRates)
			if err != nil {
				log.Println("Error processing job:", err)
			}
		}(job)
	}

	wg.Wait()

	return nil
}

func (s *FetcherService) processJob(ctx context.Context, job entity.FetchJob, rubRates dto.RubRates) error {
	var err error
	defer func() {
		if err != nil {
			s.failJob(ctx, job, err.Error())
		}
	}()

	rate := entity.Rate{
		Name: job.Name,
		Date: rubRates.Date,
		Rate: rubRates.Rates[job.Name],
	}

	_, err = s.rateService.Save(ctx, rate)
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
