package scheduler

import (
	"context"
	"log"
	"time"
)

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
	// creating a ticker
	// one of the campos its a channel
	// one time for each 24 hours it sends a value through that channel
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		// select stays blocked waiting for one of the channels to execute
		select {
		// receives a value of the ticker.c channel
		// in this case we arent interested in the value but the moment of the ticker
		case <-ticker.C:
			err := s.metricService.DeleteOldMetrics(ctx)
			if err != nil {
				log.Printf("metric cleanup error: %v", err)
			}

		case <-ctx.Done():
			return
		}
	}
}
