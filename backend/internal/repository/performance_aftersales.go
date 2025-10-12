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

	// Здесь нужно будет реализовать SQL запрос для получения данных
	// Пока возвращаем пустой список
	return []*model.PerformanceAfterSales{}, nil
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

	// Здесь нужно будет реализовать SQL запрос для получения данных
	// Пока возвращаем заглушку
	return nil, nil
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

	// Здесь нужно будет реализовать SQL запрос для обновления записи
	return nil
}

// Delete - удаление записи производительности послепродажного обслуживания.
func (r *PerformanceAfterSalesRepository) Delete(ctx context.Context, id int) error {
	r.logger.Info("Deleting performance aftersales", "id", id)

	// Здесь нужно будет реализовать SQL запрос для удаления записи
	return nil
}
