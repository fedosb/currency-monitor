package cron

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	scheduler gocron.Scheduler

	fetchInterval  time.Duration
	fetcherService FetcherService
}

type FetcherService interface {
	FetchAndUpdateRates(ctx context.Context) error
}

func NewScheduler(interval time.Duration, svc FetcherService) *Scheduler {
	s := &Scheduler{
		fetchInterval:  interval,
		fetcherService: svc,
	}

	return s
}

func (s *Scheduler) Setup(ctx context.Context) error {
	var err error

	s.scheduler, err = gocron.NewScheduler()
	if err != nil {
		return fmt.Errorf("new scheduler: %w", err)
	}

	if err = s.scheduleFetcher(ctx); err != nil {
		return fmt.Errorf("set up scheduler: %w", err)
	}

	return nil
}

func (s *Scheduler) scheduleFetcher(ctx context.Context) error {
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

func (s *Scheduler) Start() {
	s.scheduler.Start()
}

func (s *Scheduler) Stop() error {
	if err := s.scheduler.Shutdown(); err != nil {
		return fmt.Errorf("shutdown scheduler: %w", err)
	}

	return nil
}
