package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// PerformanceSalesRepository - репозиторий для работы с производительностью продаж.
type PerformanceSalesRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewPerformanceSalesRepository - конструктор репозитория производительности продаж.
func NewPerformanceSalesRepository(pool *pgxpool.Pool, logger *slog.Logger) *PerformanceSalesRepository {
	return &PerformanceSalesRepository{
		pool:   pool,
		logger: logger,
	}
}

// GetAllByPeriod - получение всех записей производительности продаж за период.
func (r *PerformanceSalesRepository) GetAllByPeriod(ctx context.Context, period time.Time) ([]*model.PerformanceSales, error) {
	r.logger.Info("Getting performance sales by period", "period", period)

	// Здесь нужно будет реализовать SQL запрос для получения данных
	// Пока возвращаем пустой список
	return []*model.PerformanceSales{}, nil
}

// GetByDealerID - получение производительности продаж по ID дилера.
func (r *PerformanceSalesRepository) GetByDealerID(ctx context.Context, dealerID int) ([]*model.PerformanceSales, error) {
	r.logger.Info("Getting performance sales by dealer ID", "dealerID", dealerID)

	// Здесь нужно будет реализовать SQL запрос для получения данных
	// Пока возвращаем пустой список
	return []*model.PerformanceSales{}, nil
}

// GetByDealerIDAndPeriod - получение производительности продаж по ID дилера и периоду.
func (r *PerformanceSalesRepository) GetByDealerIDAndPeriod(ctx context.Context, dealerID int, period time.Time) (*model.PerformanceSales, error) {
	r.logger.Info("Getting performance sales by dealer ID and period", "dealerID", dealerID, "period", period)

	// Здесь нужно будет реализовать SQL запрос для получения данных
	// Пока возвращаем заглушку
	return nil, nil
}

// GetByID - получение производительности продаж по ID.
func (r *PerformanceSalesRepository) GetByID(ctx context.Context, id int) (*model.PerformanceSales, error) {
	r.logger.Info("Getting performance sales by ID", "id", id)

	query := `
		SELECT id, dealer_id, period, quantity_sold, sales_revenue, sales_revenue_no_vat, sales_cost, 
		       sales_margin, sales_margin_pct, sales_profit_pct, created_at, updated_at
		FROM performance_sales 
		WHERE id = $1`

	var ps model.PerformanceSales
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&ps.ID,
		&ps.DealerID,
		&ps.Period,
		&ps.QuantitySold,
		&ps.SalesRevenue,
		&ps.SalesRevenueNoVat,
		&ps.SalesCost,
		&ps.SalesMargin,
		&ps.SalesMarginPct,
		&ps.SalesProfitPct,
		&ps.CreatedAt,
		&ps.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to get performance sales by ID", "error", err, "id", id)
		return nil, fmt.Errorf("PerformanceSalesRepository.GetByID: %w", err)
	}

	return &ps, nil
}

// Create - создание новой записи производительности продаж.
func (r *PerformanceSalesRepository) Create(ctx context.Context, perf *model.PerformanceSales) (int, error) {
	r.logger.Info("Creating performance sales", "dealerID", perf.DealerID)

	query := `
		INSERT INTO performance_sales (dealer_id, period, quantity_sold, sales_revenue, sales_revenue_no_vat, sales_cost, 
		                              sales_margin, sales_margin_pct, sales_profit_pct, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	var id int
	err := r.pool.QueryRow(ctx, query,
		perf.DealerID,
		perf.Period,
		perf.QuantitySold,
		perf.SalesRevenue,
		perf.SalesRevenueNoVat,
		perf.SalesCost,
		perf.SalesMargin,
		perf.SalesMarginPct,
		perf.SalesProfitPct,
		perf.CreatedAt,
		perf.UpdatedAt,
	).Scan(&id)

	if err != nil {
		r.logger.Error("Failed to create performance sales", "error", err, "dealer_id", perf.DealerID)
		return 0, fmt.Errorf("PerformanceSalesRepository.Create: %w", err)
	}

	r.logger.Info("Performance sales created successfully", "id", id, "dealer_id", perf.DealerID)
	return id, nil
}

// Update - обновление записи производительности продаж.
func (r *PerformanceSalesRepository) Update(ctx context.Context, perf *model.PerformanceSales) error {
	r.logger.Info("Updating performance sales", "id", perf.ID)

	// Здесь нужно будет реализовать SQL запрос для обновления записи
	return nil
}

// Delete - удаление записи производительности продаж.
func (r *PerformanceSalesRepository) Delete(ctx context.Context, id int) error {
	r.logger.Info("Deleting performance sales", "id", id)

	// Здесь нужно будет реализовать SQL запрос для удаления записи
	return nil
}
