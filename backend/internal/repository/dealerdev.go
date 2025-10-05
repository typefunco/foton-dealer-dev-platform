package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

const dealerDevTableName = "dealer_dev"

// DealerDevRepository репозиторий для работы с данными развития дилеров.
type DealerDevRepository struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

// NewDealerDevRepository создает новый экземпляр репозитория.
func NewDealerDevRepository(pool *pgxpool.Pool) *DealerDevRepository {
	return &DealerDevRepository{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Create создает новую запись развития дилера.
func (r *DealerDevRepository) Create(ctx context.Context, dd *model.DealerDev) (int64, error) {
	now := time.Now()
	dd.CreatedAt = now
	dd.UpdatedAt = now

	query := r.sq.Insert(dealerDevTableName).
		Columns(
			"dealer_id", "quarter", "year",
			"check_list_score", "dealer_ship_class", "branding",
			"marketing_investments", "dealer_dev_recommendation",
			"created_at", "updated_at",
		).
		Values(
			dd.DealerID, dd.Quarter, dd.Year,
			dd.CheckListScore, dd.DealerShipClass, dd.Branding,
			dd.MarketingInvestments, dd.Recommendation,
			dd.CreatedAt, dd.UpdatedAt,
		).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("DealerDevRepository.Create: error building query: %w", err)
	}

	var id int64
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("DealerDevRepository.Create: error inserting: %w", err)
	}

	dd.ID = id
	return id, nil
}

// GetByID получает запись развития дилера по ID.
func (r *DealerDevRepository) GetByID(ctx context.Context, id int64) (*model.DealerDev, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"check_list_score", "dealer_ship_class", "branding",
		"marketing_investments", "dealer_dev_recommendation",
		"created_at", "updated_at",
	).From(dealerDevTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerDevRepository.GetByID: error building query: %w", err)
	}

	dd := &model.DealerDev{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&dd.ID, &dd.DealerID, &dd.Quarter, &dd.Year,
		&dd.CheckListScore, &dd.DealerShipClass, &dd.Branding,
		&dd.MarketingInvestments, &dd.Recommendation,
		&dd.CreatedAt, &dd.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("DealerDevRepository.GetByID: error scanning: %w", err)
	}

	return dd, nil
}

// GetByDealerAndPeriod получает запись развития дилера по дилеру и периоду.
func (r *DealerDevRepository) GetByDealerAndPeriod(ctx context.Context, dealerID int64, quarter string, year int) (*model.DealerDev, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"check_list_score", "dealer_ship_class", "branding",
		"marketing_investments", "dealer_dev_recommendation",
		"created_at", "updated_at",
	).From(dealerDevTableName).Where(squirrel.Eq{
		"dealer_id": dealerID,
		"quarter":   quarter,
		"year":      year,
	})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerDevRepository.GetByDealerAndPeriod: error building query: %w", err)
	}

	dd := &model.DealerDev{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&dd.ID, &dd.DealerID, &dd.Quarter, &dd.Year,
		&dd.CheckListScore, &dd.DealerShipClass, &dd.Branding,
		&dd.MarketingInvestments, &dd.Recommendation,
		&dd.CreatedAt, &dd.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("DealerDevRepository.GetByDealerAndPeriod: error scanning: %w", err)
	}

	return dd, nil
}

// GetAllByPeriod получает все записи развития дилеров за указанный период.
func (r *DealerDevRepository) GetAllByPeriod(ctx context.Context, quarter string, year int) ([]*model.DealerDev, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"check_list_score", "dealer_ship_class", "branding",
		"marketing_investments", "dealer_dev_recommendation",
		"created_at", "updated_at",
	).From(dealerDevTableName).Where(squirrel.Eq{
		"quarter": quarter,
		"year":    year,
	}).OrderBy("dealer_ship_class ASC", "check_list_score DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerDevRepository.GetAllByPeriod: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("DealerDevRepository.GetAllByPeriod: error querying: %w", err)
	}
	defer rows.Close()

	var dealerDevs []*model.DealerDev
	for rows.Next() {
		dd := &model.DealerDev{}
		err = rows.Scan(
			&dd.ID, &dd.DealerID, &dd.Quarter, &dd.Year,
			&dd.CheckListScore, &dd.DealerShipClass, &dd.Branding,
			&dd.MarketingInvestments, &dd.Recommendation,
			&dd.CreatedAt, &dd.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DealerDevRepository.GetAllByPeriod: error scanning: %w", err)
		}
		dealerDevs = append(dealerDevs, dd)
	}

	return dealerDevs, nil
}

// Update обновляет запись развития дилера (частичное обновление).
func (r *DealerDevRepository) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("DealerDevRepository.Update: no fields to update")
	}

	updates["updated_at"] = time.Now()

	query := r.sq.Update(dealerDevTableName).
		Where(squirrel.Eq{"id": id}).
		SetMap(updates)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerDevRepository.Update: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerDevRepository.Update: error updating: %w", err)
	}

	return nil
}

// UpdateFull обновляет всю запись развития дилера.
func (r *DealerDevRepository) UpdateFull(ctx context.Context, dd *model.DealerDev) error {
	dd.UpdatedAt = time.Now()

	query := r.sq.Update(dealerDevTableName).
		Set("dealer_id", dd.DealerID).
		Set("quarter", dd.Quarter).
		Set("year", dd.Year).
		Set("check_list_score", dd.CheckListScore).
		Set("dealer_ship_class", dd.DealerShipClass).
		Set("branding", dd.Branding).
		Set("marketing_investments", dd.MarketingInvestments).
		Set("dealer_dev_recommendation", dd.Recommendation).
		Set("updated_at", dd.UpdatedAt).
		Where(squirrel.Eq{"id": dd.ID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerDevRepository.UpdateFull: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerDevRepository.UpdateFull: error updating: %w", err)
	}

	return nil
}

// Delete удаляет запись развития дилера по ID.
func (r *DealerDevRepository) Delete(ctx context.Context, id int64) error {
	query := r.sq.Delete(dealerDevTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerDevRepository.Delete: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerDevRepository.Delete: error deleting: %w", err)
	}

	return nil
}

// GetWithDetailsByPeriod получает записи развития дилеров с деталями за период.
func (r *DealerDevRepository) GetWithDetailsByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.DealerDevWithDetails, error) {
	dealerRepo := NewDealerRepository(r.pool)

	queryBuilder := r.sq.Select(
		"dd.id", "dd.dealer_id", "dd.quarter", "dd.year",
		"dd.check_list_score", "dd.dealer_ship_class", "dd.branding",
		"dd.marketing_investments", "dd.dealer_dev_recommendation",
		"dd.created_at", "dd.updated_at",
		"d.name as dealer_name", "d.city", "d.region", "d.manager",
	).From(dealerDevTableName+" dd").
		Join("dealers d ON dd.dealer_id = d.id").
		Where(squirrel.Eq{
			"dd.quarter": quarter,
			"dd.year":    year,
		}).OrderBy("dd.dealer_ship_class ASC", "dd.check_list_score DESC")

	if region != "" && region != "all-russia" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"d.region": region})
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerDevRepository.GetWithDetailsByPeriod: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("DealerDevRepository.GetWithDetailsByPeriod: error querying: %w", err)
	}
	defer rows.Close()

	var results []*model.DealerDevWithDetails
	for rows.Next() {
		ddwd := &model.DealerDevWithDetails{}
		err = rows.Scan(
			&ddwd.ID, &ddwd.DealerID, &ddwd.Quarter, &ddwd.Year,
			&ddwd.CheckListScore, &ddwd.DealerShipClass, &ddwd.Branding,
			&ddwd.MarketingInvestments, &ddwd.Recommendation,
			&ddwd.CreatedAt, &ddwd.UpdatedAt,
			&ddwd.DealerName, &ddwd.City, &ddwd.Region, &ddwd.Manager,
		)
		if err != nil {
			return nil, fmt.Errorf("DealerDevRepository.GetWithDetailsByPeriod: error scanning: %w", err)
		}

		// Получаем бренды дилера
		brands, err := dealerRepo.GetBrands(ctx, ddwd.DealerID)
		if err != nil {
			// Логируем ошибку, но продолжаем
			brands = []string{}
		}
		ddwd.Brands = brands
		ddwd.BrandsCount = len(brands)

		// Получаем побочные бизнесы
		businesses, err := dealerRepo.GetBusinesses(ctx, ddwd.DealerID)
		if err != nil {
			businesses = []string{}
		}
		ddwd.BySideBusinesses = businesses

		results = append(results, ddwd)
	}

	return results, nil
}
