package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

const dealerTableName = "dealers"

// DealerRepository репозиторий для работы с дилерами.
type DealerRepository struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

// NewDealerRepository конструктор.
func NewDealerRepository(pool *pgxpool.Pool) *DealerRepository {
	return &DealerRepository{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Create создает нового дилера.
func (r *DealerRepository) Create(ctx context.Context, dealer *model.Dealer) (int, error) {
	now := time.Now()
	dealer.CreatedAt = now
	dealer.UpdatedAt = now

	query := r.sq.Insert(dealerTableName).
		Columns("name", "city", "region", "manager", "created_at", "updated_at").
		Values(dealer.DealerNameRu, dealer.City, dealer.Region, dealer.Manager, dealer.CreatedAt, dealer.UpdatedAt).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("DealerRepository.Create: error building query: %w", err)
	}

	var id int64
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("DealerRepository.Create: error inserting: %w", err)
	}

	dealer.DealerID = int(id)
	return int(id), nil
}

// GetByID получает дилера по ID.
func (r *DealerRepository) GetByID(ctx context.Context, id int) (*model.Dealer, error) {
	query := r.sq.Select("id", "name", "city", "region", "manager", "created_at", "updated_at").
		From(dealerTableName).
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetByID: error building query: %w", err)
	}

	dealer := &model.Dealer{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&dealer.DealerID, &dealer.DealerNameRu, &dealer.City, &dealer.Region, &dealer.Manager,
		&dealer.CreatedAt, &dealer.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetByID: error scanning: %w", err)
	}

	return dealer, nil
}

// GetAll получает всех дилеров.
func (r *DealerRepository) GetAll(ctx context.Context) ([]*model.Dealer, error) {
	query := r.sq.Select("id", "name", "city", "region", "manager", "created_at", "updated_at").
		From(dealerTableName).
		OrderBy("name")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetAll: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetAll: error querying: %w", err)
	}
	defer rows.Close()

	var dealers []*model.Dealer
	for rows.Next() {
		dealer := &model.Dealer{}
		err = rows.Scan(
			&dealer.DealerID, &dealer.DealerNameRu, &dealer.City, &dealer.Region, &dealer.Manager,
			&dealer.CreatedAt, &dealer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DealerRepository.GetAll: error scanning: %w", err)
		}
		dealers = append(dealers, dealer)
	}

	return dealers, nil
}

// GetByRegion получает дилеров по региону.
func (r *DealerRepository) GetByRegion(ctx context.Context, region string) ([]*model.Dealer, error) {
	query := r.sq.Select("id", "name", "city", "region", "manager", "created_at", "updated_at").
		From(dealerTableName).
		Where(squirrel.Eq{"region": region}).
		OrderBy("name")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetByRegion: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetByRegion: error querying: %w", err)
	}
	defer rows.Close()

	var dealers []*model.Dealer
	for rows.Next() {
		dealer := &model.Dealer{}
		err = rows.Scan(
			&dealer.DealerID, &dealer.DealerNameRu, &dealer.City, &dealer.Region, &dealer.Manager,
			&dealer.CreatedAt, &dealer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DealerRepository.GetByRegion: error scanning: %w", err)
		}
		dealers = append(dealers, dealer)
	}

	return dealers, nil
}

// GetWithFilters получает дилеров с применением фильтров.
func (r *DealerRepository) GetWithFilters(ctx context.Context, filters *model.FilterParams) ([]*model.Dealer, error) {
	query := r.sq.Select("id", "name", "city", "region", "manager", "created_at", "updated_at").
		From(dealerTableName)

	// Применяем фильтры
	if filters.HasRegionFilter() {
		query = query.Where(squirrel.Eq{"region": filters.Region})
	}

	if filters.HasDealerFilter() {
		query = query.Where(squirrel.Eq{"id": filters.DealerIDs})
	}

	// Сортировка
	if filters.SortBy != "" {
		order := "ASC"
		if filters.SortOrder == "desc" {
			order = "DESC"
		}
		query = query.OrderBy(fmt.Sprintf("%s %s", filters.SortBy, order))
	} else {
		query = query.OrderBy("name")
	}

	// Пагинация
	if filters.Limit > 0 {
		query = query.Limit(uint64(filters.Limit))
	}
	if filters.Offset > 0 {
		query = query.Offset(uint64(filters.Offset))
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetWithFilters: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetWithFilters: error querying: %w", err)
	}
	defer rows.Close()

	var dealers []*model.Dealer
	for rows.Next() {
		dealer := &model.Dealer{}
		err = rows.Scan(
			&dealer.DealerID, &dealer.DealerNameRu, &dealer.City, &dealer.Region, &dealer.Manager,
			&dealer.CreatedAt, &dealer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DealerRepository.GetWithFilters: error scanning: %w", err)
		}
		dealers = append(dealers, dealer)
	}

	return dealers, nil
}

// Update обновляет данные дилера.
func (r *DealerRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	updates["updated_at"] = time.Now()
	query := r.sq.Update(dealerTableName).SetMap(updates).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.Update: error building query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.Update: error updating: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("DealerRepository.Update: no rows affected, record with id %d not found", id)
	}

	return nil
}

// Delete удаляет дилера.
func (r *DealerRepository) Delete(ctx context.Context, id int) error {
	query := r.sq.Delete(dealerTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.Delete: error building query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.Delete: error deleting: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("DealerRepository.Delete: no rows affected, record with id %d not found", id)
	}

	return nil
}

// GetDealerCardData получает данные карточки дилера за период.
func (r *DealerRepository) GetDealerCardData(ctx context.Context, dealerID int, period time.Time) (*model.DealerCardData, error) {
	// Получаем основную информацию о дилере
	dealer, err := r.GetByID(ctx, dealerID)
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetDealerCardData: error getting dealer: %w", err)
	}

	// Создаем структуру карточки дилера
	cardData := &model.DealerCardData{
		DealerID:      dealer.DealerID,
		Ruft:          dealer.Ruft,
		DealerNameRu:  dealer.DealerNameRu,
		DealerNameEn:  dealer.DealerNameEn,
		City:          dealer.City,
		Region:        dealer.Region,
		Manager:       dealer.Manager,
		JointDecision: dealer.JointDecision,
		Period:        period,
	}

	// Здесь можно добавить логику для получения дополнительных данных из других таблиц
	// Например, данные из dealer_development, sales, aftersales, performance_sales, performance_aftersales

	return cardData, nil
}

// AddBrand добавляет бренд дилеру.
func (r *DealerRepository) AddBrand(ctx context.Context, dealerID int, brandName string) error {
	query := r.sq.Insert("dealer_brands").
		Columns("dealer_id", "brand_name", "created_at").
		Values(dealerID, brandName, time.Now())

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.AddBrand: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.AddBrand: error inserting: %w", err)
	}

	return nil
}

// RemoveBrand удаляет бренд у дилера.
func (r *DealerRepository) RemoveBrand(ctx context.Context, dealerID int, brandName string) error {
	query := r.sq.Delete("dealer_brands").
		Where(squirrel.Eq{"dealer_id": dealerID, "brand_name": brandName})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.RemoveBrand: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.RemoveBrand: error deleting: %w", err)
	}

	return nil
}

// GetBrands получает список брендов дилера.
func (r *DealerRepository) GetBrands(ctx context.Context, dealerID int) ([]string, error) {
	query := r.sq.Select("brand_name").
		From("dealer_brands").
		Where(squirrel.Eq{"dealer_id": dealerID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetBrands: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetBrands: error querying: %w", err)
	}
	defer rows.Close()

	var brands []string
	for rows.Next() {
		var brand string
		err = rows.Scan(&brand)
		if err != nil {
			return nil, fmt.Errorf("DealerRepository.GetBrands: error scanning: %w", err)
		}
		brands = append(brands, brand)
	}

	return brands, nil
}

// AddBusiness добавляет тип бизнеса дилеру.
func (r *DealerRepository) AddBusiness(ctx context.Context, dealerID int, businessType string) error {
	query := r.sq.Insert("dealer_businesses").
		Columns("dealer_id", "business_type", "created_at").
		Values(dealerID, businessType, time.Now())

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.AddBusiness: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.AddBusiness: error inserting: %w", err)
	}

	return nil
}

// RemoveBusiness удаляет тип бизнеса у дилера.
func (r *DealerRepository) RemoveBusiness(ctx context.Context, dealerID int, businessType string) error {
	query := r.sq.Delete("dealer_businesses").
		Where(squirrel.Eq{"dealer_id": dealerID, "business_type": businessType})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.RemoveBusiness: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.RemoveBusiness: error deleting: %w", err)
	}

	return nil
}

// GetBusinesses получает список типов бизнеса дилера.
func (r *DealerRepository) GetBusinesses(ctx context.Context, dealerID int) ([]string, error) {
	query := r.sq.Select("business_type").
		From("dealer_businesses").
		Where(squirrel.Eq{"dealer_id": dealerID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetBusinesses: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetBusinesses: error querying: %w", err)
	}
	defer rows.Close()

	var businesses []string
	for rows.Next() {
		var business string
		err = rows.Scan(&business)
		if err != nil {
			return nil, fmt.Errorf("DealerRepository.GetBusinesses: error scanning: %w", err)
		}
		businesses = append(businesses, business)
	}

	return businesses, nil
}

// UpdateFull обновляет всю запись дилера целиком.
func (r *DealerRepository) UpdateFull(ctx context.Context, dealer *model.Dealer) error {
	dealer.UpdatedAt = time.Now()

	query := r.sq.Update(dealerTableName).
		Set("name", dealer.DealerNameRu).
		Set("region", dealer.Region).
		Set("city", dealer.City).
		Set("manager", dealer.Manager).
		Set("updated_at", dealer.UpdatedAt).
		Where(squirrel.Eq{"id": dealer.DealerID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.UpdateFull: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.UpdateFull: error updating: %w", err)
	}

	return nil
}
