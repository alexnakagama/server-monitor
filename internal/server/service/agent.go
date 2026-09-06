package service

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type AgentRepository interface {
	Create(ctx context.Context, agent model.Agent) error
	GetByID(ctx context.Context, agentID int) (model.Agent, error)
}

type AgentService struct {
	repository AgentRepository
}

func NewAgentService(repository AgentRepository) *AgentService {}
