package service

import (
	"context"
	"fmt"

	"github.com/fedosb/currency-monitor/services/gateway/dto"
	"github.com/fedosb/currency-monitor/services/gateway/internal/entity"
	hashutil "github.com/fedosb/currency-monitor/services/gateway/internal/utils/hash"
)

type AuthService struct {
	userRepository UserRepository
	authGateway    AuthGateway
}

type UserRepository interface {
	Create(ctx context.Context, user entity.User) error
	GetByLogin(ctx context.Context, login string) (entity.User, error)
}

type AuthGateway interface {
	GenerateToken(ctx context.Context, login string) (string, error)
	ValidateToken(ctx context.Context, token string) error
}

func NewAuthService(userRepository UserRepository, authGateway AuthGateway) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		authGateway:    authGateway,
	}
}

func (s *AuthService) Init(ctx context.Context) error {
	users := map[string]string{
		"user":  "password",
		"fedor": "qwerty",
	}

	for login, password := range users {
		hash, err := hashutil.CreatePasswordHash(password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		user := entity.User{
			Login:        login,
			PasswordHash: hash,
		}

		err = s.userRepository.Create(ctx, user)
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}
	}

	return nil
}

func (s *AuthService) SignIn(ctx context.Context, req dto.SignInRequest) (dto.SignInResponse, error) {
	err := req.Validate()
	if err != nil {
		return dto.SignInResponse{}, fmt.Errorf("validate request: %w", err)
	}

	user, err := s.userRepository.GetByLogin(ctx, req.Login)
	if err != nil {
		return dto.SignInResponse{}, fmt.Errorf("get user by login: %w", err)
	}

	if !hashutil.CheckPasswordHash(req.Password, user.PasswordHash) {
		return dto.SignInResponse{}, fmt.Errorf("invalid password")
	}

	token, err := s.authGateway.GenerateToken(ctx, user.Login)
	if err != nil {
		return dto.SignInResponse{}, fmt.Errorf("generate token: %w", err)
	}

	return dto.SignInResponse{Token: token}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req dto.ValidateTokenRequest) error {
	err := s.authGateway.ValidateToken(ctx, req.Token)
	if err != nil {
		return fmt.Errorf("validate token: %w", err)
	}

	return nil
}
