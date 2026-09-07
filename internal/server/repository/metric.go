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
	query := `
		INSERT INTO server_metrics (
			server_id,
			cpu_usage,
			memory_usage,
			disk_usage,
			network_receive,
			network_sent
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		metric.ServerID,
		metric.CPUUsage,
		metric.MemoryUsage,
		metric.DiskUsage,
		metric.NetworkReceive,
		metric.NetworkSent,
	)

	return err
}
