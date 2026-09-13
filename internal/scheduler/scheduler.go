package scheduler

import "context"

type metricService interface {
	DeleteOldMetrics(ctx context.Context) error
}

type Scheduler struct {
	metricService metricService
}

func NewScheduler(metricService metricService) *Scheduler {
	return &Scheduler{
		metricService: metricService,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
}
