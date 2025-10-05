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

// SalesRepository репозиторий для работы с данными продаж дилеров.
type SalesRepository struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

// NewSalesRepository создает новый экземпляр репозитория.
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
			"dealer_id", "quarter", "year",
			"sales_target", "stock_hdt", "stock_mdt", "stock_ldt",
			"buyout_hdt", "buyout_mdt", "buyout_ldt",
			"foton_salesmen", "sales_trainings", "service_contracts_sales",
			"sales_decision", "created_at", "updated_at",
		).
		Values(
			sales.DealerID, sales.Quarter, sales.Year,
			sales.SalesTarget, sales.StockHDT, sales.StockMDT, sales.StockLDT,
			sales.BuyoutHDT, sales.BuyoutMDT, sales.BuyoutLDT,
			sales.FotonSalesmen, sales.SalesTrainings, sales.ServiceContractsSales,
			sales.SalesDecision, sales.CreatedAt, sales.UpdatedAt,
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

	sales.ID = id
	return id, nil
}

// GetByID получает запись продаж по ID.
func (r *SalesRepository) GetByID(ctx context.Context, id int64) (*model.Sales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"sales_target", "stock_hdt", "stock_mdt", "stock_ldt",
		"buyout_hdt", "buyout_mdt", "buyout_ldt",
		"foton_salesmen", "sales_trainings", "service_contracts_sales",
		"sales_decision", "created_at", "updated_at",
	).From(salesTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("SalesRepository.GetByID: error building query: %w", err)
	}

	sales := &model.Sales{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&sales.ID, &sales.DealerID, &sales.Quarter, &sales.Year,
		&sales.SalesTarget, &sales.StockHDT, &sales.StockMDT, &sales.StockLDT,
		&sales.BuyoutHDT, &sales.BuyoutMDT, &sales.BuyoutLDT,
		&sales.FotonSalesmen, &sales.SalesTrainings, &sales.ServiceContractsSales,
		&sales.SalesDecision, &sales.CreatedAt, &sales.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("SalesRepository.GetByID: error scanning: %w", err)
	}

	return sales, nil
}

// GetByDealerAndPeriod получает запись продаж по дилеру и периоду.
func (r *SalesRepository) GetByDealerAndPeriod(ctx context.Context, dealerID int64, quarter string, year int) (*model.Sales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"sales_target", "stock_hdt", "stock_mdt", "stock_ldt",
		"buyout_hdt", "buyout_mdt", "buyout_ldt",
		"foton_salesmen", "sales_trainings", "service_contracts_sales",
		"sales_decision", "created_at", "updated_at",
	).From(salesTableName).Where(squirrel.Eq{
		"dealer_id": dealerID,
		"quarter":   quarter,
		"year":      year,
	})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("SalesRepository.GetByDealerAndPeriod: error building query: %w", err)
	}

	sales := &model.Sales{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&sales.ID, &sales.DealerID, &sales.Quarter, &sales.Year,
		&sales.SalesTarget, &sales.StockHDT, &sales.StockMDT, &sales.StockLDT,
		&sales.BuyoutHDT, &sales.BuyoutMDT, &sales.BuyoutLDT,
		&sales.FotonSalesmen, &sales.SalesTrainings, &sales.ServiceContractsSales,
		&sales.SalesDecision, &sales.CreatedAt, &sales.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("SalesRepository.GetByDealerAndPeriod: error scanning: %w", err)
	}

	return sales, nil
}

// GetAllByPeriod получает все записи продаж за указанный период.
func (r *SalesRepository) GetAllByPeriod(ctx context.Context, quarter string, year int) ([]*model.Sales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"sales_target", "stock_hdt", "stock_mdt", "stock_ldt",
		"buyout_hdt", "buyout_mdt", "buyout_ldt",
		"foton_salesmen", "sales_trainings", "service_contracts_sales",
		"sales_decision", "created_at", "updated_at",
	).From(salesTableName).Where(squirrel.Eq{
		"quarter": quarter,
		"year":    year,
	}).OrderBy("foton_salesmen DESC")

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
			&sales.ID, &sales.DealerID, &sales.Quarter, &sales.Year,
			&sales.SalesTarget, &sales.StockHDT, &sales.StockMDT, &sales.StockLDT,
			&sales.BuyoutHDT, &sales.BuyoutMDT, &sales.BuyoutLDT,
			&sales.FotonSalesmen, &sales.SalesTrainings, &sales.ServiceContractsSales,
			&sales.SalesDecision, &sales.CreatedAt, &sales.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("SalesRepository.GetAllByPeriod: error scanning: %w", err)
		}
		salesList = append(salesList, sales)
	}

	return salesList, nil
}

// Update обновляет запись продаж (частичное обновление).
func (r *SalesRepository) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("SalesRepository.Update: no fields to update")
	}

	updates["updated_at"] = time.Now()

	query := r.sq.Update(salesTableName).
		Where(squirrel.Eq{"id": id}).
		SetMap(updates)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("SalesRepository.Update: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SalesRepository.Update: error updating: %w", err)
	}

	return nil
}

// UpdateFull обновляет всю запись продаж.
func (r *SalesRepository) UpdateFull(ctx context.Context, sales *model.Sales) error {
	sales.UpdatedAt = time.Now()

	query := r.sq.Update(salesTableName).
		Set("dealer_id", sales.DealerID).
		Set("quarter", sales.Quarter).
		Set("year", sales.Year).
		Set("sales_target", sales.SalesTarget).
		Set("stock_hdt", sales.StockHDT).
		Set("stock_mdt", sales.StockMDT).
		Set("stock_ldt", sales.StockLDT).
		Set("buyout_hdt", sales.BuyoutHDT).
		Set("buyout_mdt", sales.BuyoutMDT).
		Set("buyout_ldt", sales.BuyoutLDT).
		Set("foton_salesmen", sales.FotonSalesmen).
		Set("sales_trainings", sales.SalesTrainings).
		Set("service_contracts_sales", sales.ServiceContractsSales).
		Set("sales_decision", sales.SalesDecision).
		Set("updated_at", sales.UpdatedAt).
		Where(squirrel.Eq{"id": sales.ID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("SalesRepository.UpdateFull: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SalesRepository.UpdateFull: error updating: %w", err)
	}

	return nil
}

// Delete удаляет запись продаж по ID.
func (r *SalesRepository) Delete(ctx context.Context, id int64) error {
	query := r.sq.Delete(salesTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("SalesRepository.Delete: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SalesRepository.Delete: error deleting: %w", err)
	}

	return nil
}

// GetWithDetailsByPeriod получает записи продаж с деталями дилера за период.
func (r *SalesRepository) GetWithDetailsByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.SalesWithDetails, error) {
	queryBuilder := r.sq.Select(
		"s.id", "s.dealer_id", "s.quarter", "s.year",
		"s.sales_target", "s.stock_hdt", "s.stock_mdt", "s.stock_ldt",
		"s.buyout_hdt", "s.buyout_mdt", "s.buyout_ldt",
		"s.foton_salesmen", "s.sales_trainings", "s.service_contracts_sales",
		"s.sales_decision", "s.created_at", "s.updated_at",
		"d.name as dealer_name", "d.city", "d.region", "d.manager",
	).From(salesTableName + " s").
		Join("dealers d ON s.dealer_id = d.id").
		Where(squirrel.Eq{
			"s.quarter": quarter,
			"s.year":    year,
		}).OrderBy("s.foton_salesmen DESC")

	if region != "" && region != "all-russia" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"d.region": region})
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("SalesRepository.GetWithDetailsByPeriod: error building query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("SalesRepository.GetWithDetailsByPeriod: error querying: %w", err)
	}
	defer rows.Close()

	var results []*model.SalesWithDetails
	for rows.Next() {
		swd := &model.SalesWithDetails{}
		err = rows.Scan(
			&swd.ID, &swd.DealerID, &swd.Quarter, &swd.Year,
			&swd.SalesTarget, &swd.StockHDT, &swd.StockMDT, &swd.StockLDT,
			&swd.BuyoutHDT, &swd.BuyoutMDT, &swd.BuyoutLDT,
			&swd.FotonSalesmen, &swd.SalesTrainings, &swd.ServiceContractsSales,
			&swd.SalesDecision, &swd.CreatedAt, &swd.UpdatedAt,
			&swd.DealerName, &swd.City, &swd.Region, &swd.Manager,
		)
		if err != nil {
			return nil, fmt.Errorf("SalesRepository.GetWithDetailsByPeriod: error scanning: %w", err)
		}

		// Форматируем значения stock и buyout
		swd.StockHdtMdtLdt = fmt.Sprintf("%d/%d/%d", swd.StockHDT, swd.StockMDT, swd.StockLDT)
		swd.BuyoutHdtMdtLdt = fmt.Sprintf("%d/%d/%d", swd.BuyoutHDT, swd.BuyoutMDT, swd.BuyoutLDT)

		results = append(results, swd)
	}

	return results, nil
}
