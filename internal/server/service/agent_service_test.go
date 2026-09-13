package service

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
)

var _ AgentRepository = (*AgentRepositoryMock)(nil)

type AgentRepositoryMock struct {
	create         func(ctx context.Context, agent model.Agent) error
	getByID        func(ctx context.Context, agentID int) (model.Agent, error)
	deleteByID     func(ctx context.Context, agentID int) error
	getByTokenHash func(ctx context.Context, tokenHash string) (model.Agent, error)
}

func (m *AgentRepositoryMock) Create(ctx context.Context, agent model.Agent) error {
}
