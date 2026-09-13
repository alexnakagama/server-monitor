package scheduler

import "context"

type metricService interface {
	DeleteOldMetrics(ctx context.Context) error
}

type Scheduler struct {
}
