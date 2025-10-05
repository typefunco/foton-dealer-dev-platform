package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

const performanceTableName = "performance"

// PerformanceRepository репозиторий для работы с данными производительности дилеров.
type PerformanceRepository struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

// NewPerformanceRepository создает новый экземпляр репозитория.
func NewPerformanceRepository(pool *pgxpool.Pool) *PerformanceRepository {
	return &PerformanceRepository{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Create создает новую запись производительности.
func (r *PerformanceRepository) Create(ctx context.Context, perf *model.Performance) (int64, error) {
	now := time.Now()
	perf.CreatedAt = now
	perf.UpdatedAt = now

	query := r.sq.Insert(performanceTableName).
		Columns(
			"dealer_id", "quarter", "year",
			"sales_revenue_rub", "sales_profit_rub", "sales_profit_percent", "sales_margin_percent",
			"after_sales_revenue_rub", "after_sales_profit_rub", "after_sales_margin_percent",
			"marketing_investment", "foton_rank", "performance_decision",
			"created_at", "updated_at",
		).
		Values(
			perf.DealerID, perf.Quarter, perf.Year,
			perf.SalesRevenueRub, perf.SalesProfitRub, perf.SalesProfitPercent, perf.SalesMarginPercent,
			perf.AfterSalesRevenueRub, perf.AfterSalesProfitRub, perf.AfterSalesMarginPercent,
			perf.MarketingInvestment, perf.FotonRank, perf.PerformanceDecision,
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

	perf.ID = id
	return id, nil
}

// GetByID получает запись производительности по ID.
func (r *PerformanceRepository) GetByID(ctx context.Context, id int64) (*model.Performance, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"sales_revenue_rub", "sales_profit_rub", "sales_profit_percent", "sales_margin_percent",
		"after_sales_revenue_rub", "after_sales_profit_rub", "after_sales_margin_percent",
		"marketing_investment", "foton_rank", "performance_decision",
		"created_at", "updated_at",
	).From(performanceTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.GetByID: error building query: %w", err)
	}

	perf := &model.Performance{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&perf.ID, &perf.DealerID, &perf.Quarter, &perf.Year,
		&perf.SalesRevenueRub, &perf.SalesProfitRub, &perf.SalesProfitPercent, &perf.SalesMarginPercent,
		&perf.AfterSalesRevenueRub, &perf.AfterSalesProfitRub, &perf.AfterSalesMarginPercent,
		&perf.MarketingInvestment, &perf.FotonRank, &perf.PerformanceDecision,
		&perf.CreatedAt, &perf.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.GetByID: error scanning: %w", err)
	}

	return perf, nil
}

// GetByDealerAndPeriod получает запись производительности по дилеру и периоду.
func (r *PerformanceRepository) GetByDealerAndPeriod(ctx context.Context, dealerID int64, quarter string, year int) (*model.Performance, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"sales_revenue_rub", "sales_profit_rub", "sales_profit_percent", "sales_margin_percent",
		"after_sales_revenue_rub", "after_sales_profit_rub", "after_sales_margin_percent",
		"marketing_investment", "foton_rank", "performance_decision",
		"created_at", "updated_at",
	).From(performanceTableName).Where(squirrel.Eq{
		"dealer_id": dealerID,
		"quarter":   quarter,
		"year":      year,
	})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.GetByDealerAndPeriod: error building query: %w", err)
	}

	perf := &model.Performance{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&perf.ID, &perf.DealerID, &perf.Quarter, &perf.Year,
		&perf.SalesRevenueRub, &perf.SalesProfitRub, &perf.SalesProfitPercent, &perf.SalesMarginPercent,
		&perf.AfterSalesRevenueRub, &perf.AfterSalesProfitRub, &perf.AfterSalesMarginPercent,
		&perf.MarketingInvestment, &perf.FotonRank, &perf.PerformanceDecision,
		&perf.CreatedAt, &perf.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.GetByDealerAndPeriod: error scanning: %w", err)
	}

	return perf, nil
}

// GetAllByPeriod получает все записи производительности за указанный период.
func (r *PerformanceRepository) GetAllByPeriod(ctx context.Context, quarter string, year int) ([]*model.Performance, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"sales_revenue_rub", "sales_profit_rub", "sales_profit_percent", "sales_margin_percent",
		"after_sales_revenue_rub", "after_sales_profit_rub", "after_sales_margin_percent",
		"marketing_investment", "foton_rank", "performance_decision",
		"created_at", "updated_at",
	).From(performanceTableName).Where(squirrel.Eq{
		"quarter": quarter,
		"year":    year,
	}).OrderBy("foton_rank ASC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.GetAllByPeriod: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.GetAllByPeriod: error querying: %w", err)
	}
	defer rows.Close()

	var performances []*model.Performance
	for rows.Next() {
		perf := &model.Performance{}
		err = rows.Scan(
			&perf.ID, &perf.DealerID, &perf.Quarter, &perf.Year,
			&perf.SalesRevenueRub, &perf.SalesProfitRub, &perf.SalesProfitPercent, &perf.SalesMarginPercent,
			&perf.AfterSalesRevenueRub, &perf.AfterSalesProfitRub, &perf.AfterSalesMarginPercent,
			&perf.MarketingInvestment, &perf.FotonRank, &perf.PerformanceDecision,
			&perf.CreatedAt, &perf.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("PerformanceRepository.GetAllByPeriod: error scanning: %w", err)
		}
		performances = append(performances, perf)
	}

	return performances, nil
}

// Update обновляет запись производительности.
// Принимает карту полей для обновления, что позволяет обновлять как одно поле, так и всю модель.
func (r *PerformanceRepository) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("PerformanceRepository.Update: no fields to update")
	}

	// Автоматически обновляем updated_at
	updates["updated_at"] = time.Now()

	query := r.sq.Update(performanceTableName).
		Where(squirrel.Eq{"id": id}).
		SetMap(updates)

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

// UpdateFull обновляет всю запись производительности целиком.
func (r *PerformanceRepository) UpdateFull(ctx context.Context, perf *model.Performance) error {
	perf.UpdatedAt = time.Now()

	query := r.sq.Update(performanceTableName).
		Set("dealer_id", perf.DealerID).
		Set("quarter", perf.Quarter).
		Set("year", perf.Year).
		Set("sales_revenue_rub", perf.SalesRevenueRub).
		Set("sales_profit_rub", perf.SalesProfitRub).
		Set("sales_profit_percent", perf.SalesProfitPercent).
		Set("sales_margin_percent", perf.SalesMarginPercent).
		Set("after_sales_revenue_rub", perf.AfterSalesRevenueRub).
		Set("after_sales_profit_rub", perf.AfterSalesProfitRub).
		Set("after_sales_margin_percent", perf.AfterSalesMarginPercent).
		Set("marketing_investment", perf.MarketingInvestment).
		Set("foton_rank", perf.FotonRank).
		Set("performance_decision", perf.PerformanceDecision).
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

// Delete удаляет запись производительности по ID.
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

// GetWithDetailsByPeriod получает записи производительности с деталями дилера за период.
func (r *PerformanceRepository) GetWithDetailsByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.PerformanceWithDetails, error) {
	queryBuilder := r.sq.Select(
		"p.id", "p.dealer_id", "p.quarter", "p.year",
		"p.sales_revenue_rub", "p.sales_profit_rub", "p.sales_profit_percent", "p.sales_margin_percent",
		"p.after_sales_revenue_rub", "p.after_sales_profit_rub", "p.after_sales_margin_percent",
		"p.marketing_investment", "p.foton_rank", "p.performance_decision",
		"p.created_at", "p.updated_at",
		"d.name as dealer_name", "d.city", "d.region", "d.manager",
	).From(performanceTableName + " p").
		Join("dealers d ON p.dealer_id = d.id").
		Where(squirrel.Eq{
			"p.quarter": quarter,
			"p.year":    year,
		}).OrderBy("p.foton_rank ASC")

	// Если указан конкретный регион, добавляем фильтр
	if region != "" && region != "all-russia" {
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
			&pwd.ID, &pwd.DealerID, &pwd.Quarter, &pwd.Year,
			&pwd.SalesRevenueRub, &pwd.SalesProfitRub, &pwd.SalesProfitPercent, &pwd.SalesMarginPercent,
			&pwd.AfterSalesRevenueRub, &pwd.AfterSalesProfitRub, &pwd.AfterSalesMarginPercent,
			&pwd.MarketingInvestment, &pwd.FotonRank, &pwd.PerformanceDecision,
			&pwd.CreatedAt, &pwd.UpdatedAt,
			&pwd.DealerName, &pwd.City, &pwd.Region, &pwd.Manager,
		)
		if err != nil {
			return nil, fmt.Errorf("PerformanceRepository.GetWithDetailsByPeriod: error scanning: %w", err)
		}

		// Форматируем денежные значения
		pwd.SalesRevenueFormatted = formatMoney(pwd.SalesRevenueRub)
		pwd.SalesProfitFormatted = formatMoney(pwd.SalesProfitRub)
		pwd.AfterSalesRevenueFormatted = formatMoney(pwd.AfterSalesRevenueRub)
		pwd.AfterSalesProfitFormatted = formatMoney(pwd.AfterSalesProfitRub)

		results = append(results, pwd)
	}

	return results, nil
}

// FindPerformances получает записи производительности по региону (для обратной совместимости).
// Deprecated: Используйте GetWithDetailsByPeriod для получения данных с деталями.
func (r *PerformanceRepository) FindPerformances(ctx context.Context, region string) ([]*model.Performance, error) {
	// Получаем данные за текущий квартал (можно настроить)
	// Для примера используем q1 2025
	return r.GetAllByPeriod(ctx, "q1", 2025)
}

// formatMoney форматирует денежное значение в строку вида "5 555 555".
func formatMoney(amount int64) string {
	if amount == 0 {
		return "0"
	}

	str := fmt.Sprintf("%d", amount)
	n := len(str)
	if n <= 3 {
		return str
	}

	var result []byte
	for i, digit := range str {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, ' ')
		}
		result = append(result, byte(digit))
	}

	return string(result)
}
