package cron

import (
	"context"
	"fmt"
	"github.com/fedosb/currency-monitor/services/currency/internal/config"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	scheduler gocron.Scheduler

	fetchInterval   time.Duration
	requeueInterval time.Duration

	fetcherService FetcherService
}

type FetcherService interface {
	FetchAndUpdateRates(ctx context.Context) error
	RequeueFailedJobs(ctx context.Context) error
}

func NewScheduler(cfg config.FetcherConfig, svc FetcherService) *Scheduler {
	s := &Scheduler{
		fetchInterval:   cfg.GetFetchInterval(),
		requeueInterval: cfg.GetJobRescheduleInterval(),
		fetcherService:  svc,
	}

	return s
}

func (s *Scheduler) Setup(ctx context.Context) error {
	var err error

	s.scheduler, err = gocron.NewScheduler()
	if err != nil {
		return fmt.Errorf("new scheduler: %w", err)
	}

	if err = s.setupFetchAndUpdateRates(ctx); err != nil {
		return fmt.Errorf("set up fetch and update scheduler: %w", err)
	}

	if err = s.setupRequeueFailedJobs(ctx); err != nil {
		return fmt.Errorf("set up requeue failed scheduler: %w", err)
	}

	return nil
}

func (s *Scheduler) setupFetchAndUpdateRates(ctx context.Context) error {
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(s.fetchInterval),
		gocron.NewTask(
			func(svc FetcherService, ctx context.Context) {
				log.Println("Fetching rates")

				ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.fetchInterval)
				defer cancel()

				err := svc.FetchAndUpdateRates(ctx)
				if err != nil {
					log.Println(err)
				}
			}, s.fetcherService, ctx,
		),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		return fmt.Errorf("failed to schedule fetcher: %w", err)
	}

	return nil
}

func (s *Scheduler) setupRequeueFailedJobs(ctx context.Context) error {
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(s.requeueInterval),
		gocron.NewTask(
			func(svc FetcherService, ctx context.Context) {
				log.Println("Requeuing failed jobs")

				ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.requeueInterval)
				defer cancel()

				err := svc.RequeueFailedJobs(ctx)
				if err != nil {
					log.Println(err)
				}
			}, s.fetcherService, ctx,
		),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		return fmt.Errorf("failed to schedule fetcher: %w", err)
	}

	return nil
}

func (s *Scheduler) Start() {
	s.scheduler.Start()
}

func (s *Scheduler) Stop() error {
	if err := s.scheduler.Shutdown(); err != nil {
		return fmt.Errorf("shutdown scheduler: %w", err)
	}

	return nil
}
