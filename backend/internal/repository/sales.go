package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

const salesTableName = "sales"

// SalesRepository репозиторий для работы с данными продаж.
type SalesRepository struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

// NewSalesRepository конструктор.
func NewSalesRepository(pool *pgxpool.Pool) *SalesRepository {
	return &SalesRepository{
		pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Create создает новую запись продаж.
func (r *SalesRepository) Create(ctx context.Context, sales *model.Sales) (int64, error) {
	now := time.Now()
	sales.CreatedAt = now
	sales.UpdatedAt = now

	query := r.sq.Insert(salesTableName).
		Columns(
			"dealer_id", "period",
			"stock_hdt", "stock_mdt", "stock_ldt",
			"buyout_hdt", "buyout_mdt", "buyout_ldt",
			"foton_sales_personnel", "sales_target_plan", "sales_target_fact",
			"service_contracts_sales", "sales_trainings", "sales_recommendation",
			"created_at", "updated_at",
		).
		Values(
			sales.DealerID, sales.Period,
			sales.StockHDT, sales.StockMDT, sales.StockLDT,
			sales.BuyoutHDT, sales.BuyoutMDT, sales.BuyoutLDT,
			sales.FotonSalesPersonnel, sales.SalesTargetPlan, sales.SalesTargetFact,
			sales.ServiceContractsSales, sales.SalesTrainings, sales.SalesRecommendation,
			sales.CreatedAt, sales.UpdatedAt,
		).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("SalesRepository.Create: error building query: %w", err)
	}

	var id int64
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("SalesRepository.Create: error inserting: %w", err)
	}

	sales.ID = int(id)
	return id, nil
}

// GetByID получает запись продаж по ID.
func (r *SalesRepository) GetByID(ctx context.Context, id int64) (*model.Sales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "period",
		"stock_hdt", "stock_mdt", "stock_ldt",
		"buyout_hdt", "buyout_mdt", "buyout_ldt",
		"foton_sales_personnel", "sales_target_plan", "sales_target_fact",
		"service_contracts_sales", "sales_trainings", "sales_recommendation",
		"created_at", "updated_at",
	).From(salesTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("SalesRepository.GetByID: error building query: %w", err)
	}

	sales := &model.Sales{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&sales.ID, &sales.DealerID, &sales.Period,
		&sales.StockHDT, &sales.StockMDT, &sales.StockLDT,
		&sales.BuyoutHDT, &sales.BuyoutMDT, &sales.BuyoutLDT,
		&sales.FotonSalesPersonnel, &sales.SalesTargetPlan, &sales.SalesTargetFact,
		&sales.ServiceContractsSales, &sales.SalesTrainings, &sales.SalesRecommendation,
		&sales.CreatedAt, &sales.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("SalesRepository.GetByID: error scanning: %w", err)
	}

	return sales, nil
}

// GetByDealerAndPeriod получает запись продаж по дилеру и периоду.
func (r *SalesRepository) GetByDealerAndPeriod(ctx context.Context, dealerID int64, quarter string, year int) (*model.Sales, error) {
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

// GetByDealerAndPeriodTime получает запись продаж по дилеру и периоду.
func (r *SalesRepository) GetByDealerAndPeriodTime(ctx context.Context, dealerID int, period time.Time) (*model.Sales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "period",
		"stock_hdt", "stock_mdt", "stock_ldt",
		"buyout_hdt", "buyout_mdt", "buyout_ldt",
		"foton_sales_personnel", "sales_target_plan", "sales_target_fact",
		"service_contracts_sales", "sales_trainings", "sales_recommendation",
		"created_at", "updated_at",
	).From(salesTableName).Where(squirrel.Eq{
		"dealer_id": dealerID,
		"period":    period,
	})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("SalesRepository.GetByDealerAndPeriod: error building query: %w", err)
	}

	sales := &model.Sales{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&sales.ID, &sales.DealerID, &sales.Period,
		&sales.StockHDT, &sales.StockMDT, &sales.StockLDT,
		&sales.BuyoutHDT, &sales.BuyoutMDT, &sales.BuyoutLDT,
		&sales.FotonSalesPersonnel, &sales.SalesTargetPlan, &sales.SalesTargetFact,
		&sales.ServiceContractsSales, &sales.SalesTrainings, &sales.SalesRecommendation,
		&sales.CreatedAt, &sales.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("SalesRepository.GetByDealerAndPeriod: error scanning: %w", err)
	}

	return sales, nil
}

// GetAllByPeriod получает все записи продаж за указанный период.
func (r *SalesRepository) GetAllByPeriod(ctx context.Context, quarter string, year int) ([]*model.Sales, error) {
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

// GetAllByPeriodTime получает все записи продаж за указанный период.
func (r *SalesRepository) GetAllByPeriodTime(ctx context.Context, period time.Time) ([]*model.Sales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "period",
		"stock_hdt", "stock_mdt", "stock_ldt",
		"buyout_hdt", "buyout_mdt", "buyout_ldt",
		"foton_sales_personnel", "sales_target_plan", "sales_target_fact",
		"service_contracts_sales", "sales_trainings", "sales_recommendation",
		"created_at", "updated_at",
	).From(salesTableName).Where(squirrel.Eq{"period": period})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("SalesRepository.GetAllByPeriod: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("SalesRepository.GetAllByPeriod: error querying: %w", err)
	}
	defer rows.Close()

	var salesList []*model.Sales
	for rows.Next() {
		sales := &model.Sales{}
		err = rows.Scan(
			&sales.ID, &sales.DealerID, &sales.Period,
			&sales.StockHDT, &sales.StockMDT, &sales.StockLDT,
			&sales.BuyoutHDT, &sales.BuyoutMDT, &sales.BuyoutLDT,
			&sales.FotonSalesPersonnel, &sales.SalesTargetPlan, &sales.SalesTargetFact,
			&sales.ServiceContractsSales, &sales.SalesTrainings, &sales.SalesRecommendation,
			&sales.CreatedAt, &sales.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("SalesRepository.GetAllByPeriod: error scanning: %w", err)
		}
		salesList = append(salesList, sales)
	}

	return salesList, nil
}

// UpdateFull обновляет всю запись продаж целиком.
func (r *SalesRepository) UpdateFull(ctx context.Context, sales *model.Sales) error {
	sales.UpdatedAt = time.Now()

	query := r.sq.Update(salesTableName).
		Set("dealer_id", sales.DealerID).
		Set("period", sales.Period).
		Set("stock_hdt", sales.StockHDT).
		Set("stock_mdt", sales.StockMDT).
		Set("stock_ldt", sales.StockLDT).
		Set("buyout_hdt", sales.BuyoutHDT).
		Set("buyout_mdt", sales.BuyoutMDT).
		Set("buyout_ldt", sales.BuyoutLDT).
		Set("foton_sales_personnel", sales.FotonSalesPersonnel).
		Set("sales_target_plan", sales.SalesTargetPlan).
		Set("sales_target_fact", sales.SalesTargetFact).
		Set("service_contracts_sales", sales.ServiceContractsSales).
		Set("sales_trainings", sales.SalesTrainings).
		Set("sales_recommendation", sales.SalesRecommendation).
		Set("updated_at", sales.UpdatedAt).
		Where(squirrel.Eq{"id": sales.ID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("SalesRepository.UpdateFull: error building query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SalesRepository.UpdateFull: error updating: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("SalesRepository.UpdateFull: no rows affected, record with id %d not found", sales.ID)
	}

	return nil
}

// Delete удаляет запись продаж.
func (r *SalesRepository) Delete(ctx context.Context, id int64) error {
	query := r.sq.Delete(salesTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("SalesRepository.Delete: error building query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SalesRepository.Delete: error deleting: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("SalesRepository.Delete: no rows affected, record with id %d not found", id)
	}

	return nil
}

// Update обновляет данные продаж.
func (r *SalesRepository) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	updates["updated_at"] = time.Now()
	query := r.sq.Update(salesTableName).SetMap(updates).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("SalesRepository.Update: error building query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SalesRepository.Update: error updating: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("SalesRepository.Update: no rows affected, record with id %d not found", id)
	}

	return nil
}

// GetWithDetailsByPeriod получает записи продаж с деталями за период.
func (r *SalesRepository) GetWithDetailsByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.SalesWithDetails, error) {
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

// GetWithDetailsByPeriodTime получает записи продаж с деталями за период.
func (r *SalesRepository) GetWithDetailsByPeriodTime(ctx context.Context, period time.Time, region string) ([]*model.SalesWithDetails, error) {
	queryBuilder := r.sq.Select(
		"s.id", "s.dealer_id", "s.period",
		"s.stock_hdt", "s.stock_mdt", "s.stock_ldt",
		"s.buyout_hdt", "s.buyout_mdt", "s.buyout_ldt",
		"s.foton_sales_personnel", "s.sales_target_plan", "s.sales_target_fact",
		"s.service_contracts_sales", "s.sales_trainings", "s.sales_recommendation",
		"s.created_at", "s.updated_at",
		"d.dealer_name_ru", "d.dealer_name_en", "d.city", "d.region", "d.manager", "d.ruft",
	).
		From(salesTableName + " s").
		Join("dealers d ON s.dealer_id = d.dealer_id").
		Where(squirrel.Eq{"s.period": period})

	if region != "all-russia" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"d.region": region})
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("SalesRepository.GetWithDetailsByPeriodTime: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("SalesRepository.GetWithDetailsByPeriodTime: error querying: %w", err)
	}
	defer rows.Close()

	var results []*model.SalesWithDetails
	for rows.Next() {
		swd := &model.SalesWithDetails{}
		err = rows.Scan(
			&swd.ID, &swd.DealerID, &swd.Period,
			&swd.StockHDT, &swd.StockMDT, &swd.StockLDT,
			&swd.BuyoutHDT, &swd.BuyoutMDT, &swd.BuyoutLDT,
			&swd.FotonSalesPersonnel, &swd.SalesTargetPlan, &swd.SalesTargetFact,
			&swd.ServiceContractsSales, &swd.SalesTrainings, &swd.SalesRecommendation,
			&swd.CreatedAt, &swd.UpdatedAt,
			&swd.DealerNameRu, &swd.DealerNameEn, &swd.City, &swd.Region, &swd.Manager, &swd.Ruft,
		)
		if err != nil {
			return nil, fmt.Errorf("SalesRepository.GetWithDetailsByPeriodTime: error scanning: %w", err)
		}
		results = append(results, swd)
	}

	return results, nil
}
