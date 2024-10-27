package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/fedosb/currency-monitor/services/gateway/internal/config"
	"github.com/fedosb/currency-monitor/services/gateway/internal/gateway/auth"
	"github.com/fedosb/currency-monitor/services/gateway/internal/gateway/currency"
	handlerhttp "github.com/fedosb/currency-monitor/services/gateway/internal/handler/http"
	"github.com/fedosb/currency-monitor/services/gateway/internal/repository"
	"github.com/fedosb/currency-monitor/services/gateway/internal/service"
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
		return fmt.Errorf("config: %w", err)
	}

	repo := repository.NewUserRepository()

	authGW := auth.New(cfg.AuthAPI)
	currencyGW := currency.New(cfg.CurrencyService)

	authSvc := service.NewAuthService(repo, authGW)
	currencySvc := service.NewCurrencyService(currencyGW)

	err = authSvc.Init(ctx)
	if err != nil {
		return fmt.Errorf("auth service init: %w", err)
	}

	handler := handlerhttp.NewHandler(authSvc, currencySvc)

	httpServer := http.Server{
		Addr:    cfg.Net.GetHTTPAddress(),
		Handler: handler.HTTPHandler(),
	}

	var runGroup errgroup.Group
	runGroup.Go(func() error {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}

		return nil
	})

	runGroup.Go(func() error {
		<-ctx.Done()

		stopCtx := context.Background()

		log.Println("Shutting down http server...")
		err := httpServer.Shutdown(stopCtx)
		if err != nil {
			return fmt.Errorf("http server shutdown: %w", err)
		}

		return nil
	})

	err = runGroup.Wait()
	if err != nil {
		return fmt.Errorf("run group: %w", err)
	}

	return nil
}
