package agent

import (
	"context"
	"log"
	"time"

	"github.com/alexnakagama/server-monitor/internal/agent/collector"
)

type Agent struct {
	interval time.Duration
	Client   *Client
}

func New(interval time.Duration, client *Client) *Agent {
	return &Agent{
		interval: interval,
		Client:   client,
	}
}

func (a *Agent) Run(ctx context.Context) {
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	for {
		metric, err := collector.CollectMetrics()
		if err != nil {
			log.Println(err)
			continue
		}

		err = a.Client.SendMetric(metric)
		if err != nil {
			log.Println(err)
			return
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}
