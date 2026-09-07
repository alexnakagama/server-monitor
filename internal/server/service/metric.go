package service

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type MetricRepository interface {
	Create(ctx context.Context, metric model.Metric) error
}
