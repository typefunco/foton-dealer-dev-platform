package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

const afterSalesTableName = "aftersales"

// AfterSalesRepository репозиторий для работы с данными послепродажного обслуживания.
type AfterSalesRepository struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

// NewAfterSalesRepository конструктор.
func NewAfterSalesRepository(pool *pgxpool.Pool) *AfterSalesRepository {
	return &AfterSalesRepository{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Create создает новую запись послепродажного обслуживания.
func (r *AfterSalesRepository) Create(ctx context.Context, as *model.AfterSales) (int64, error) {
	now := time.Now()
	as.CreatedAt = now
	as.UpdatedAt = now

	query := r.sq.Insert(afterSalesTableName).
		Columns(
			"dealer_id", "period",
			"recommended_stock_pct", "warranty_stock_pct", "foton_labor_hours_pct",
			"warranty_hours", "service_contracts_hours", "as_trainings",
			"spare_parts_sales_q", "spare_parts_sales_ytd_pct", "as_recommendation",
			"created_at", "updated_at",
		).
		Values(
			as.DealerID, as.Period,
			as.RecommendedStockPct, as.WarrantyStockPct, as.FotonLaborHoursPct,
			as.WarrantyHours, as.ServiceContractsHours, as.ASTrainings,
			as.SparePartsSalesQ, as.SparePartsSalesYtdPct, as.ASRecommendation,
			as.CreatedAt, as.UpdatedAt,
		).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("AfterSalesRepository.Create: error building query: %w", err)
	}

	var id int64
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("AfterSalesRepository.Create: error inserting: %w", err)
	}

	as.ID = int(id)
	return id, nil
}

// GetByID получает запись послепродажного обслуживания по ID.
func (r *AfterSalesRepository) GetByID(ctx context.Context, id int64) (*model.AfterSales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "period",
		"recommended_stock_pct", "warranty_stock_pct", "foton_labor_hours_pct",
		"warranty_hours", "service_contracts_hours", "as_trainings",
		"spare_parts_sales_q", "spare_parts_sales_ytd_pct", "as_recommendation",
		"created_at", "updated_at",
	).From(afterSalesTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("AfterSalesRepository.GetByID: error building query: %w", err)
	}

	as := &model.AfterSales{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&as.ID, &as.DealerID, &as.Period,
		&as.RecommendedStockPct, &as.WarrantyStockPct, &as.FotonLaborHoursPct,
		&as.WarrantyHours, &as.ServiceContractsHours, &as.ASTrainings,
		&as.SparePartsSalesQ, &as.SparePartsSalesYtdPct, &as.ASRecommendation,
		&as.CreatedAt, &as.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("AfterSalesRepository.GetByID: error scanning: %w", err)
	}

	return as, nil
}

// GetByDealerAndPeriodTime получает запись послепродажного обслуживания по дилеру и периоду.
func (r *AfterSalesRepository) GetByDealerAndPeriodTime(ctx context.Context, dealerID int, period time.Time) (*model.AfterSales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "period",
		"recommended_stock_pct", "warranty_stock_pct", "foton_labor_hours_pct",
		"warranty_hours", "service_contracts_hours", "as_trainings",
		"spare_parts_sales_q", "spare_parts_sales_ytd_pct", "as_recommendation",
		"created_at", "updated_at",
	).From(afterSalesTableName).Where(squirrel.Eq{
		"dealer_id": dealerID,
		"period":    period,
	})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("AfterSalesRepository.GetByDealerAndPeriod: error building query: %w", err)
	}

	as := &model.AfterSales{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&as.ID, &as.DealerID, &as.Period,
		&as.RecommendedStockPct, &as.WarrantyStockPct, &as.FotonLaborHoursPct,
		&as.WarrantyHours, &as.ServiceContractsHours, &as.ASTrainings,
		&as.SparePartsSalesQ, &as.SparePartsSalesYtdPct, &as.ASRecommendation,
		&as.CreatedAt, &as.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("AfterSalesRepository.GetByDealerAndPeriod: error scanning: %w", err)
	}

	return as, nil
}

// GetAllByPeriodTime получает все записи послепродажного обслуживания за указанный период.
func (r *AfterSalesRepository) GetAllByPeriodTime(ctx context.Context, period time.Time) ([]*model.AfterSales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "period",
		"recommended_stock_pct", "warranty_stock_pct", "foton_labor_hours_pct",
		"warranty_hours", "service_contracts_hours", "as_trainings",
		"spare_parts_sales_q", "spare_parts_sales_ytd_pct", "as_recommendation",
		"created_at", "updated_at",
	).From(afterSalesTableName).Where(squirrel.Eq{"period": period})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("AfterSalesRepository.GetAllByPeriod: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("AfterSalesRepository.GetAllByPeriod: error querying: %w", err)
	}
	defer rows.Close()

	var afterSalesList []*model.AfterSales
	for rows.Next() {
		as := &model.AfterSales{}
		err = rows.Scan(
			&as.ID, &as.DealerID, &as.Period,
			&as.RecommendedStockPct, &as.WarrantyStockPct, &as.FotonLaborHoursPct,
			&as.WarrantyHours, &as.ServiceContractsHours, &as.ASTrainings,
			&as.SparePartsSalesQ, &as.SparePartsSalesYtdPct, &as.ASRecommendation,
			&as.CreatedAt, &as.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("AfterSalesRepository.GetAllByPeriod: error scanning: %w", err)
		}
		afterSalesList = append(afterSalesList, as)
	}

	return afterSalesList, nil
}

// UpdateFull обновляет всю запись послепродажного обслуживания целиком.
func (r *AfterSalesRepository) UpdateFull(ctx context.Context, as *model.AfterSales) error {
	as.UpdatedAt = time.Now()

	query := r.sq.Update(afterSalesTableName).
		Set("dealer_id", as.DealerID).
		Set("period", as.Period).
		Set("recommended_stock_pct", as.RecommendedStockPct).
		Set("warranty_stock_pct", as.WarrantyStockPct).
		Set("foton_labor_hours_pct", as.FotonLaborHoursPct).
		Set("warranty_hours", as.WarrantyHours).
		Set("service_contracts_hours", as.ServiceContractsHours).
		Set("as_trainings", as.ASTrainings).
		Set("spare_parts_sales_q", as.SparePartsSalesQ).
		Set("spare_parts_sales_ytd_pct", as.SparePartsSalesYtdPct).
		Set("as_recommendation", as.ASRecommendation).
		Set("updated_at", as.UpdatedAt).
		Where(squirrel.Eq{"id": as.ID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("AfterSalesRepository.UpdateFull: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AfterSalesRepository.UpdateFull: error updating: %w", err)
	}

	return nil
}

// GetWithDetailsByPeriodTime получает записи послепродажного обслуживания с деталями за период.
func (r *AfterSalesRepository) GetWithDetailsByPeriodTime(ctx context.Context, period time.Time, region string) ([]*model.AfterSalesWithDetails, error) {

	queryBuilder := r.sq.Select(
		"aftersales.id", "aftersales.dealer_id", "aftersales.period",
		"aftersales.recommended_stock_pct", "aftersales.warranty_stock_pct", "aftersales.foton_labor_hours_pct",
		"aftersales.warranty_hours", "aftersales.service_contracts_hours", "aftersales.as_trainings",
		"aftersales.spare_parts_sales_q", "aftersales.spare_parts_sales_ytd_pct", "aftersales.as_recommendation",
		"aftersales.created_at", "aftersales.updated_at",
		"d.dealer_name_ru", "d.dealer_name_en", "d.city", "d.region", "d.manager", "d.ruft",
	).
		From(afterSalesTableName + " aftersales").
		Join("dealers d ON aftersales.dealer_id = d.dealer_id").
		Where(squirrel.Eq{"aftersales.period": period})

	if region != "all-russia" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"d.region": region})
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("AfterSalesRepository.GetWithDetailsByPeriod: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("AfterSalesRepository.GetWithDetailsByPeriod: error querying: %w", err)
	}
	defer rows.Close()

	var results []*model.AfterSalesWithDetails
	for rows.Next() {
		aswd := &model.AfterSalesWithDetails{}
		err = rows.Scan(
			&aswd.ID, &aswd.DealerID, &aswd.Period,
			&aswd.RecommendedStockPct, &aswd.WarrantyStockPct, &aswd.FotonLaborHoursPct,
			&aswd.WarrantyHours, &aswd.ServiceContractsHours, &aswd.ASTrainings,
			&aswd.SparePartsSalesQ, &aswd.SparePartsSalesYtdPct, &aswd.ASRecommendation,
			&aswd.CreatedAt, &aswd.UpdatedAt,
			&aswd.DealerNameRu, &aswd.DealerNameEn, &aswd.City, &aswd.Region, &aswd.Manager, &aswd.Ruft,
		)
		if err != nil {
			return nil, fmt.Errorf("AfterSalesRepository.GetWithDetailsByPeriod: error scanning: %w", err)
		}
		results = append(results, aswd)
	}

	return results, nil
}

// Delete удаляет запись послепродажного обслуживания.
func (r *AfterSalesRepository) Delete(ctx context.Context, id int64) error {
	query := r.sq.Delete(afterSalesTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("AfterSalesRepository.Delete: error building query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AfterSalesRepository.Delete: error deleting: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("AfterSalesRepository.Delete: no rows affected, record with id %d not found", id)
	}

	return nil
}

// Update обновляет данные послепродажного обслуживания.
func (r *AfterSalesRepository) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	updates["updated_at"] = time.Now()
	query := r.sq.Update(afterSalesTableName).SetMap(updates).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("AfterSalesRepository.Update: error building query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AfterSalesRepository.Update: error updating: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("AfterSalesRepository.Update: no rows affected, record with id %d not found", id)
	}

	return nil
}

// GetByDealerAndPeriod получает запись послепродажного обслуживания по дилеру и периоду (с кварталом и годом).
func (r *AfterSalesRepository) GetByDealerAndPeriod(ctx context.Context, dealerID int64, quarter string, year int) (*model.AfterSales, error) {
	// Преобразуем quarter/year в period
	var month int
	switch quarter {
	case "q1":
		month = 1
	case "q2":
		month = 4
	case "q3":
		month = 7
	case "q4":
		month = 10
	default:
		return nil, fmt.Errorf("invalid quarter: %s", quarter)
	}
	period := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	return r.GetByDealerAndPeriodTime(ctx, int(dealerID), period)
}

// GetAllByPeriod получает все записи послепродажного обслуживания за указанный период (с кварталом и годом).
func (r *AfterSalesRepository) GetAllByPeriod(ctx context.Context, quarter string, year int) ([]*model.AfterSales, error) {
	// Преобразуем quarter/year в period
	var month int
	switch quarter {
	case "q1":
		month = 1
	case "q2":
		month = 4
	case "q3":
		month = 7
	case "q4":
		month = 10
	default:
		return nil, fmt.Errorf("invalid quarter: %s", quarter)
	}
	period := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	return r.GetAllByPeriodTime(ctx, period)
}

// GetWithDetailsByPeriod получает записи послепродажного обслуживания с деталями за период (с кварталом и годом).
func (r *AfterSalesRepository) GetWithDetailsByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.AfterSalesWithDetails, error) {
	// Преобразуем quarter/year в period
	var month int
	switch quarter {
	case "q1":
		month = 1
	case "q2":
		month = 4
	case "q3":
		month = 7
	case "q4":
		month = 10
	default:
		return nil, fmt.Errorf("invalid quarter: %s", quarter)
	}
	period := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	return r.GetWithDetailsByPeriodTime(ctx, period, region)
}
