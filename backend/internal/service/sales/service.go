package sales

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// Repository интерфейс репозитория Sales.
type Repository interface {
	Create(ctx context.Context, sales *model.Sales) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Sales, error)
	GetByDealerAndPeriod(ctx context.Context, dealerID int64, quarter string, year int) (*model.Sales, error)
	GetAllByPeriod(ctx context.Context, quarter string, year int) ([]*model.Sales, error)
	Update(ctx context.Context, id int64, updates map[string]interface{}) error
	UpdateFull(ctx context.Context, sales *model.Sales) error
	Delete(ctx context.Context, id int64) error
	GetWithDetailsByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.SalesWithDetails, error)
}

// Service сервис для работы с данными продаж.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService создает новый экземпляр сервиса Sales.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// GetSalesByPeriod возвращает список данных продаж за период.
func (s *Service) GetSalesByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.SalesWithDetails, error) {
	// Валидация квартала
	if !isValidQuarter(quarter) {
		return nil, fmt.Errorf("SalesService.GetSalesByPeriod: invalid quarter: %s", quarter)
	}

	// Валидация года
	if year < 2020 || year > 2030 {
		return nil, fmt.Errorf("SalesService.GetSalesByPeriod: invalid year: %d", year)
	}

	salesList, err := s.repo.GetWithDetailsByPeriod(ctx, quarter, year, region)
	if err != nil {
		s.logger.Error("SalesService.GetSalesByPeriod: failed to get sales data",
			"quarter", quarter,
			"year", year,
			"region", region,
			"error", err,
		)
		return nil, fmt.Errorf("SalesService.GetSalesByPeriod: %w", err)
	}

	s.logger.Info("SalesService.GetSalesByPeriod: successfully retrieved data",
		"quarter", quarter,
		"year", year,
		"region", region,
		"count", len(salesList),
	)

	return salesList, nil
}

// GetSalesByID возвращает данные продаж по ID.
func (s *Service) GetSalesByID(ctx context.Context, id int64) (*model.Sales, error) {
	if id <= 0 {
		return nil, fmt.Errorf("SalesService.GetSalesByID: invalid ID: %d", id)
	}

	sales, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("SalesService.GetSalesByID: failed to get sales",
			"id", id,
			"error", err,
		)
		return nil, fmt.Errorf("SalesService.GetSalesByID: %w", err)
	}

	return sales, nil
}

// CreateSales создает новую запись продаж.
func (s *Service) CreateSales(ctx context.Context, sales *model.Sales) (int64, error) {
	// Валидация
	if err := s.validateSales(sales); err != nil {
		return 0, fmt.Errorf("SalesService.CreateSales: validation failed: %w", err)
	}

	id, err := s.repo.Create(ctx, sales)
	if err != nil {
		s.logger.Error("SalesService.CreateSales: failed to create",
			"dealer_id", sales.DealerID,
			"period", sales.Period,
			"error", err,
		)
		return 0, fmt.Errorf("SalesService.CreateSales: %w", err)
	}

	s.logger.Info("SalesService.CreateSales: successfully created",
		"id", id,
		"dealer_id", sales.DealerID,
		"period", sales.Period,
	)

	return id, nil
}

// UpdateSales обновляет данные продаж.
func (s *Service) UpdateSales(ctx context.Context, id int64, updates map[string]interface{}) error {
	if id <= 0 {
		return fmt.Errorf("SalesService.UpdateSales: invalid ID: %d", id)
	}

	if len(updates) == 0 {
		return fmt.Errorf("SalesService.UpdateSales: no fields to update")
	}

	err := s.repo.Update(ctx, id, updates)
	if err != nil {
		s.logger.Error("SalesService.UpdateSales: failed to update",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("SalesService.UpdateSales: %w", err)
	}

	s.logger.Info("SalesService.UpdateSales: successfully updated",
		"id", id,
	)

	return nil
}

// DeleteSales удаляет запись продаж.
func (s *Service) DeleteSales(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("SalesService.DeleteSales: invalid ID: %d", id)
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("SalesService.DeleteSales: failed to delete",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("SalesService.DeleteSales: %w", err)
	}

	s.logger.Info("SalesService.DeleteSales: successfully deleted",
		"id", id,
	)

	return nil
}

// validateSales валидирует данные продаж.
func (s *Service) validateSales(sales *model.Sales) error {
	if sales.DealerID <= 0 {
		return fmt.Errorf("dealer_id is required")
	}

	// Валидация периода
	if sales.Period.IsZero() {
		return fmt.Errorf("period is required")
	}

	// Валидация решения (если не nil)
	if sales.SalesRecommendation != nil {
		validDecisions := []string{"Planned Result", "Needs Development", "Find New Candidate", "Close Down"}
		isValid := false
		for _, decision := range validDecisions {
			if *sales.SalesRecommendation == decision {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid sales_recommendation: %s", *sales.SalesRecommendation)
		}
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
