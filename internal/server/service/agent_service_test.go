package service

import (
	"context"
	"testing"

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
	return m.create(ctx, agent)
}

func (m *AgentRepositoryMock) GetByID(ctx context.Context, agentID int) (model.Agent, error) {
	return m.getByID(ctx, agentID)
}

func (m *AgentRepositoryMock) DeleteByID(ctx context.Context, agentID int) error {
	return m.deleteByID(ctx, agentID)
}

func (m *AgentRepositoryMock) GetByTokenHash(ctx context.Context, tokenHash string) (model.Agent, error) {
	return m.getByTokenHash(ctx, tokenHash)
}

func TestAgentService_Create(t *testing.T) {
	var createdAgent model.Agent

	repository := &AgentRepositoryMock{
		create: func(ctx context.Context, agent model.Agent) error {
			createdAgent = agent
			return nil
		},
	}

	service := NewAgentService(repository)

	token, err := service.Create(context.Background(), 1, "my-server-agent")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if createdAgent.ServerID != 1 {
		t.Errorf("expected server id 1, got: %d", createdAgent.ServerID)
	}

	if createdAgent.Name != "my-server-agent" {
		t.Errorf("expected name my-server-agent, got: %s", createdAgent.Name)
	}

	if token == "" {
		t.Errorf("expected token to be generated")
	}

	expectedHash := HashAgentToken(token)

	if createdAgent.TokenHash != expectedHash {
		t.Errorf("expected token hash %s, got: %s", expectedHash, createdAgent.TokenHash)
	}
}

func TestAgentService_GetByID(t *testing.T) {}

func TestAgentService_DeleteByID(t *testing.T) {}

func TestAgentService_GetByTokenHash(t *testing.T) {}
