package service

import (
	"context"
	"time"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type MetricRepositoryMock struct {
	create          func(ctx context.Context, metric model.Metric) error
	getByID         func(ctx context.Context, metricID int) (model.Metric, error)
	getByServerID   func(ctx context.Context, serverID int, filters model.MetricFilters) ([]model.Metric, error)
	deleteOlderThan func(ctx context.Context, before time.Time) error
}
