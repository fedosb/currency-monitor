package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/fedosb/currency-monitor/services/currency/internal/config"
	"github.com/fedosb/currency-monitor/services/currency/internal/cron"
	"github.com/fedosb/currency-monitor/services/currency/internal/db/postgres"
	"github.com/fedosb/currency-monitor/services/currency/internal/gateway/currency_api"
	"github.com/fedosb/currency-monitor/services/currency/internal/repository"
	"github.com/fedosb/currency-monitor/services/currency/internal/service"
	grpcTransport "github.com/fedosb/currency-monitor/services/currency/internal/transport/grpc"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	db, err := postgres.New(cfg.DB)
	if err != nil {
		return fmt.Errorf("creating db: %w", err)
	}

	if cfg.DB.GetApplyMigrations() {
		log.Println("Applying migrations...")
		if err := db.ApplyMigrations(); err != nil {
			return fmt.Errorf("applying migrations: %w", err)
		}
	}

	rateRepo := repository.NewRateRepository(db)
	fetchJobsRepo := repository.NewFetchJobRepository(db)

	currencyApiGateway := currency_api.New(cfg.CurrencyAPI)

	rateSvc := service.NewRateService(rateRepo)
	fetchSvc := service.NewFetcherService(
		currencyApiGateway,
		fetchJobsRepo,
		rateSvc,
		cfg.Fetcher,
	)

	fetchSvc.Init(ctx, cfg.Fetcher.GetCurrencies())

	if cfg.Fetcher.GetRunImmediately() {
		log.Println("Running fetcher immediately...")
		err := fetchSvc.FetchAndUpdateRates(ctx)
		if err != nil {
			return fmt.Errorf("running fetcher: %w", err)
		}
	}

	grpcServer := grpcTransport.NewGRPCServer(rateSvc)

	scheduler := cron.NewScheduler(cfg.Fetcher, fetchSvc)
	err = scheduler.Setup(ctx)
	if err != nil {
		return fmt.Errorf("scheduler setup: %w", err)
	}

	var runGroup errgroup.Group
	runGroup.Go(func() error {
		log.Println("Starting gRPC server...")
		err := grpcServer.Serve(cfg.Net)
		if err != nil {
			return err
		}

		return nil
	})

	runGroup.Go(func() error {
		log.Println("Starting scheduler...")
		scheduler.Start()
		return nil
	})

	runGroup.Go(func() error {
		<-ctx.Done()

		log.Println("Shutting down scheduler...")
		err := scheduler.Stop()
		if err != nil {
			log.Println("scheduler stop:", err)
		}

		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()

		return nil
	})

	err = runGroup.Wait()
	if err != nil {
		return fmt.Errorf("run group: %w", err)
	}

	return nil
}