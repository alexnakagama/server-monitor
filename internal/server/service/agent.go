package service

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type AgentRepository interface {
	Create(ctx context.Context, agent model.Agent) error
}
