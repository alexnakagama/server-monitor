package repository

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MetricRepository struct {
	db *pgxpool.Pool
}

func NewMetricRepository(db *pgxpool.Pool) *MetricRepository {
	return &MetricRepository{
		db: db,
	}
}

func (r *MetricRepository) Create(ctx context.Context, metric model.Metric) error {
}
