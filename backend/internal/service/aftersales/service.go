package aftersales

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// Repository интерфейс репозитория AfterSales.
type Repository interface {
	Create(ctx context.Context, as *model.AfterSales) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.AfterSales, error)
	GetByDealerAndPeriod(ctx context.Context, dealerID int64, quarter string, year int) (*model.AfterSales, error)
	GetAllByPeriod(ctx context.Context, quarter string, year int) ([]*model.AfterSales, error)
	Update(ctx context.Context, id int64, updates map[string]interface{}) error
	UpdateFull(ctx context.Context, as *model.AfterSales) error
	Delete(ctx context.Context, id int64) error
	GetWithDetailsByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.AfterSalesWithDetails, error)
}

// Service сервис для работы с данными послепродажного обслуживания.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService создает новый экземпляр сервиса AfterSales.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// GetAfterSalesByPeriod возвращает список данных послепродажного обслуживания за период.
func (s *Service) GetAfterSalesByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.AfterSalesWithDetails, error) {
	// Валидация квартала
	if !isValidQuarter(quarter) {
		return nil, fmt.Errorf("AfterSalesService.GetAfterSalesByPeriod: invalid quarter: %s", quarter)
	}

	// Валидация года
	if year < 2020 || year > 2030 {
		return nil, fmt.Errorf("AfterSalesService.GetAfterSalesByPeriod: invalid year: %d", year)
	}

	afterSalesList, err := s.repo.GetWithDetailsByPeriod(ctx, quarter, year, region)
	if err != nil {
		s.logger.Error("AfterSalesService.GetAfterSalesByPeriod: failed to get after sales data",
			"quarter", quarter,
			"year", year,
			"region", region,
			"error", err,
		)
		return nil, fmt.Errorf("AfterSalesService.GetAfterSalesByPeriod: %w", err)
	}

	s.logger.Info("AfterSalesService.GetAfterSalesByPeriod: successfully retrieved data",
		"quarter", quarter,
		"year", year,
		"region", region,
		"count", len(afterSalesList),
	)

	return afterSalesList, nil
}

// GetAfterSalesByID возвращает данные послепродажного обслуживания по ID.
func (s *Service) GetAfterSalesByID(ctx context.Context, id int64) (*model.AfterSales, error) {
	if id <= 0 {
		return nil, fmt.Errorf("AfterSalesService.GetAfterSalesByID: invalid ID: %d", id)
	}

	afterSales, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("AfterSalesService.GetAfterSalesByID: failed to get after sales",
			"id", id,
			"error", err,
		)
		return nil, fmt.Errorf("AfterSalesService.GetAfterSalesByID: %w", err)
	}

	return afterSales, nil
}

// CreateAfterSales создает новую запись послепродажного обслуживания.
func (s *Service) CreateAfterSales(ctx context.Context, as *model.AfterSales) (int64, error) {
	// Валидация
	if err := s.validateAfterSales(as); err != nil {
		return 0, fmt.Errorf("AfterSalesService.CreateAfterSales: validation failed: %w", err)
	}

	id, err := s.repo.Create(ctx, as)
	if err != nil {
		s.logger.Error("AfterSalesService.CreateAfterSales: failed to create",
			"dealer_id", as.DealerID,
			"quarter", as.Quarter,
			"year", as.Year,
			"error", err,
		)
		return 0, fmt.Errorf("AfterSalesService.CreateAfterSales: %w", err)
	}

	s.logger.Info("AfterSalesService.CreateAfterSales: successfully created",
		"id", id,
		"dealer_id", as.DealerID,
		"quarter", as.Quarter,
		"year", as.Year,
	)

	return id, nil
}

// UpdateAfterSales обновляет данные послепродажного обслуживания.
func (s *Service) UpdateAfterSales(ctx context.Context, id int64, updates map[string]interface{}) error {
	if id <= 0 {
		return fmt.Errorf("AfterSalesService.UpdateAfterSales: invalid ID: %d", id)
	}

	if len(updates) == 0 {
		return fmt.Errorf("AfterSalesService.UpdateAfterSales: no fields to update")
	}

	err := s.repo.Update(ctx, id, updates)
	if err != nil {
		s.logger.Error("AfterSalesService.UpdateAfterSales: failed to update",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("AfterSalesService.UpdateAfterSales: %w", err)
	}

	s.logger.Info("AfterSalesService.UpdateAfterSales: successfully updated",
		"id", id,
	)

	return nil
}

// DeleteAfterSales удаляет запись послепродажного обслуживания.
func (s *Service) DeleteAfterSales(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("AfterSalesService.DeleteAfterSales: invalid ID: %d", id)
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("AfterSalesService.DeleteAfterSales: failed to delete",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("AfterSalesService.DeleteAfterSales: %w", err)
	}

	s.logger.Info("AfterSalesService.DeleteAfterSales: successfully deleted",
		"id", id,
	)

	return nil
}

// validateAfterSales валидирует данные послепродажного обслуживания.
func (s *Service) validateAfterSales(as *model.AfterSales) error {
	if as.DealerID <= 0 {
		return fmt.Errorf("dealer_id is required")
	}

	if !isValidQuarter(as.Quarter) {
		return fmt.Errorf("invalid quarter: %s (must be q1, q2, q3, or q4)", as.Quarter)
	}

	if as.Year < 2020 || as.Year > 2030 {
		return fmt.Errorf("invalid year: %d (must be between 2020 and 2030)", as.Year)
	}

	if as.RecommendedStock < 0 || as.RecommendedStock > 100 {
		return fmt.Errorf("recommended_stock must be between 0 and 100")
	}

	if as.WarrantyStock < 0 || as.WarrantyStock > 100 {
		return fmt.Errorf("warranty_stock must be between 0 and 100")
	}

	if as.FotonLaborHours < 0 || as.FotonLaborHours > 100 {
		return fmt.Errorf("foton_labor_hours must be between 0 and 100")
	}

	if as.ServiceContracts < 0 {
		return fmt.Errorf("service_contracts cannot be negative")
	}

	// Валидация решения
	validDecisions := map[model.AfterSalesDecision]bool{
		model.AfterSalesDecisionPlannedResult:    true,
		model.AfterSalesDecisionNeedsDevelopment: true,
		model.AfterSalesDecisionFindNewCandidate: true,
		model.AfterSalesDecisionCloseDown:        true,
	}

	if !validDecisions[as.AfterSalesDecision] {
		return fmt.Errorf("invalid as_decision: %s", as.AfterSalesDecision)
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
