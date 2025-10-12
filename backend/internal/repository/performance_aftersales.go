package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// PerformanceAfterSalesRepository - репозиторий для работы с производительностью послепродажного обслуживания.
type PerformanceAfterSalesRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewPerformanceAfterSalesRepository - конструктор репозитория производительности послепродажного обслуживания.
func NewPerformanceAfterSalesRepository(pool *pgxpool.Pool, logger *slog.Logger) *PerformanceAfterSalesRepository {
	return &PerformanceAfterSalesRepository{
		pool:   pool,
		logger: logger,
	}
}

// GetAllByPeriod - получение всех записей производительности послепродажного обслуживания за период.
func (r *PerformanceAfterSalesRepository) GetAllByPeriod(ctx context.Context, period time.Time) ([]*model.PerformanceAfterSales, error) {
	r.logger.Info("Getting performance aftersales by period", "period", period)

	query := `
		SELECT id, dealer_id, period, as_revenue, as_revenue_no_vat, as_cost, 
		       as_margin, as_margin_pct, as_profit_pct, created_at, updated_at
		FROM performance_aftersales 
		WHERE period = $1
		ORDER BY dealer_id`

	rows, err := r.pool.Query(ctx, query, period)
	if err != nil {
		r.logger.Error("Failed to get performance aftersales by period", "error", err, "period", period)
		return nil, fmt.Errorf("PerformanceAfterSalesRepository.GetAllByPeriod: %w", err)
	}
	defer rows.Close()

	var performances []*model.PerformanceAfterSales
	for rows.Next() {
		var pas model.PerformanceAfterSales
		err := rows.Scan(
			&pas.ID,
			&pas.DealerID,
			&pas.Period,
			&pas.ASRevenue,
			&pas.ASRevenueNoVat,
			&pas.ASCost,
			&pas.ASMargin,
			&pas.ASMarginPct,
			&pas.ASProfitPct,
			&pas.CreatedAt,
			&pas.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan performance aftersales", "error", err)
			return nil, fmt.Errorf("PerformanceAfterSalesRepository.GetAllByPeriod: %w", err)
		}
		performances = append(performances, &pas)
	}

	return performances, nil
}

// GetByDealerID - получение производительности послепродажного обслуживания по ID дилера.
func (r *PerformanceAfterSalesRepository) GetByDealerID(ctx context.Context, dealerID int) ([]*model.PerformanceAfterSales, error) {
	r.logger.Info("Getting performance aftersales by dealer ID", "dealerID", dealerID)

	// Здесь нужно будет реализовать SQL запрос для получения данных
	// Пока возвращаем пустой список
	return []*model.PerformanceAfterSales{}, nil
}

// GetByDealerIDAndPeriod - получение производительности послепродажного обслуживания по ID дилера и периоду.
func (r *PerformanceAfterSalesRepository) GetByDealerIDAndPeriod(ctx context.Context, dealerID int, period time.Time) (*model.PerformanceAfterSales, error) {
	r.logger.Info("Getting performance aftersales by dealer ID and period", "dealerID", dealerID, "period", period)

	query := `
		SELECT id, dealer_id, period, as_revenue, as_revenue_no_vat, as_cost, 
		       as_margin, as_margin_pct, as_profit_pct, created_at, updated_at
		FROM performance_aftersales 
		WHERE dealer_id = $1 AND period = $2`

	var pas model.PerformanceAfterSales
	err := r.pool.QueryRow(ctx, query, dealerID, period).Scan(
		&pas.ID,
		&pas.DealerID,
		&pas.Period,
		&pas.ASRevenue,
		&pas.ASRevenueNoVat,
		&pas.ASCost,
		&pas.ASMargin,
		&pas.ASMarginPct,
		&pas.ASProfitPct,
		&pas.CreatedAt,
		&pas.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to get performance aftersales by dealer ID and period", "error", err, "dealerID", dealerID, "period", period)
		return nil, fmt.Errorf("PerformanceAfterSalesRepository.GetByDealerIDAndPeriod: %w", err)
	}

	return &pas, nil
}

// GetByID - получение производительности послепродажного обслуживания по ID.
func (r *PerformanceAfterSalesRepository) GetByID(ctx context.Context, id int) (*model.PerformanceAfterSales, error) {
	r.logger.Info("Getting performance aftersales by ID", "id", id)

	query := `
		SELECT id, dealer_id, period, as_revenue, as_revenue_no_vat, as_cost, 
		       as_margin, as_margin_pct, as_profit_pct, created_at, updated_at
		FROM performance_aftersales 
		WHERE id = $1`

	var pas model.PerformanceAfterSales
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&pas.ID,
		&pas.DealerID,
		&pas.Period,
		&pas.ASRevenue,
		&pas.ASRevenueNoVat,
		&pas.ASCost,
		&pas.ASMargin,
		&pas.ASMarginPct,
		&pas.ASProfitPct,
		&pas.CreatedAt,
		&pas.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to get performance aftersales by ID", "error", err, "id", id)
		return nil, fmt.Errorf("PerformanceAfterSalesRepository.GetByID: %w", err)
	}

	return &pas, nil
}

// Create - создание новой записи производительности послепродажного обслуживания.
func (r *PerformanceAfterSalesRepository) Create(ctx context.Context, perf *model.PerformanceAfterSales) (int, error) {
	r.logger.Info("Creating performance aftersales", "dealerID", perf.DealerID)

	query := `
		INSERT INTO performance_aftersales (dealer_id, period, as_revenue, as_revenue_no_vat, as_cost, 
		                                    as_margin, as_margin_pct, as_profit_pct, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	var id int
	err := r.pool.QueryRow(ctx, query,
		perf.DealerID,
		perf.Period,
		perf.ASRevenue,
		perf.ASRevenueNoVat,
		perf.ASCost,
		perf.ASMargin,
		perf.ASMarginPct,
		perf.ASProfitPct,
		perf.CreatedAt,
		perf.UpdatedAt,
	).Scan(&id)

	if err != nil {
		r.logger.Error("Failed to create performance aftersales", "error", err, "dealer_id", perf.DealerID)
		return 0, fmt.Errorf("PerformanceAfterSalesRepository.Create: %w", err)
	}

	r.logger.Info("Performance aftersales created successfully", "id", id, "dealer_id", perf.DealerID)
	return id, nil
}

// Update - обновление записи производительности послепродажного обслуживания.
func (r *PerformanceAfterSalesRepository) Update(ctx context.Context, perf *model.PerformanceAfterSales) error {
	r.logger.Info("Updating performance aftersales", "id", perf.ID)

	query := `
		UPDATE performance_aftersales 
		SET dealer_id = $1, period = $2, as_revenue = $3, as_revenue_no_vat = $4, as_cost = $5, 
		    as_margin = $6, as_margin_pct = $7, as_profit_pct = $8, updated_at = $9
		WHERE id = $10`

	result, err := r.pool.Exec(ctx, query,
		perf.DealerID,
		perf.Period,
		perf.ASRevenue,
		perf.ASRevenueNoVat,
		perf.ASCost,
		perf.ASMargin,
		perf.ASMarginPct,
		perf.ASProfitPct,
		perf.UpdatedAt,
		perf.ID,
	)

	if err != nil {
		r.logger.Error("Failed to update performance aftersales", "error", err, "id", perf.ID)
		return fmt.Errorf("PerformanceAfterSalesRepository.Update: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("PerformanceAfterSalesRepository.Update: no rows affected, record with id %d not found", perf.ID)
	}

	return nil
}

// Delete - удаление записи производительности послепродажного обслуживания.
func (r *PerformanceAfterSalesRepository) Delete(ctx context.Context, id int) error {
	r.logger.Info("Deleting performance aftersales", "id", id)

	query := `DELETE FROM performance_aftersales WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete performance aftersales", "error", err, "id", id)
		return fmt.Errorf("PerformanceAfterSalesRepository.Delete: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("PerformanceAfterSalesRepository.Delete: no rows affected, record with id %d not found", id)
	}

	return nil
}
