package dealer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// Repository интерфейс репозитория Dealer.
type Repository interface {
	Create(ctx context.Context, dealer *model.Dealer) (int, error)
	GetByID(ctx context.Context, id int) (*model.Dealer, error)
	GetAll(ctx context.Context) ([]*model.Dealer, error)
	GetByRegion(ctx context.Context, region string) ([]*model.Dealer, error)
	GetWithFilters(ctx context.Context, filters *model.FilterParams) ([]*model.Dealer, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	UpdateFull(ctx context.Context, dealer *model.Dealer) error
	Delete(ctx context.Context, id int) error
	GetDealerCardData(ctx context.Context, dealerID int, period time.Time) (*model.DealerCardData, error)
	AddBrand(ctx context.Context, dealerID int, brandName string) error
	RemoveBrand(ctx context.Context, dealerID int, brandName string) error
	GetBrands(ctx context.Context, dealerID int) ([]string, error)
	AddBusiness(ctx context.Context, dealerID int, businessType string) error
	RemoveBusiness(ctx context.Context, dealerID int, businessType string) error
	GetBusinesses(ctx context.Context, dealerID int) ([]string, error)
}

// ExcelRepository интерфейс репозитория для работы с Excel данными дилеров.
type ExcelRepository interface {
	GetDealerNetTableName(year int, quarter string) string
	TableExists(ctx context.Context, year int, quarter string) (bool, error)
	GetDealersWithFilters(ctx context.Context, year int, quarter string, filters *model.FilterParams) ([]*model.Dealer, error)
	GetDealerCardData(ctx context.Context, year int, quarter string, dealerName string) (*model.DealerCardData, error)
	GetAvailableRegions(ctx context.Context, year int, quarter string) ([]string, error)
}

// Service сервис для работы с дилерами.
type Service struct {
	repo      Repository
	excelRepo ExcelRepository
	logger    *slog.Logger
}

// NewService создает новый экземпляр сервиса Dealer.
func NewService(repo Repository, excelRepo ExcelRepository, logger *slog.Logger) *Service {
	return &Service{
		repo:      repo,
		excelRepo: excelRepo,
		logger:    logger,
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

	// Преобразуем quarter/year в period
	period, err := parseQuarterToPeriod(quarter, year)
	if err != nil {
		return nil, fmt.Errorf("DealerService.GetDealerCard: invalid quarter/year: %w", err)
	}

	cardData, err := s.repo.GetDealerCardData(ctx, int(dealerID), period)
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
		"dealer_name", cardData.DealerNameRu,
		"quarter", quarter,
		"year", year,
	)

	return cardData, nil
}

// GetDealerByID возвращает дилера по ID.
func (s *Service) GetDealerByID(ctx context.Context, id int) (*model.Dealer, error) {
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

// GetDealersWithFilters возвращает дилеров с применением фильтров.
func (s *Service) GetDealersWithFilters(ctx context.Context, filters *model.FilterParams) ([]*model.Dealer, error) {
	// Валидация фильтров
	if err := filters.Validate(); err != nil {
		return nil, fmt.Errorf("DealerService.GetDealersWithFilters: validation failed: %w", err)
	}

	// Если нет параметров year и quarter, используем данные из Excel по умолчанию
	// Используем последний доступный квартал (2025 Q3)
	if filters.Year == 0 && filters.Quarter == "" {
		s.logger.Info("No year/quarter specified, using default Excel data",
			slog.Int("default_year", 2025),
			slog.String("default_quarter", "Q3"),
		)

		// Проверяем существование Excel таблицы
		exists, err := s.excelRepo.TableExists(ctx, 2025, "Q3")
		if err != nil {
			s.logger.Error("Failed to check Excel table existence", "error", err)
			return nil, fmt.Errorf("failed to check Excel table existence: %w", err)
		}

		if exists {
			// Используем Excel данные
			dealers, err := s.excelRepo.GetDealersWithFilters(ctx, 2025, "Q3", filters)
			if err != nil {
				s.logger.Error("DealerService.GetDealersWithFilters: failed to get dealers from Excel",
					"filters", filters,
					"error", err,
				)
				return nil, fmt.Errorf("DealerService.GetDealersWithFilters: %w", err)
			}

			s.logger.Info("DealerService.GetDealersWithFilters: successfully retrieved dealers from Excel",
				"filters", filters,
				"count", len(dealers),
			)

			return dealers, nil
		} else {
			s.logger.Warn("Excel table does not exist, falling back to old data")
		}
	}

	// Используем старые данные как fallback
	dealers, err := s.repo.GetWithFilters(ctx, filters)
	if err != nil {
		s.logger.Error("DealerService.GetDealersWithFilters: failed to get dealers",
			"filters", filters,
			"error", err,
		)
		return nil, fmt.Errorf("DealerService.GetDealersWithFilters: %w", err)
	}

	s.logger.Info("DealerService.GetDealersWithFilters: successfully retrieved dealers",
		"filters", filters,
		"count", len(dealers),
	)

	return dealers, nil
}

// CreateDealer создает нового дилера.
func (s *Service) CreateDealer(ctx context.Context, dealer *model.Dealer) (int, error) {
	// Валидация
	if err := s.validateDealer(dealer); err != nil {
		return 0, fmt.Errorf("DealerService.CreateDealer: validation failed: %w", err)
	}

	id, err := s.repo.Create(ctx, dealer)
	if err != nil {
		s.logger.Error("DealerService.CreateDealer: failed to create dealer",
			"name", dealer.DealerNameRu,
			"error", err,
		)
		return 0, fmt.Errorf("DealerService.CreateDealer: %w", err)
	}

	s.logger.Info("DealerService.CreateDealer: successfully created dealer",
		"id", id,
		"name", dealer.DealerNameRu,
	)

	return id, nil
}

// UpdateDealer обновляет данные дилера.
func (s *Service) UpdateDealer(ctx context.Context, id int, updates map[string]interface{}) error {
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
func (s *Service) DeleteDealer(ctx context.Context, id int) error {
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
	if dealer.DealerNameRu == "" {
		return fmt.Errorf("dealer_name_ru is required")
	}

	if dealer.DealerNameEn == "" {
		return fmt.Errorf("dealer_name_en is required")
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

// parseQuarterToPeriod преобразует quarter и year в time.Time.
func parseQuarterToPeriod(quarter string, year int) (time.Time, error) {
	var month int
	switch quarter {
	case "q1", "Q1":
		month = 1 // Январь
	case "q2", "Q2":
		month = 4 // Апрель
	case "q3", "Q3":
		month = 7 // Июль
	case "q4", "Q4":
		month = 10 // Октябрь
	default:
		return time.Time{}, fmt.Errorf("invalid quarter: %s", quarter)
	}

	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
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

// GetDealersFromExcel получает дилеров из Excel таблицы с фильтрами.
func (s *Service) GetDealersFromExcel(ctx context.Context, year int, quarter string, filters *model.FilterParams) ([]*model.Dealer, error) {
	s.logger.Info("Getting dealers from Excel table",
		slog.Int("year", year),
		slog.String("quarter", quarter),
		slog.String("region", filters.Region),
		slog.Int("limit", filters.Limit),
	)

	// Проверяем валидность параметров
	if year < 2020 || year > 2030 {
		return nil, fmt.Errorf("invalid year: %d", year)
	}

	if !isValidQuarter(quarter) {
		return nil, fmt.Errorf("invalid quarter: %s", quarter)
	}

	// Проверяем существование таблицы
	exists, err := s.excelRepo.TableExists(ctx, year, quarter)
	if err != nil {
		return nil, fmt.Errorf("failed to check table existence: %w", err)
	}

	if !exists {
		s.logger.Warn("Excel table does not exist",
			slog.Int("year", year),
			slog.String("quarter", quarter),
		)
		return []*model.Dealer{}, nil
	}

	// Получаем дилеров из Excel таблицы
	dealers, err := s.excelRepo.GetDealersWithFilters(ctx, year, quarter, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get dealers from Excel: %w", err)
	}

	s.logger.Info("Dealers retrieved from Excel",
		slog.Int("year", year),
		slog.String("quarter", quarter),
		slog.Int("count", len(dealers)),
	)

	return dealers, nil
}

// GetDealerCardFromExcel получает карточку дилера из Excel таблицы.
func (s *Service) GetDealerCardFromExcel(ctx context.Context, year int, quarter string, dealerName string) (*model.DealerCardData, error) {
	s.logger.Info("Getting dealer card from Excel table",
		slog.Int("year", year),
		slog.String("quarter", quarter),
		slog.String("dealer_name", dealerName),
	)

	// Проверяем валидность параметров
	if year < 2020 || year > 2030 {
		return nil, fmt.Errorf("invalid year: %d", year)
	}

	if !isValidQuarter(quarter) {
		return nil, fmt.Errorf("invalid quarter: %s", quarter)
	}

	if dealerName == "" {
		return nil, fmt.Errorf("dealer name cannot be empty")
	}

	// Проверяем существование таблицы
	exists, err := s.excelRepo.TableExists(ctx, year, quarter)
	if err != nil {
		return nil, fmt.Errorf("failed to check table existence: %w", err)
	}

	if !exists {
		s.logger.Warn("Excel table does not exist",
			slog.Int("year", year),
			slog.String("quarter", quarter),
		)
		return &model.DealerCardData{}, nil
	}

	// Получаем карточку дилера из Excel таблицы
	cardData, err := s.excelRepo.GetDealerCardData(ctx, year, quarter, dealerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get dealer card from Excel: %w", err)
	}

	s.logger.Info("Dealer card retrieved from Excel",
		slog.Int("year", year),
		slog.String("quarter", quarter),
		slog.String("dealer_name", dealerName),
	)

	return cardData, nil
}

// GetAvailableRegionsFromExcel получает список доступных регионов из Excel таблицы.
func (s *Service) GetAvailableRegionsFromExcel(ctx context.Context, year int, quarter string) ([]string, error) {
	s.logger.Info("Getting available regions from Excel table",
		slog.Int("year", year),
		slog.String("quarter", quarter),
	)

	// Проверяем валидность параметров
	if year < 2020 || year > 2030 {
		return nil, fmt.Errorf("invalid year: %d", year)
	}

	if !isValidQuarter(quarter) {
		return nil, fmt.Errorf("invalid quarter: %s", quarter)
	}

	// Проверяем существование таблицы
	exists, err := s.excelRepo.TableExists(ctx, year, quarter)
	if err != nil {
		return nil, fmt.Errorf("failed to check table existence: %w", err)
	}

	if !exists {
		s.logger.Warn("Excel table does not exist",
			slog.Int("year", year),
			slog.String("quarter", quarter),
		)
		return []string{}, nil
	}

	// Получаем регионы из Excel таблицы
	regions, err := s.excelRepo.GetAvailableRegions(ctx, year, quarter)
	if err != nil {
		return nil, fmt.Errorf("failed to get regions from Excel: %w", err)
	}

	s.logger.Info("Available regions retrieved from Excel",
		slog.Int("year", year),
		slog.String("quarter", quarter),
		slog.Int("count", len(regions)),
	)

	return regions, nil
}
