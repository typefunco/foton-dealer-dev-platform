package dealerdev

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// Repository интерфейс репозитория DealerDev.
type Repository interface {
	Create(ctx context.Context, dd *model.DealerDev) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.DealerDev, error)
	GetByDealerAndPeriod(ctx context.Context, dealerID int64, quarter string, year int) (*model.DealerDev, error)
	GetAllByPeriod(ctx context.Context, quarter string, year int) ([]*model.DealerDev, error)
	Update(ctx context.Context, id int64, updates map[string]interface{}) error
	UpdateFull(ctx context.Context, dd *model.DealerDev) error
	Delete(ctx context.Context, id int64) error
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
func (s *Service) GetDealerDevByID(ctx context.Context, id int64) (*model.DealerDev, error) {
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
func (s *Service) CreateDealerDev(ctx context.Context, dd *model.DealerDev) (int64, error) {
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
func (s *Service) UpdateDealerDev(ctx context.Context, id int64, updates map[string]interface{}) error {
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
func (s *Service) DeleteDealerDev(ctx context.Context, id int64) error {
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
func (s *Service) validateDealerDev(dd *model.DealerDev) error {
	if dd.DealerID <= 0 {
		return fmt.Errorf("dealer_id is required")
	}

	if !isValidQuarter(dd.Quarter) {
		return fmt.Errorf("invalid quarter: %s (must be q1, q2, q3, or q4)", dd.Quarter)
	}

	if dd.Year < 2020 || dd.Year > 2030 {
		return fmt.Errorf("invalid year: %d (must be between 2020 and 2030)", dd.Year)
	}

	if dd.CheckListScore < 0 || dd.CheckListScore > 100 {
		return fmt.Errorf("check_list_score must be between 0 and 100")
	}

	// Валидация класса дилера
	validClasses := map[model.DealerDevClass]bool{
		model.DealerDevClassA: true,
		model.DealerDevClassB: true,
		model.DealerDevClassC: true,
		model.DealerDevClassD: true,
	}

	if !validClasses[dd.DealerShipClass] {
		return fmt.Errorf("invalid dealer_ship_class: %s", dd.DealerShipClass)
	}

	// Валидация рекомендации
	validRecs := map[model.DealerDevRecommendation]bool{
		model.RecommendationPlannedResult:    true,
		model.RecommendationNeedsDevelopment: true,
		model.RecommendationFindNewCandidate: true,
		model.RecommendationCloseDown:        true,
	}

	if !validRecs[dd.Recommendation] {
		return fmt.Errorf("invalid dealer_dev_recommendation: %s", dd.Recommendation)
	}

	if dd.MarketingInvestments < 0 {
		return fmt.Errorf("marketing_investments cannot be negative")
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
