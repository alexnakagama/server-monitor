package repository

import (
	"context"
	"errors"
	"time"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
	"github.com/jackc/pgx/v5"
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

func (r *MetricRepository) GetByID(ctx context.Context, metricID int) (model.Metric, error) {
	query := `
		SELECT
			id,
			server_id,
			cpu_usage,
			memory_usage,
			disk_usage,
			network_receive,
			network_sent,
			timestamp
		FROM server_metrics
		WHERE id = $1
	`

	var metric model.Metric

	err := r.db.QueryRow(
		ctx,
		query,
		metricID,
	).Scan(
		&metric.ID,
		&metric.ServerID,
		&metric.CPUUsage,
		&metric.MemoryUsage,
		&metric.DiskUsage,
		&metric.NetworkReceive,
		&metric.NetworkSent,
		&metric.Timestamp,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Metric{}, errors_custom.ErrMetricNotFound
		}

		return model.Metric{}, err
	}

	return metric, nil
}

type MetricFilters struct {
	From  *time.Time
	To    *time.Time
	Limit int
}

func (r *MetricRepository) GetByServerID(ctx context.Context, serverID int, filters MetricFilters) ([]model.Metric, error) {
	query := `
		SELECT
			id,
			server_id,
			cpu_usage,
			memory_usage,
			disk_usage,
			network_receive,
			network_sent,
			timestamp
		FROM server_metrics
		WHERE server_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, serverID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metrics := make([]model.Metric, 0)

	for rows.Next() {
		var metric model.Metric

		err := rows.Scan(
			&metric.ID,
			&metric.ServerID,
			&metric.CPUUsage,
			&metric.MemoryUsage,
			&metric.DiskUsage,
			&metric.NetworkReceive,
			&metric.NetworkSent,
			&metric.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, metric)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
