package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

const dealerDevTableName = "dealer_development"

// DealerDevRepository репозиторий для работы с данными развития дилеров.
type DealerDevRepository struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

// NewDealerDevRepository конструктор.
func NewDealerDevRepository(pool *pgxpool.Pool) *DealerDevRepository {
	return &DealerDevRepository{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Create создает новую запись развития дилера.
func (r *DealerDevRepository) Create(ctx context.Context, dd *model.DealerDevelopment) (int, error) {
	now := time.Now()
	dd.CreatedAt = now
	dd.UpdatedAt = now

	query := r.sq.Insert(dealerDevTableName).
		Columns(
			"dealer_id", "period",
			"check_list_score", "dealership_class", "brands", "branding",
			"marketing_investments", "by_side_businesses", "dd_recommendation",
			"created_at", "updated_at",
		).
		Values(
			dd.DealerID, dd.Period,
			dd.CheckListScore, dd.DealershipClass, dd.Brands, dd.Branding,
			dd.MarketingInvestments, dd.BySideBusinesses, dd.DDRecommendation,
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

	dd.ID = int(id)
	return int(id), nil
}

// GetByID получает запись развития дилера по ID.
func (r *DealerDevRepository) GetByID(ctx context.Context, id int) (*model.DealerDevelopment, error) {
	query := r.sq.Select(
		"id", "dealer_id", "period",
		"check_list_score", "dealership_class", "brands", "branding",
		"marketing_investments", "by_side_businesses", "dd_recommendation",
		"created_at", "updated_at",
	).From(dealerDevTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerDevRepository.GetByID: error building query: %w", err)
	}

	dd := &model.DealerDevelopment{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&dd.ID, &dd.DealerID, &dd.Period,
		&dd.CheckListScore, &dd.DealershipClass, &dd.Brands, &dd.Branding,
		&dd.MarketingInvestments, &dd.BySideBusinesses, &dd.DDRecommendation,
		&dd.CreatedAt, &dd.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("DealerDevRepository.GetByID: error scanning: %w", err)
	}

	return dd, nil
}

// GetByDealerAndPeriod получает запись развития дилера по дилеру и периоду.
func (r *DealerDevRepository) GetByDealerAndPeriod(ctx context.Context, dealerID int, period time.Time) (*model.DealerDevelopment, error) {
	query := r.sq.Select(
		"id", "dealer_id", "period",
		"check_list_score", "dealership_class", "brands", "branding",
		"marketing_investments", "by_side_businesses", "dd_recommendation",
		"created_at", "updated_at",
	).From(dealerDevTableName).Where(squirrel.Eq{
		"dealer_id": dealerID,
		"period":    period,
	})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerDevRepository.GetByDealerAndPeriod: error building query: %w", err)
	}

	dd := &model.DealerDevelopment{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&dd.ID, &dd.DealerID, &dd.Period,
		&dd.CheckListScore, &dd.DealershipClass, &dd.Brands, &dd.Branding,
		&dd.MarketingInvestments, &dd.BySideBusinesses, &dd.DDRecommendation,
		&dd.CreatedAt, &dd.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("DealerDevRepository.GetByDealerAndPeriod: error scanning: %w", err)
	}

	return dd, nil
}

// GetAllByPeriod получает все записи развития дилеров за указанный период.
func (r *DealerDevRepository) GetAllByPeriod(ctx context.Context, period time.Time) ([]*model.DealerDevelopment, error) {
	query := r.sq.Select(
		"id", "dealer_id", "period",
		"check_list_score", "dealership_class", "brands", "branding",
		"marketing_investments", "by_side_businesses", "dd_recommendation",
		"created_at", "updated_at",
	).From(dealerDevTableName).Where(squirrel.Eq{"period": period})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("DealerDevRepository.GetAllByPeriod: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("DealerDevRepository.GetAllByPeriod: error querying: %w", err)
	}
	defer rows.Close()

	var dealerDevs []*model.DealerDevelopment
	for rows.Next() {
		dd := &model.DealerDevelopment{}
		err = rows.Scan(
			&dd.ID, &dd.DealerID, &dd.Period,
			&dd.CheckListScore, &dd.DealershipClass, &dd.Brands, &dd.Branding,
			&dd.MarketingInvestments, &dd.BySideBusinesses, &dd.DDRecommendation,
			&dd.CreatedAt, &dd.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DealerDevRepository.GetAllByPeriod: error scanning: %w", err)
		}
		dealerDevs = append(dealerDevs, dd)
	}

	return dealerDevs, nil
}

// UpdateFull обновляет всю запись развития дилера целиком.
func (r *DealerDevRepository) UpdateFull(ctx context.Context, dd *model.DealerDevelopment) error {
	dd.UpdatedAt = time.Now()

	query := r.sq.Update(dealerDevTableName).
		Set("dealer_id", dd.DealerID).
		Set("period", dd.Period).
		Set("check_list_score", dd.CheckListScore).
		Set("dealership_class", dd.DealershipClass).
		Set("brands", dd.Brands).
		Set("branding", dd.Branding).
		Set("marketing_investments", dd.MarketingInvestments).
		Set("by_side_businesses", dd.BySideBusinesses).
		Set("dd_recommendation", dd.DDRecommendation).
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

// GetWithDetailsByPeriod получает записи развития дилеров с деталями за период.
func (r *DealerDevRepository) GetWithDetailsByPeriod(ctx context.Context, period time.Time, region string) ([]*model.DealerDevWithDetails, error) {
	queryBuilder := r.sq.Select(
		"dd.id", "dd.dealer_id", "dd.period",
		"dd.check_list_score", "dd.dealership_class", "dd.brands", "dd.branding",
		"dd.marketing_investments", "dd.by_side_businesses", "dd.dd_recommendation",
		"dd.created_at", "dd.updated_at",
		"d.dealer_name_ru", "d.dealer_name_en", "d.city", "d.region", "d.manager", "d.ruft",
	).
		From(dealerDevTableName + " dd").
		Join("dealers d ON dd.dealer_id = d.dealer_id").
		Where(squirrel.Eq{"dd.period": period})

	if region != "all-russia" {
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
			&ddwd.ID, &ddwd.DealerID, &ddwd.Period,
			&ddwd.CheckListScore, &ddwd.DealershipClass, &ddwd.Brands, &ddwd.Branding,
			&ddwd.MarketingInvestments, &ddwd.BySideBusinesses, &ddwd.DDRecommendation,
			&ddwd.CreatedAt, &ddwd.UpdatedAt,
			&ddwd.DealerNameRu, &ddwd.DealerNameEn, &ddwd.City, &ddwd.Region, &ddwd.Manager, &ddwd.Ruft,
		)
		if err != nil {
			return nil, fmt.Errorf("DealerDevRepository.GetWithDetailsByPeriod: error scanning: %w", err)
		}
		results = append(results, ddwd)
	}

	return results, nil
}

// Delete удаляет запись развития дилера.
func (r *DealerDevRepository) Delete(ctx context.Context, id int) error {
	query := r.sq.Delete(dealerDevTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerDevRepository.Delete: error building query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerDevRepository.Delete: error deleting: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("DealerDevRepository.Delete: no rows affected, record with id %d not found", id)
	}

	return nil
}

// Update обновляет данные развития дилера.
func (r *DealerDevRepository) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	updates["updated_at"] = time.Now()
	query := r.sq.Update(dealerDevTableName).SetMap(updates).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DealerDevRepository.Update: error building query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DealerDevRepository.Update: error updating: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("DealerDevRepository.Update: no rows affected, record with id %d not found", id)
	}

	return nil
}
