package performance

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/typefunco/dealer_dev_platform/internal/model"
)

type Repository interface {
	FindPerformances(ctx context.Context, region string) ([]*model.PerformanceSales, error)
	GetWithDetailsByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.PerformanceWithDetails, error)
	GetByID(ctx context.Context, id int64) (*model.PerformanceSales, error)
	Create(ctx context.Context, perf *model.PerformanceSales) (int64, error)
	Update(ctx context.Context, id int64, updates map[string]interface{}) error
	Delete(ctx context.Context, id int64) error
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

// FindPerformances возвращает производительность по региону (deprecated).
func (s *Service) FindPerformances(ctx context.Context, region string) ([]*model.PerformanceSales, error) {
	return s.repository.FindPerformances(ctx, region)
}

// GetPerformanceByPeriod возвращает список данных производительности за период.
func (s *Service) GetPerformanceByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.PerformanceWithDetails, error) {
	// Валидация квартала
	if !isValidQuarter(quarter) {
		return nil, fmt.Errorf("PerformanceService.GetPerformanceByPeriod: invalid quarter: %s", quarter)
	}

	// Валидация года
	if year < 2020 || year > 2030 {
		return nil, fmt.Errorf("PerformanceService.GetPerformanceByPeriod: invalid year: %d", year)
	}

	perfList, err := s.repository.GetWithDetailsByPeriod(ctx, quarter, year, region)
	if err != nil {
		s.logger.Error("PerformanceService.GetPerformanceByPeriod: failed to get performance data",
			"quarter", quarter,
			"year", year,
			"region", region,
			"error", err,
		)
		return nil, fmt.Errorf("PerformanceService.GetPerformanceByPeriod: %w", err)
	}

	s.logger.Info("PerformanceService.GetPerformanceByPeriod: successfully retrieved data",
		"quarter", quarter,
		"year", year,
		"region", region,
		"count", len(perfList),
	)

	return perfList, nil
}

// GetPerformanceByID возвращает данные производительности по ID.
func (s *Service) GetPerformanceByID(ctx context.Context, id int64) (*model.PerformanceSales, error) {
	if id <= 0 {
		return nil, fmt.Errorf("PerformanceService.GetPerformanceByID: invalid ID: %d", id)
	}

	perf, err := s.repository.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("PerformanceService.GetPerformanceByID: failed to get performance",
			"id", id,
			"error", err,
		)
		return nil, fmt.Errorf("PerformanceService.GetPerformanceByID: %w", err)
	}

	return perf, nil
}

// CreatePerformance создает новую запись производительности.
func (s *Service) CreatePerformance(ctx context.Context, perf *model.PerformanceSales) (int64, error) {
	// Валидация
	if err := s.validatePerformance(perf); err != nil {
		return 0, fmt.Errorf("PerformanceService.CreatePerformance: validation failed: %w", err)
	}

	id, err := s.repository.Create(ctx, perf)
	if err != nil {
		s.logger.Error("PerformanceService.CreatePerformance: failed to create",
			"dealer_id", perf.DealerID,
			"period", perf.Period,
			"error", err,
		)
		return 0, fmt.Errorf("PerformanceService.CreatePerformance: %w", err)
	}

	s.logger.Info("PerformanceService.CreatePerformance: successfully created",
		"id", id,
		"dealer_id", perf.DealerID,
		"period", perf.Period,
	)

	return id, nil
}

// UpdatePerformance обновляет данные производительности.
func (s *Service) UpdatePerformance(ctx context.Context, id int64, updates map[string]interface{}) error {
	if id <= 0 {
		return fmt.Errorf("PerformanceService.UpdatePerformance: invalid ID: %d", id)
	}

	if len(updates) == 0 {
		return fmt.Errorf("PerformanceService.UpdatePerformance: no fields to update")
	}

	err := s.repository.Update(ctx, id, updates)
	if err != nil {
		s.logger.Error("PerformanceService.UpdatePerformance: failed to update",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("PerformanceService.UpdatePerformance: %w", err)
	}

	s.logger.Info("PerformanceService.UpdatePerformance: successfully updated",
		"id", id,
	)

	return nil
}

// DeletePerformance удаляет запись производительности.
func (s *Service) DeletePerformance(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("PerformanceService.DeletePerformance: invalid ID: %d", id)
	}

	err := s.repository.Delete(ctx, id)
	if err != nil {
		s.logger.Error("PerformanceService.DeletePerformance: failed to delete",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("PerformanceService.DeletePerformance: %w", err)
	}

	s.logger.Info("PerformanceService.DeletePerformance: successfully deleted",
		"id", id,
	)

	return nil
}

// validatePerformance валидирует данные производительности.
func (s *Service) validatePerformance(perf *model.PerformanceSales) error {
	if perf.DealerID <= 0 {
		return fmt.Errorf("dealer_id is required")
	}

	// Валидация периода
	if perf.Period.IsZero() {
		return fmt.Errorf("period is required")
	}

	// Валидация финансовых показателей (если не nil)
	if perf.SalesRevenue != nil && *perf.SalesRevenue < 0 {
		return fmt.Errorf("sales_revenue cannot be negative")
	}

	if perf.SalesCost != nil && *perf.SalesCost < 0 {
		return fmt.Errorf("sales_cost cannot be negative")
	}

	if perf.SalesMargin != nil && *perf.SalesMargin < 0 {
		return fmt.Errorf("sales_margin cannot be negative")
	}

	return nil
}

// isValidQuarter проверяет валидность квартала.
func isValidQuarter(quarter string) bool {
	validQuarters := map[string]bool{
		"q1": true,
		"q2": true,
		"q3": true,
		"q4": true,
	}
	return validQuarters[quarter]
}
