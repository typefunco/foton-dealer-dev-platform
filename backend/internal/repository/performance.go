package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

const performanceTableName = "performance_sales"

// PerformanceRepository репозиторий для работы с данными производительности.
type PerformanceRepository struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

// NewPerformanceRepository конструктор.
func NewPerformanceRepository(pool *pgxpool.Pool) *PerformanceRepository {
	return &PerformanceRepository{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Create создает новую запись производительности продаж.
func (r *PerformanceRepository) Create(ctx context.Context, perf *model.PerformanceSales) (int64, error) {
	now := time.Now()
	perf.CreatedAt = now
	perf.UpdatedAt = now

	query := r.sq.Insert(performanceTableName).
		Columns(
			"dealer_id", "period",
			"quantity_sold", "sales_revenue", "sales_revenue_no_vat", "sales_cost",
			"sales_margin", "sales_margin_pct", "sales_profit_pct",
			"created_at", "updated_at",
		).
		Values(
			perf.DealerID, perf.Period,
			perf.QuantitySold, perf.SalesRevenue, perf.SalesRevenueNoVat, perf.SalesCost,
			perf.SalesMargin, perf.SalesMarginPct, perf.SalesProfitPct,
			perf.CreatedAt, perf.UpdatedAt,
		).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("PerformanceRepository.Create: error building query: %w", err)
	}

	var id int64
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("PerformanceRepository.Create: error inserting: %w", err)
	}

	perf.ID = int(id)
	return id, nil
}

// GetByID получает запись производительности продаж по ID.
func (r *PerformanceRepository) GetByID(ctx context.Context, id int64) (*model.PerformanceSales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "period",
		"quantity_sold", "sales_revenue", "sales_revenue_no_vat", "sales_cost",
		"sales_margin", "sales_margin_pct", "sales_profit_pct",
		"created_at", "updated_at",
	).From(performanceTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.GetByID: error building query: %w", err)
	}

	perf := &model.PerformanceSales{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&perf.ID, &perf.DealerID, &perf.Period,
		&perf.QuantitySold, &perf.SalesRevenue, &perf.SalesRevenueNoVat, &perf.SalesCost,
		&perf.SalesMargin, &perf.SalesMarginPct, &perf.SalesProfitPct,
		&perf.CreatedAt, &perf.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.GetByID: error scanning: %w", err)
	}

	return perf, nil
}

// GetByDealerAndPeriod получает запись производительности продаж по дилеру и периоду.
func (r *PerformanceRepository) GetByDealerAndPeriod(ctx context.Context, dealerID int64, quarter string, year int) (*model.PerformanceSales, error) {
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

	query := r.sq.Select(
		"id", "dealer_id", "period",
		"quantity_sold", "sales_revenue", "sales_revenue_no_vat", "sales_cost",
		"sales_margin", "sales_margin_pct", "sales_profit_pct",
		"created_at", "updated_at",
	).From(performanceTableName).Where(squirrel.Eq{
		"dealer_id": dealerID,
		"period":    period,
	})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.GetByDealerAndPeriod: error building query: %w", err)
	}

	perf := &model.PerformanceSales{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&perf.ID, &perf.DealerID, &perf.Period,
		&perf.QuantitySold, &perf.SalesRevenue, &perf.SalesRevenueNoVat, &perf.SalesCost,
		&perf.SalesMargin, &perf.SalesMarginPct, &perf.SalesProfitPct,
		&perf.CreatedAt, &perf.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.GetByDealerAndPeriod: error scanning: %w", err)
	}

	return perf, nil
}

// GetAllByPeriod получает все записи производительности продаж за указанный период.
func (r *PerformanceRepository) GetAllByPeriod(ctx context.Context, quarter string, year int) ([]*model.PerformanceSales, error) {
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

	query := r.sq.Select(
		"id", "dealer_id", "period",
		"quantity_sold", "sales_revenue", "sales_revenue_no_vat", "sales_cost",
		"sales_margin", "sales_margin_pct", "sales_profit_pct",
		"created_at", "updated_at",
	).From(performanceTableName).Where(squirrel.Eq{"period": period})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.GetAllByPeriod: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.GetAllByPeriod: error querying: %w", err)
	}
	defer rows.Close()

	var performances []*model.PerformanceSales
	for rows.Next() {
		perf := &model.PerformanceSales{}
		err = rows.Scan(
			&perf.ID, &perf.DealerID, &perf.Period,
			&perf.QuantitySold, &perf.SalesRevenue, &perf.SalesRevenueNoVat, &perf.SalesCost,
			&perf.SalesMargin, &perf.SalesMarginPct, &perf.SalesProfitPct,
			&perf.CreatedAt, &perf.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("PerformanceRepository.GetAllByPeriod: error scanning: %w", err)
		}
		performances = append(performances, perf)
	}

	return performances, nil
}

// UpdateFull обновляет всю запись производительности продаж целиком.
func (r *PerformanceRepository) UpdateFull(ctx context.Context, perf *model.PerformanceSales) error {
	perf.UpdatedAt = time.Now()

	query := r.sq.Update(performanceTableName).
		Set("dealer_id", perf.DealerID).
		Set("period", perf.Period).
		Set("quantity_sold", perf.QuantitySold).
		Set("sales_revenue", perf.SalesRevenue).
		Set("sales_revenue_no_vat", perf.SalesRevenueNoVat).
		Set("sales_cost", perf.SalesCost).
		Set("sales_margin", perf.SalesMargin).
		Set("sales_margin_pct", perf.SalesMarginPct).
		Set("sales_profit_pct", perf.SalesProfitPct).
		Set("updated_at", perf.UpdatedAt).
		Where(squirrel.Eq{"id": perf.ID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("PerformanceRepository.UpdateFull: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PerformanceRepository.UpdateFull: error updating: %w", err)
	}

	return nil
}

// GetWithDetailsByPeriod получает записи производительности продаж с деталями за период.
func (r *PerformanceRepository) GetWithDetailsByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.PerformanceWithDetails, error) {
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

	queryBuilder := r.sq.Select(
		"ps.id", "ps.dealer_id", "ps.period",
		"ps.quantity_sold", "ps.sales_revenue", "ps.sales_revenue_no_vat", "ps.sales_cost",
		"ps.sales_margin", "ps.sales_margin_pct", "ps.sales_profit_pct",
		"ps.created_at", "ps.updated_at",
		"d.dealer_name_ru", "d.city", "d.region", "d.manager",
	).
		From(performanceTableName + " ps").
		Join("dealers d ON ps.dealer_id = d.dealer_id").
		Where(squirrel.Eq{"ps.period": period})

	if region != "all-russia" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"d.region": region})
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.GetWithDetailsByPeriod: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.GetWithDetailsByPeriod: error querying: %w", err)
	}
	defer rows.Close()

	var results []*model.PerformanceWithDetails
	for rows.Next() {
		pwd := &model.PerformanceWithDetails{}
		err = rows.Scan(
			&pwd.DealerID, &pwd.DealerName, &pwd.City, &pwd.Region, &pwd.Manager,
			&pwd.FotonRank, &pwd.SalesRevenue, &pwd.SalesProfit, &pwd.SalesMargin,
			&pwd.AsRevenue, &pwd.AsProfit, &pwd.AsMargin, &pwd.Marketing, &pwd.Decision,
		)
		if err != nil {
			return nil, fmt.Errorf("PerformanceRepository.GetWithDetailsByPeriod: error scanning: %w", err)
		}
		results = append(results, pwd)
	}

	return results, nil
}

// Delete удаляет запись производительности продаж.
func (r *PerformanceRepository) Delete(ctx context.Context, id int64) error {
	query := r.sq.Delete(performanceTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("PerformanceRepository.Delete: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PerformanceRepository.Delete: error deleting: %w", err)
	}

	return nil
}

// Update обновляет данные производительности продаж.
func (r *PerformanceRepository) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	updates["updated_at"] = time.Now()
	query := r.sq.Update(performanceTableName).SetMap(updates).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("PerformanceRepository.Update: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PerformanceRepository.Update: error updating: %w", err)
	}

	return nil
}

// FindPerformances получает записи производительности продаж по региону (для обратной совместимости).
// Deprecated: Используйте GetWithDetailsByPeriod для получения данных с деталями.
func (r *PerformanceRepository) FindPerformances(ctx context.Context, region string) ([]*model.PerformanceSales, error) {
	// Получаем данные за текущий квартал (можно настроить)
	// Для примера используем q1 2025
	return r.GetAllByPeriod(ctx, "q1", 2025)
}
