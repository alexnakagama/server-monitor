package agent

import "time"

type Agent struct {
	interval time.Duration
}

func New(interval time.Duration) *Agent {
	return &Agent{
		interval: interval,
	}
}

func (a *Agent) Run() {}
