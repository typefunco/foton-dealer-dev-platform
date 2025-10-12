package dealerdev

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// Repository интерфейс репозитория DealerDevelopment.
type Repository interface {
	Create(ctx context.Context, dd *model.DealerDevelopment) (int, error)
	GetByID(ctx context.Context, id int) (*model.DealerDevelopment, error)
	GetByDealerAndPeriod(ctx context.Context, dealerID int, period time.Time) (*model.DealerDevelopment, error)
	GetAllByPeriod(ctx context.Context, period time.Time) ([]*model.DealerDevelopment, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	UpdateFull(ctx context.Context, dd *model.DealerDevelopment) error
	Delete(ctx context.Context, id int) error
	GetWithDetailsByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.DealerDevWithDetails, error)
}

// Service сервис для работы с данными развития дилеров.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService создает новый экземпляр сервиса DealerDev.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// GetDealerDevByPeriod возвращает список данных развития дилеров за период.
func (s *Service) GetDealerDevByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.DealerDevWithDetails, error) {
	// Валидация квартала
	if !isValidQuarter(quarter) {
		return nil, fmt.Errorf("DealerDevService.GetDealerDevByPeriod: invalid quarter: %s", quarter)
	}

	// Валидация года
	if year < 2020 || year > 2030 {
		return nil, fmt.Errorf("DealerDevService.GetDealerDevByPeriod: invalid year: %d", year)
	}

	ddList, err := s.repo.GetWithDetailsByPeriod(ctx, quarter, year, region)
	if err != nil {
		s.logger.Error("DealerDevService.GetDealerDevByPeriod: failed to get dealer dev data",
			"quarter", quarter,
			"year", year,
			"region", region,
			"error", err,
		)
		return nil, fmt.Errorf("DealerDevService.GetDealerDevByPeriod: %w", err)
	}

	s.logger.Info("DealerDevService.GetDealerDevByPeriod: successfully retrieved data",
		"quarter", quarter,
		"year", year,
		"region", region,
		"count", len(ddList),
	)

	return ddList, nil
}

// GetDealerDevByID возвращает данные развития дилера по ID.
func (s *Service) GetDealerDevByID(ctx context.Context, id int) (*model.DealerDevelopment, error) {
	if id <= 0 {
		return nil, fmt.Errorf("DealerDevService.GetDealerDevByID: invalid ID: %d", id)
	}

	dd, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("DealerDevService.GetDealerDevByID: failed to get dealer dev",
			"id", id,
			"error", err,
		)
		return nil, fmt.Errorf("DealerDevService.GetDealerDevByID: %w", err)
	}

	return dd, nil
}

// CreateDealerDev создает новую запись развития дилера.
func (s *Service) CreateDealerDev(ctx context.Context, dd *model.DealerDevelopment) (int, error) {
	// Валидация
	if err := s.validateDealerDev(dd); err != nil {
		return 0, fmt.Errorf("DealerDevService.CreateDealerDev: validation failed: %w", err)
	}

	id, err := s.repo.Create(ctx, dd)
	if err != nil {
		s.logger.Error("DealerDevService.CreateDealerDev: failed to create",
			"dealer_id", dd.DealerID,
			"quarter", dd.Quarter,
			"year", dd.Year,
			"error", err,
		)
		return 0, fmt.Errorf("DealerDevService.CreateDealerDev: %w", err)
	}

	s.logger.Info("DealerDevService.CreateDealerDev: successfully created",
		"id", id,
		"dealer_id", dd.DealerID,
		"quarter", dd.Quarter,
		"year", dd.Year,
	)

	return id, nil
}

// UpdateDealerDev обновляет данные развития дилера.
func (s *Service) UpdateDealerDev(ctx context.Context, id int, updates map[string]interface{}) error {
	if id <= 0 {
		return fmt.Errorf("DealerDevService.UpdateDealerDev: invalid ID: %d", id)
	}

	if len(updates) == 0 {
		return fmt.Errorf("DealerDevService.UpdateDealerDev: no fields to update")
	}

	err := s.repo.Update(ctx, id, updates)
	if err != nil {
		s.logger.Error("DealerDevService.UpdateDealerDev: failed to update",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("DealerDevService.UpdateDealerDev: %w", err)
	}

	s.logger.Info("DealerDevService.UpdateDealerDev: successfully updated",
		"id", id,
	)

	return nil
}

// DeleteDealerDev удаляет запись развития дилера.
func (s *Service) DeleteDealerDev(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("DealerDevService.DeleteDealerDev: invalid ID: %d", id)
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("DealerDevService.DeleteDealerDev: failed to delete",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("DealerDevService.DeleteDealerDev: %w", err)
	}

	s.logger.Info("DealerDevService.DeleteDealerDev: successfully deleted",
		"id", id,
	)

	return nil
}

// validateDealerDev валидирует данные развития дилера.
func (s *Service) validateDealerDev(dd *model.DealerDevelopment) error {
	if dd.DealerID <= 0 {
		return fmt.Errorf("dealer_id is required")
	}

	// Валидация квартала и года
	if dd.Quarter == "" {
		return fmt.Errorf("quarter is required")
	}
	if dd.Year <= 0 {
		return fmt.Errorf("year is required")
	}

	if dd.CheckListScore < 0 || dd.CheckListScore > 100 {
		return fmt.Errorf("check_list_score must be between 0 and 100")
	}

	// Валидация класса дилера
	validClasses := []string{"A", "B", "C", "D"}
	isValid := false
	for _, class := range validClasses {
		if dd.DealershipClass == class {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid dealership_class: %s", dd.DealershipClass)
	}

	// Валидация рекомендации
	validRecs := []string{"Planned Result", "Needs Development", "Find New Candidate", "Close Down"}
	recValid := false
	for _, rec := range validRecs {
		if dd.DDRecommendation == rec {
			recValid = true
			break
		}
	}
	if !recValid {
		return fmt.Errorf("invalid dealer_dev_recommendation: %s", dd.DDRecommendation)
	}

	if dd.MarketingInvestments < 0 {
		return fmt.Errorf("marketing_investments cannot be negative")
	}

	return nil
}

// isValidQuarter проверяет валидность квартала.
func isValidQuarter(quarter string) bool {
	validQuarters := map[string]bool{
		"Q1": true,
		"Q2": true,
		"Q3": true,
		"Q4": true,
	}
	return validQuarters[quarter]
}

// parseQuarterToPeriod преобразует quarter и year в time.Time.
func parseQuarterToPeriod(quarter string, year int) (time.Time, error) {
	var month int
	switch quarter {
	case "Q1":
		month = 1 // Январь
	case "Q2":
		month = 4 // Апрель
	case "Q3":
		month = 7 // Июль
	case "Q4":
		month = 10 // Октябрь
	default:
		return time.Time{}, fmt.Errorf("invalid quarter: %s", quarter)
	}

	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
}
