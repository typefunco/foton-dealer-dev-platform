package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/typefunco/dealer_dev_platform/internal/model"
	"github.com/typefunco/dealer_dev_platform/internal/utils/jwt"
)

type Repository interface {
	CreateUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, login string) error
	GetUser(ctx context.Context, login string) (*model.User, error)
	Ping(ctx context.Context) error
}

type JWTRepository interface {
	ValidateJWT(jwt string) (*jwt.JWTClaims, error)
	ValidateJWTLegacy(jwt string) error
	GenerateJWT(login string, isAdmin bool, role string) (string, error)
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
func (s *Service) Login(ctx context.Context, login string, password string) (*jwt.JWTClaims, error) {
	if login == "" || password == "" {
		return nil, fmt.Errorf("AuthService.Login username or password is empty")
	}

	user, err := s.repo.GetUser(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("AuthService.GetUser no user %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("AuthService.GetUser user not found")
	}

	// В реальном приложении здесь должна быть проверка хеша пароля
	if user.Password != password {
		return nil, fmt.Errorf("AuthService.Login invalid password")
	}

	// Генерируем JWT с информацией о пользователе
	token, err := s.jwt.GenerateJWT(user.Login, user.IsAdmin, string(user.Role))
	if err != nil {
		s.logger.Error("AuthService.GenerateJWT failed to generate JWT", "error", err)
		return nil, fmt.Errorf("AuthService.GenerateJWT failed to generate JWT: %w", err)
	}

	// Валидируем токен для получения claims
	claims, err := s.jwt.ValidateJWT(token)
	if err != nil {
		s.logger.Error("AuthService.ValidateJWT failed to validate generated JWT", "error", err)
		return nil, fmt.Errorf("AuthService.ValidateJWT failed to validate generated JWT: %w", err)
	}

	return claims, nil
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

	jwt, err := s.jwt.GenerateJWT(user.Login, user.IsAdmin, string(user.Role))
	if err != nil {
		return "", fmt.Errorf("AuthService.GenerateJWT error %w", err)
	}

	return jwt, nil
}

func (s *Service) PingDatabase(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

// GenerateToken генерирует JWT токен для пользователя
func (s *Service) GenerateToken(login string, isAdmin bool, role string) (string, error) {
	return s.jwt.GenerateJWT(login, isAdmin, role)
}
