package agent

import (
	"log"
	"time"

	"github.com/alexnakagama/server-monitor/internal/agent/collector"
)

type Agent struct {
	interval time.Duration
	Client   *Client
}

func New(interval time.Duration) *Agent {
	return &Agent{
		interval: interval,
	}
}

func (a *Agent) Run() {
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

		time.Sleep(a.interval)
	}
}
