package main

import (
	"context"
	"fmt"
	"github.com/fedosb/currency-monitor/services/gateway/internal/dto"
	"github.com/fedosb/currency-monitor/services/gateway/internal/service"
	"log"
	"os/signal"
	"syscall"

	"github.com/fedosb/currency-monitor/services/gateway/internal/gateway/auth"
	"github.com/fedosb/currency-monitor/services/gateway/internal/repository"
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

	authSvc := service.NewAuthService(repo, authGW)

	err := authSvc.Init(ctx)
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
