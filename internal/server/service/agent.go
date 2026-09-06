package service

import (
	"context"
	"errors"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
)

type AgentRepository interface {
	Create(ctx context.Context, agent model.Agent) error
	GetByID(ctx context.Context, agentID int) (model.Agent, error)
}

type AgentService struct {
	repository AgentRepository
}

func NewAgentService(repository AgentRepository) *AgentService {
	return &AgentService{
		repository: repository,
	}
}

func (s *AgentService) Create(ctx context.Context, serverID int, name string) error {
	err := model.ValidateServerID(serverID)
	if err != nil {
		return errors_custom.ErrInvalidServerID
	}

	err = model.ValidateAgentName(name)
	if err != nil {
		return errors_custom.ErrNameRequired
	}

	// Todo function to generate the token hash to pass to the agent to create it

	agent := model.Agent{
		ServerID: serverID,
		Name:     name,
	}

	return s.repository.Create(ctx, agent)
}

func (s *AgentService) GetByID(ctx context.Context, agentID int) (model.Agent, error) {
	return s.repository.GetByID(ctx, agentID)
}
