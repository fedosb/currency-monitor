package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/fedosb/currency-monitor/services/currency/internal/config"
	"github.com/fedosb/currency-monitor/services/currency/internal/db/postgres"
	"github.com/fedosb/currency-monitor/services/currency/internal/repository"
	"github.com/fedosb/currency-monitor/services/currency/internal/service"
	grpcTransport "github.com/fedosb/currency-monitor/services/currency/transport/grpc"
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

	if cfg.DB.Migrate() {
		fmt.Println("Applying migrations...")
		if err := db.ApplyMigrations(); err != nil {
			return fmt.Errorf("applying migrations: %w", err)
		}
	}

	rateRepo := repository.NewRateRepository(db)

	rateSvc := service.NewRateService(rateRepo)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	grpcServer := grpcTransport.NewGRPCServer(rateSvc)

	var runGroup errgroup.Group

	runGroup.Go(func() error {
		fmt.Println("Starting gRPC server...")
		err := grpcServer.Serve(listener)
		if err != nil {
			return err
		}

		return nil
	})

	runGroup.Go(func() error {
		<-ctx.Done()

		grpcServer.GracefulStop()
		fmt.Println("gRPC server has been stopped")

		return nil
	})

	err = runGroup.Wait()
	if err != nil {
		return fmt.Errorf("run group: %w", err)
	}

	return nil
}
