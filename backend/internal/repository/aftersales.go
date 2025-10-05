package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

const afterSalesTableName = "after_sales"

// AfterSalesRepository репозиторий для работы с данными послепродажного обслуживания.
type AfterSalesRepository struct {
	pool *pgxpool.Pool
	sq   squirrel.StatementBuilderType
}

// NewAfterSalesRepository создает новый экземпляр репозитория.
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
			"dealer_id", "quarter", "year",
			"recommended_stock", "warranty_stock", "foton_labor_hours",
			"service_contracts", "as_trainings", "csi", "foton_warranty_hours",
			"as_decision", "created_at", "updated_at",
		).
		Values(
			as.DealerID, as.Quarter, as.Year,
			as.RecommendedStock, as.WarrantyStock, as.FotonLaborHours,
			as.ServiceContracts, as.ASTrainings, as.CSI, as.FotonWarrantyHours,
			as.AfterSalesDecision, as.CreatedAt, as.UpdatedAt,
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

	as.ID = id
	return id, nil
}

// GetByID получает запись послепродажного обслуживания по ID.
func (r *AfterSalesRepository) GetByID(ctx context.Context, id int64) (*model.AfterSales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"recommended_stock", "warranty_stock", "foton_labor_hours",
		"service_contracts", "as_trainings", "csi", "foton_warranty_hours",
		"as_decision", "created_at", "updated_at",
	).From(afterSalesTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("AfterSalesRepository.GetByID: error building query: %w", err)
	}

	as := &model.AfterSales{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&as.ID, &as.DealerID, &as.Quarter, &as.Year,
		&as.RecommendedStock, &as.WarrantyStock, &as.FotonLaborHours,
		&as.ServiceContracts, &as.ASTrainings, &as.CSI, &as.FotonWarrantyHours,
		&as.AfterSalesDecision, &as.CreatedAt, &as.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("AfterSalesRepository.GetByID: error scanning: %w", err)
	}

	return as, nil
}

// GetByDealerAndPeriod получает запись послепродажного обслуживания по дилеру и периоду.
func (r *AfterSalesRepository) GetByDealerAndPeriod(ctx context.Context, dealerID int64, quarter string, year int) (*model.AfterSales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"recommended_stock", "warranty_stock", "foton_labor_hours",
		"service_contracts", "as_trainings", "csi", "foton_warranty_hours",
		"as_decision", "created_at", "updated_at",
	).From(afterSalesTableName).Where(squirrel.Eq{
		"dealer_id": dealerID,
		"quarter":   quarter,
		"year":      year,
	})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("AfterSalesRepository.GetByDealerAndPeriod: error building query: %w", err)
	}

	as := &model.AfterSales{}
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&as.ID, &as.DealerID, &as.Quarter, &as.Year,
		&as.RecommendedStock, &as.WarrantyStock, &as.FotonLaborHours,
		&as.ServiceContracts, &as.ASTrainings, &as.CSI, &as.FotonWarrantyHours,
		&as.AfterSalesDecision, &as.CreatedAt, &as.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("AfterSalesRepository.GetByDealerAndPeriod: error scanning: %w", err)
	}

	return as, nil
}

// GetAllByPeriod получает все записи послепродажного обслуживания за указанный период.
func (r *AfterSalesRepository) GetAllByPeriod(ctx context.Context, quarter string, year int) ([]*model.AfterSales, error) {
	query := r.sq.Select(
		"id", "dealer_id", "quarter", "year",
		"recommended_stock", "warranty_stock", "foton_labor_hours",
		"service_contracts", "as_trainings", "csi", "foton_warranty_hours",
		"as_decision", "created_at", "updated_at",
	).From(afterSalesTableName).Where(squirrel.Eq{
		"quarter": quarter,
		"year":    year,
	}).OrderBy("service_contracts DESC")

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
			&as.ID, &as.DealerID, &as.Quarter, &as.Year,
			&as.RecommendedStock, &as.WarrantyStock, &as.FotonLaborHours,
			&as.ServiceContracts, &as.ASTrainings, &as.CSI, &as.FotonWarrantyHours,
			&as.AfterSalesDecision, &as.CreatedAt, &as.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("AfterSalesRepository.GetAllByPeriod: error scanning: %w", err)
		}
		afterSalesList = append(afterSalesList, as)
	}

	return afterSalesList, nil
}

// Update обновляет запись послепродажного обслуживания (частичное обновление).
func (r *AfterSalesRepository) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("AfterSalesRepository.Update: no fields to update")
	}

	updates["updated_at"] = time.Now()

	query := r.sq.Update(afterSalesTableName).
		Where(squirrel.Eq{"id": id}).
		SetMap(updates)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("AfterSalesRepository.Update: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AfterSalesRepository.Update: error updating: %w", err)
	}

	return nil
}

// UpdateFull обновляет всю запись послепродажного обслуживания.
func (r *AfterSalesRepository) UpdateFull(ctx context.Context, as *model.AfterSales) error {
	as.UpdatedAt = time.Now()

	query := r.sq.Update(afterSalesTableName).
		Set("dealer_id", as.DealerID).
		Set("quarter", as.Quarter).
		Set("year", as.Year).
		Set("recommended_stock", as.RecommendedStock).
		Set("warranty_stock", as.WarrantyStock).
		Set("foton_labor_hours", as.FotonLaborHours).
		Set("service_contracts", as.ServiceContracts).
		Set("as_trainings", as.ASTrainings).
		Set("csi", as.CSI).
		Set("foton_warranty_hours", as.FotonWarrantyHours).
		Set("as_decision", as.AfterSalesDecision).
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

// Delete удаляет запись послепродажного обслуживания по ID.
func (r *AfterSalesRepository) Delete(ctx context.Context, id int64) error {
	query := r.sq.Delete(afterSalesTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("AfterSalesRepository.Delete: error building query: %w", err)
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AfterSalesRepository.Delete: error deleting: %w", err)
	}

	return nil
}

// GetWithDetailsByPeriod получает записи послепродажного обслуживания с деталями дилера за период.
func (r *AfterSalesRepository) GetWithDetailsByPeriod(ctx context.Context, quarter string, year int, region string) ([]*model.AfterSalesWithDetails, error) {
	queryBuilder := r.sq.Select(
		"a.id", "a.dealer_id", "a.quarter", "a.year",
		"a.recommended_stock", "a.warranty_stock", "a.foton_labor_hours",
		"a.service_contracts", "a.as_trainings", "a.csi", "a.foton_warranty_hours",
		"a.as_decision", "a.created_at", "a.updated_at",
		"d.name as dealer_name", "d.city", "d.region", "d.manager",
	).From(afterSalesTableName + " a").
		Join("dealers d ON a.dealer_id = d.id").
		Where(squirrel.Eq{
			"a.quarter": quarter,
			"a.year":    year,
		}).OrderBy("a.service_contracts DESC")

	if region != "" && region != "all-russia" {
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
			&aswd.ID, &aswd.DealerID, &aswd.Quarter, &aswd.Year,
			&aswd.RecommendedStock, &aswd.WarrantyStock, &aswd.FotonLaborHours,
			&aswd.ServiceContracts, &aswd.ASTrainings, &aswd.CSI, &aswd.FotonWarrantyHours,
			&aswd.AfterSalesDecision, &aswd.CreatedAt, &aswd.UpdatedAt,
			&aswd.DealerName, &aswd.City, &aswd.Region, &aswd.Manager,
		)
		if err != nil {
			return nil, fmt.Errorf("AfterSalesRepository.GetWithDetailsByPeriod: error scanning: %w", err)
		}

		results = append(results, aswd)
	}

	return results, nil
}
