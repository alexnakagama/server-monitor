package service

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type MetricRepository interface {
	Create(ctx context.Context, metric model.Metric) error
	GetByID(ctx context.Context, metricID int) (model.Metric, error)
}

type MetricService struct {
	repository MetricRepository
}
