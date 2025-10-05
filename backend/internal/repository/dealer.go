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

// NewDealerRepository создает новый экземпляр репозитория.
func NewDealerRepository(pool *pgxpool.Pool) *DealerRepository {
	return &DealerRepository{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Create создает нового дилера.
func (r *DealerRepository) Create(ctx context.Context, dealer *model.Dealer) (int64, error) {
	now := time.Now()
	dealer.CreatedAt = now
	dealer.UpdatedAt = now

	query := r.sq.Insert(dealerTableName).
		Columns("name", "city", "region", "manager", "created_at", "updated_at").
		Values(dealer.Name, dealer.City, dealer.Region, dealer.Manager, dealer.CreatedAt, dealer.UpdatedAt).
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

	dealer.ID = id
	return id, nil
}

// GetByID получает дилера по ID.
func (r *DealerRepository) GetByID(ctx context.Context, id int64) (*model.Dealer, error) {
	query := r.sq.Select("id", "name", "city", "region", "manager", "created_at", "updated_at").
		From(dealerTableName).
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerRepository.GetByID: error building query: %w", err)
	}

	dealer := &model.Dealer{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&dealer.ID, &dealer.Name, &dealer.City, &dealer.Region, &dealer.Manager,
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
		OrderBy("name ASC")

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
			&dealer.ID, &dealer.Name, &dealer.City, &dealer.Region, &dealer.Manager,
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
	queryBuilder := r.sq.Select("id", "name", "city", "region", "manager", "created_at", "updated_at").
		From(dealerTableName).
		OrderBy("name ASC")

	// Если регион "all-russia", получаем всех дилеров
	if region != "" && region != "all-russia" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"region": region})
	}

	sql, args, err := queryBuilder.ToSql()
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
			&dealer.ID, &dealer.Name, &dealer.City, &dealer.Region, &dealer.Manager,
			&dealer.CreatedAt, &dealer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DealerRepository.GetByRegion: error scanning: %w", err)
		}
		dealers = append(dealers, dealer)
	}

	return dealers, nil
}

// Update обновляет дилера (частичное обновление).
func (r *DealerRepository) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("DealerRepository.Update: no fields to update")
	}

	updates["updated_at"] = time.Now()

	query := r.sq.Update(dealerTableName).
		Where(squirrel.Eq{"id": id}).
		SetMap(updates)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.Update: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.Update: error updating: %w", err)
	}

	return nil
}

// UpdateFull обновляет дилера полностью.
func (r *DealerRepository) UpdateFull(ctx context.Context, dealer *model.Dealer) error {
	dealer.UpdatedAt = time.Now()

	query := r.sq.Update(dealerTableName).
		Set("name", dealer.Name).
		Set("city", dealer.City).
		Set("region", dealer.Region).
		Set("manager", dealer.Manager).
		Set("updated_at", dealer.UpdatedAt).
		Where(squirrel.Eq{"id": dealer.ID})

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

// Delete удаляет дилера по ID.
func (r *DealerRepository) Delete(ctx context.Context, id int64) error {
	query := r.sq.Delete(dealerTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerRepository.Delete: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerRepository.Delete: error deleting: %w", err)
	}

	return nil
}

// AddBrand добавляет бренд к дилеру.
func (r *DealerRepository) AddBrand(ctx context.Context, dealerID int64, brandName string) error {
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
func (r *DealerRepository) RemoveBrand(ctx context.Context, dealerID int64, brandName string) error {
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
func (r *DealerRepository) GetBrands(ctx context.Context, dealerID int64) ([]string, error) {
	query := r.sq.Select("brand_name").
		From("dealer_brands").
		Where(squirrel.Eq{"dealer_id": dealerID}).
		OrderBy("brand_name ASC")

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

// AddBusiness добавляет побочный бизнес к дилеру.
func (r *DealerRepository) AddBusiness(ctx context.Context, dealerID int64, businessType string) error {
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

// RemoveBusiness удаляет побочный бизнес у дилера.
func (r *DealerRepository) RemoveBusiness(ctx context.Context, dealerID int64, businessType string) error {
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

// GetBusinesses получает список побочных бизнесов дилера.
func (r *DealerRepository) GetBusinesses(ctx context.Context, dealerID int64) ([]string, error) {
	query := r.sq.Select("business_type").
		From("dealer_businesses").
		Where(squirrel.Eq{"dealer_id": dealerID}).
		OrderBy("business_type ASC")

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
