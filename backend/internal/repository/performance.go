package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

const performanceTableName = "dealer_performance"

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

// FindPerformances получение performance по региону.
func (p *PerformanceRepository) FindPerformances(ctx context.Context, region string) ([]*model.Performance, error) {
	var performances []*model.Performance
	query := p.sq.Select("dealer_name", "region", "city", "manager",
		"sales_revenues", "sales_profits", "sales_margin",
		"as_revenues", "as_profits", "as_margin",
		"foton_rank", "dealer_dev_recommendation").
		From(performanceTableName).Where(squirrel.Eq{"region": region})
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.FindPerformances error selecting performances: %w", err)
	}

	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PerformanceRepository.FindPerformances error selecting performances: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		performance := &model.Performance{}
		err = rows.Scan(
			&performance.DealerName,
			&performance.Region,
			&performance.City,
			&performance.Manager,
			&performance.SalesRevenues,
			&performance.SalesProfits,
			&performance.SalesMargin,
			&performance.ASRevenues,
			&performance.ASProfits,
			&performance.ASMargin,
			&performance.FotonRank,
			&performance.Recommendation,
		)
		if err != nil {
			return nil, fmt.Errorf("PerformanceRepository.FindPerformances error scanning performances: %w", err)
		}
		performances = append(performances, performance)
	}

	return performances, nil
}
