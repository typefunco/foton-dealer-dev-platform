package dealer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// Repository интерфейс репозитория Dealer.
type Repository interface {
	Create(ctx context.Context, dealer *model.Dealer) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Dealer, error)
	GetAll(ctx context.Context) ([]*model.Dealer, error)
	GetByRegion(ctx context.Context, region string) ([]*model.Dealer, error)
	Update(ctx context.Context, id int64, updates map[string]interface{}) error
	UpdateFull(ctx context.Context, dealer *model.Dealer) error
	Delete(ctx context.Context, id int64) error
	GetDealerCardData(ctx context.Context, dealerID int64, quarter string, year int) (*model.DealerCardData, error)
	AddBrand(ctx context.Context, dealerID int64, brandName string) error
	RemoveBrand(ctx context.Context, dealerID int64, brandName string) error
	GetBrands(ctx context.Context, dealerID int64) ([]string, error)
	AddBusiness(ctx context.Context, dealerID int64, businessType string) error
	RemoveBusiness(ctx context.Context, dealerID int64, businessType string) error
	GetBusinesses(ctx context.Context, dealerID int64) ([]string, error)
}

// Service сервис для работы с дилерами.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService создает новый экземпляр сервиса Dealer.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// GetDealerCard возвращает полную информацию о дилере для карточки.
func (s *Service) GetDealerCard(ctx context.Context, dealerID int64, quarter string, year int) (*model.DealerCardData, error) {
	if dealerID <= 0 {
		return nil, fmt.Errorf("DealerService.GetDealerCard: invalid dealer ID: %d", dealerID)
	}

	// Валидация квартала
	if !isValidQuarter(quarter) {
		return nil, fmt.Errorf("DealerService.GetDealerCard: invalid quarter: %s", quarter)
	}

	// Валидация года
	if year < 2020 || year > 2030 {
		return nil, fmt.Errorf("DealerService.GetDealerCard: invalid year: %d", year)
	}

	cardData, err := s.repo.GetDealerCardData(ctx, dealerID, quarter, year)
	if err != nil {
		s.logger.Error("DealerService.GetDealerCard: failed to get dealer card",
			"dealer_id", dealerID,
			"quarter", quarter,
			"year", year,
			"error", err,
		)
		return nil, fmt.Errorf("DealerService.GetDealerCard: %w", err)
	}

	s.logger.Info("DealerService.GetDealerCard: successfully retrieved card",
		"dealer_id", dealerID,
		"dealer_name", cardData.Name,
		"quarter", quarter,
		"year", year,
	)

	return cardData, nil
}

// GetDealerByID возвращает дилера по ID.
func (s *Service) GetDealerByID(ctx context.Context, id int64) (*model.Dealer, error) {
	if id <= 0 {
		return nil, fmt.Errorf("DealerService.GetDealerByID: invalid ID: %d", id)
	}

	dealer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("DealerService.GetDealerByID: failed to get dealer",
			"id", id,
			"error", err,
		)
		return nil, fmt.Errorf("DealerService.GetDealerByID: %w", err)
	}

	return dealer, nil
}

// GetAllDealers возвращает всех дилеров.
func (s *Service) GetAllDealers(ctx context.Context) ([]*model.Dealer, error) {
	dealers, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("DealerService.GetAllDealers: failed to get dealers",
			"error", err,
		)
		return nil, fmt.Errorf("DealerService.GetAllDealers: %w", err)
	}

	return dealers, nil
}

// GetDealersByRegion возвращает дилеров по региону.
func (s *Service) GetDealersByRegion(ctx context.Context, region string) ([]*model.Dealer, error) {
	dealers, err := s.repo.GetByRegion(ctx, region)
	if err != nil {
		s.logger.Error("DealerService.GetDealersByRegion: failed to get dealers",
			"region", region,
			"error", err,
		)
		return nil, fmt.Errorf("DealerService.GetDealersByRegion: %w", err)
	}

	s.logger.Info("DealerService.GetDealersByRegion: successfully retrieved dealers",
		"region", region,
		"count", len(dealers),
	)

	return dealers, nil
}

// CreateDealer создает нового дилера.
func (s *Service) CreateDealer(ctx context.Context, dealer *model.Dealer) (int64, error) {
	// Валидация
	if err := s.validateDealer(dealer); err != nil {
		return 0, fmt.Errorf("DealerService.CreateDealer: validation failed: %w", err)
	}

	id, err := s.repo.Create(ctx, dealer)
	if err != nil {
		s.logger.Error("DealerService.CreateDealer: failed to create dealer",
			"name", dealer.Name,
			"error", err,
		)
		return 0, fmt.Errorf("DealerService.CreateDealer: %w", err)
	}

	s.logger.Info("DealerService.CreateDealer: successfully created dealer",
		"id", id,
		"name", dealer.Name,
	)

	return id, nil
}

// UpdateDealer обновляет данные дилера.
func (s *Service) UpdateDealer(ctx context.Context, id int64, updates map[string]interface{}) error {
	if id <= 0 {
		return fmt.Errorf("DealerService.UpdateDealer: invalid ID: %d", id)
	}

	if len(updates) == 0 {
		return fmt.Errorf("DealerService.UpdateDealer: no fields to update")
	}

	err := s.repo.Update(ctx, id, updates)
	if err != nil {
		s.logger.Error("DealerService.UpdateDealer: failed to update dealer",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("DealerService.UpdateDealer: %w", err)
	}

	s.logger.Info("DealerService.UpdateDealer: successfully updated dealer",
		"id", id,
	)

	return nil
}

// DeleteDealer удаляет дилера.
func (s *Service) DeleteDealer(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("DealerService.DeleteDealer: invalid ID: %d", id)
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("DealerService.DeleteDealer: failed to delete dealer",
			"id", id,
			"error", err,
		)
		return fmt.Errorf("DealerService.DeleteDealer: %w", err)
	}

	s.logger.Info("DealerService.DeleteDealer: successfully deleted dealer",
		"id", id,
	)

	return nil
}

// validateDealer валидирует данные дилера.
func (s *Service) validateDealer(dealer *model.Dealer) error {
	if dealer.Name == "" {
		return fmt.Errorf("name is required")
	}

	if dealer.City == "" {
		return fmt.Errorf("city is required")
	}

	if dealer.Region == "" {
		return fmt.Errorf("region is required")
	}

	if dealer.Manager == "" {
		return fmt.Errorf("manager is required")
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
