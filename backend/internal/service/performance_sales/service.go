package performance_sales

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// Repository интерфейс репозитория PerformanceSales.
type Repository interface {
	Create(ctx context.Context, ps *model.PerformanceSales) (int, error)
	GetByID(ctx context.Context, id int) (*model.PerformanceSales, error)
	GetByDealerIDAndPeriod(ctx context.Context, dealerID int, period time.Time) (*model.PerformanceSales, error)
	GetAllByPeriod(ctx context.Context, period time.Time) ([]*model.PerformanceSales, error)
	Update(ctx context.Context, ps *model.PerformanceSales) error
	Delete(ctx context.Context, id int) error
}

// Service сервис для работы с данными производительности продаж.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService создает новый экземпляр сервиса PerformanceSales.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// Create создает новую запись производительности продаж.
func (s *Service) Create(ctx context.Context, ps *model.PerformanceSales) (int, error) {
	if ps.DealerID <= 0 {
		return 0, fmt.Errorf("PerformanceSalesService.Create: invalid dealer ID: %d", ps.DealerID)
	}

	if ps.Quarter == "" || ps.Year == 0 {
		return 0, fmt.Errorf("PerformanceSalesService.Create: invalid period")
	}

	id, err := s.repo.Create(ctx, ps)
	if err != nil {
		s.logger.Error("Failed to create performance sales", "error", err, "dealer_id", ps.DealerID)
		return 0, fmt.Errorf("PerformanceSalesService.Create: %w", err)
	}

	s.logger.Info("Performance sales created successfully", "id", id, "dealer_id", ps.DealerID)
	return id, nil
}

// GetByID получает производительность продаж по ID.
func (s *Service) GetByID(ctx context.Context, id int) (*model.PerformanceSales, error) {
	if id <= 0 {
		return nil, fmt.Errorf("PerformanceSalesService.GetByID: invalid ID: %d", id)
	}

	ps, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get performance sales", "error", err, "id", id)
		return nil, fmt.Errorf("PerformanceSalesService.GetByID: %w", err)
	}

	return ps, nil
}

// GetByDealerIDAndPeriod получает производительность продаж по ID дилера и периоду.
func (s *Service) GetByDealerIDAndPeriod(ctx context.Context, dealerID int, period time.Time) (*model.PerformanceSales, error) {
	if dealerID <= 0 {
		return nil, fmt.Errorf("PerformanceSalesService.GetByDealerIDAndPeriod: invalid dealer ID: %d", dealerID)
	}

	if period.IsZero() {
		return nil, fmt.Errorf("PerformanceSalesService.GetByDealerIDAndPeriod: invalid period")
	}

	ps, err := s.repo.GetByDealerIDAndPeriod(ctx, dealerID, period)
	if err != nil {
		s.logger.Error("Failed to get performance sales", "error", err, "dealer_id", dealerID, "period", period)
		return nil, fmt.Errorf("PerformanceSalesService.GetByDealerIDAndPeriod: %w", err)
	}

	return ps, nil
}

// GetAllByPeriod получает все записи производительности продаж за период.
func (s *Service) GetAllByPeriod(ctx context.Context, period time.Time) ([]*model.PerformanceSales, error) {
	if period.IsZero() {
		return nil, fmt.Errorf("PerformanceSalesService.GetAllByPeriod: invalid period")
	}

	performances, err := s.repo.GetAllByPeriod(ctx, period)
	if err != nil {
		s.logger.Error("Failed to get performance sales by period", "error", err, "period", period)
		return nil, fmt.Errorf("PerformanceSalesService.GetAllByPeriod: %w", err)
	}

	return performances, nil
}

// Update обновляет запись производительности продаж.
func (s *Service) Update(ctx context.Context, ps *model.PerformanceSales) error {
	if ps.ID <= 0 {
		return fmt.Errorf("PerformanceSalesService.Update: invalid ID: %d", ps.ID)
	}

	if ps.DealerID <= 0 {
		return fmt.Errorf("PerformanceSalesService.Update: invalid dealer ID: %d", ps.DealerID)
	}

	if ps.Quarter == "" || ps.Year == 0 {
		return fmt.Errorf("PerformanceSalesService.Update: invalid period")
	}

	err := s.repo.Update(ctx, ps)
	if err != nil {
		s.logger.Error("Failed to update performance sales", "error", err, "id", ps.ID)
		return fmt.Errorf("PerformanceSalesService.Update: %w", err)
	}

	s.logger.Info("Performance sales updated successfully", "id", ps.ID)
	return nil
}

// Delete удаляет запись производительности продаж.
func (s *Service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("PerformanceSalesService.Delete: invalid ID: %d", id)
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete performance sales", "error", err, "id", id)
		return fmt.Errorf("PerformanceSalesService.Delete: %w", err)
	}

	s.logger.Info("Performance sales deleted successfully", "id", id)
	return nil
}
