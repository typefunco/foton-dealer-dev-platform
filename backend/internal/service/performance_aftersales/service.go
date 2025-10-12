package performance_aftersales

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// Repository интерфейс репозитория PerformanceAfterSales.
type Repository interface {
	Create(ctx context.Context, pas *model.PerformanceAfterSales) (int, error)
	GetByID(ctx context.Context, id int) (*model.PerformanceAfterSales, error)
	GetByDealerIDAndPeriod(ctx context.Context, dealerID int, period time.Time) (*model.PerformanceAfterSales, error)
	GetAllByPeriod(ctx context.Context, period time.Time) ([]*model.PerformanceAfterSales, error)
	Update(ctx context.Context, pas *model.PerformanceAfterSales) error
	Delete(ctx context.Context, id int) error
}

// Service сервис для работы с данными производительности запчастей.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService создает новый экземпляр сервиса PerformanceAfterSales.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// Create создает новую запись производительности запчастей.
func (s *Service) Create(ctx context.Context, pas *model.PerformanceAfterSales) (int, error) {
	if pas.DealerID <= 0 {
		return 0, fmt.Errorf("PerformanceAfterSalesService.Create: invalid dealer ID: %d", pas.DealerID)
	}

	if pas.Period.IsZero() {
		return 0, fmt.Errorf("PerformanceAfterSalesService.Create: invalid period")
	}

	id, err := s.repo.Create(ctx, pas)
	if err != nil {
		s.logger.Error("Failed to create performance aftersales", "error", err, "dealer_id", pas.DealerID)
		return 0, fmt.Errorf("PerformanceAfterSalesService.Create: %w", err)
	}

	s.logger.Info("Performance aftersales created successfully", "id", id, "dealer_id", pas.DealerID)
	return id, nil
}

// GetByID получает производительность запчастей по ID.
func (s *Service) GetByID(ctx context.Context, id int) (*model.PerformanceAfterSales, error) {
	if id <= 0 {
		return nil, fmt.Errorf("PerformanceAfterSalesService.GetByID: invalid ID: %d", id)
	}

	pas, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get performance aftersales", "error", err, "id", id)
		return nil, fmt.Errorf("PerformanceAfterSalesService.GetByID: %w", err)
	}

	return pas, nil
}

// GetByDealerIDAndPeriod получает производительность запчастей по ID дилера и периоду.
func (s *Service) GetByDealerIDAndPeriod(ctx context.Context, dealerID int, period time.Time) (*model.PerformanceAfterSales, error) {
	if dealerID <= 0 {
		return nil, fmt.Errorf("PerformanceAfterSalesService.GetByDealerIDAndPeriod: invalid dealer ID: %d", dealerID)
	}

	if period.IsZero() {
		return nil, fmt.Errorf("PerformanceAfterSalesService.GetByDealerIDAndPeriod: invalid period")
	}

	pas, err := s.repo.GetByDealerIDAndPeriod(ctx, dealerID, period)
	if err != nil {
		s.logger.Error("Failed to get performance aftersales", "error", err, "dealer_id", dealerID, "period", period)
		return nil, fmt.Errorf("PerformanceAfterSalesService.GetByDealerIDAndPeriod: %w", err)
	}

	return pas, nil
}

// GetAllByPeriod получает все записи производительности запчастей за период.
func (s *Service) GetAllByPeriod(ctx context.Context, period time.Time) ([]*model.PerformanceAfterSales, error) {
	if period.IsZero() {
		return nil, fmt.Errorf("PerformanceAfterSalesService.GetAllByPeriod: invalid period")
	}

	performances, err := s.repo.GetAllByPeriod(ctx, period)
	if err != nil {
		s.logger.Error("Failed to get performance aftersales by period", "error", err, "period", period)
		return nil, fmt.Errorf("PerformanceAfterSalesService.GetAllByPeriod: %w", err)
	}

	return performances, nil
}

// Update обновляет запись производительности запчастей.
func (s *Service) Update(ctx context.Context, pas *model.PerformanceAfterSales) error {
	if pas.ID <= 0 {
		return fmt.Errorf("PerformanceAfterSalesService.Update: invalid ID: %d", pas.ID)
	}

	if pas.DealerID <= 0 {
		return fmt.Errorf("PerformanceAfterSalesService.Update: invalid dealer ID: %d", pas.DealerID)
	}

	if pas.Period.IsZero() {
		return fmt.Errorf("PerformanceAfterSalesService.Update: invalid period")
	}

	err := s.repo.Update(ctx, pas)
	if err != nil {
		s.logger.Error("Failed to update performance aftersales", "error", err, "id", pas.ID)
		return fmt.Errorf("PerformanceAfterSalesService.Update: %w", err)
	}

	s.logger.Info("Performance aftersales updated successfully", "id", pas.ID)
	return nil
}

// Delete удаляет запись производительности запчастей.
func (s *Service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("PerformanceAfterSalesService.Delete: invalid ID: %d", id)
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete performance aftersales", "error", err, "id", id)
		return fmt.Errorf("PerformanceAfterSalesService.Delete: %w", err)
	}

	s.logger.Info("Performance aftersales deleted successfully", "id", id)
	return nil
}
