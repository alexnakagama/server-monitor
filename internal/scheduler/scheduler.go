package scheduler

import (
	"context"
	"log"
	"time"

	"golang.org/x/text/cases"
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
	// checks if the context was already canceled
	// if its canceled returns
	// if it isnt continues executing the function
	select {
	case <-ctx.Done():
		return

	default:
	}

	err := s.metricService.DeleteOldMetrics(ctx)
	if err != nil {
		log.Printf("metric cleanup error: %v", err)
	}

	// creating a ticker
	// one of the campos its a channel
	// one time for each 24 hours it sends a value through that channel
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		// select blocks waiting for one of the channels operations to become ready
		select {
		// receives the time sent by ticker.c
		// we dont need the value, only the fact that the ticker fired
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
