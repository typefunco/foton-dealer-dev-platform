package auth

import (
	"context"
	"fmt"
	"github.com/typefunco/dealer_dev_platform/internal/model"
	"log/slog"
)

type Repository interface {
	CreateUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, login string) error
	GetUser(ctx context.Context, login string) (*model.User, error)
	Ping(ctx context.Context) error
}

type JWTRepository interface {
	ValidateJWT(jwt string) error
	GenerateJWT(login string) (string, error)
}

type Service struct {
	repo   Repository
	jwt    JWTRepository
	logger *slog.Logger
}

// NewService конструктор auth Service.
func NewService(repo Repository, jwt JWTRepository, logger *slog.Logger) *Service {
	return &Service{repo: repo, jwt: jwt, logger: logger}
}

// Login метод логина пользователя.
func (s *Service) Login(ctx context.Context, login string, password string) error {
	if login == "" || password == "" {
		return fmt.Errorf("AuthService.Login username or password is empty")
	}
	err := s.jwt.ValidateJWT(login)
	if err != nil {
		s.logger.Error("AuthService.ValidateJWT failed to validate JWT", "error", err)
	}

	user, err := s.repo.GetUser(ctx, login)
	if err != nil {
		return fmt.Errorf("AuthService.GetUser no user %w", err)
	}

	if user == nil {
		return fmt.Errorf("AuthService.GetUser user not found")
	}
	return nil
}

// Signup создает JWT и ошибку.
func (s *Service) Signup(ctx context.Context, user model.User) (string, error) {
	if user.Login == "" || user.Password == "" {
		return "", fmt.Errorf("AuthService.Signup username or password is empty")
	}

	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("AuthService.CreateUser error %w", err)
	}

	jwt, err := s.jwt.GenerateJWT(user.Login)
	if err != nil {
		return "", fmt.Errorf("AuthService.GenerateJWT error %w", err)
	}

	return jwt, nil
}

func (s *Service) PingDatabase(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
