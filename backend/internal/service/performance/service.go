package performance

import (
	"context"
	"github.com/typefunco/dealer_dev_platform/internal/model"
	"log/slog"
)

type Repository interface {
	FindPerformances(ctx context.Context, region string) ([]*model.Performance, error)
}

type Service struct {
	repository Repository
	logger     *slog.Logger
}

// NewService конструктор.
func NewService(repository Repository, logger *slog.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) FindPerformances(ctx context.Context, region string) ([]*model.Performance, error) {
	return s.repository.FindPerformances(ctx, region)
}
