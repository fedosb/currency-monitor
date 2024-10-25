package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/fedosb/currency-monitor/services/gateway/internal/dto"
	"github.com/fedosb/currency-monitor/services/gateway/internal/gateway/auth"
	"github.com/fedosb/currency-monitor/services/gateway/internal/gateway/currency"
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

	repo := repository.NewUserRepository()

	authGW := auth.New("http://localhost:8082")
	currencyGW := currency.New("localhost:50051")

	authSvc := service.NewAuthService(repo, authGW)
	currencySvc := service.NewCurrencyService(currencyGW)

	res, err := currencySvc.GetRateByNameAndDateRange(ctx, dto.GetByNameAndDateRangeRequest{
		Name: "usd",
		From: time.Now().Add(-time.Hour * 24),
		To:   time.Now(),
	})
	fmt.Println(res, err)

	res2, err := currencySvc.GetRateByNameAndDate(ctx, dto.GetByNameAndDateRequest{
		Name: "usd",
		Date: time.Now(),
	})
	fmt.Println(res2, err)

	err = authSvc.Init(ctx)
	fmt.Println(err)

	token, err := authSvc.SignIn(ctx, dto.SignInRequest{
		Login:    "fedor",
		Password: "qwerty",
	})
	fmt.Println(token, err)

	err = authSvc.ValidateToken(ctx, dto.ValidateTokenRequest{
		Token: token.Token,
	})
	fmt.Println(err)

	return nil
}
