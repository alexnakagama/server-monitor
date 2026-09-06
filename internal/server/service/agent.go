package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
)

type AgentRepository interface {
	Create(ctx context.Context, agent model.Agent) error
	GetByID(ctx context.Context, agentID int) (model.Agent, error)
	DeleteByID(ctx context.Context, agentID int) error
	GetByTokenHash(ctx context.Context, tokenHash string) (model.Agent, error)
}

type AgentService struct {
	repository AgentRepository
}

func NewAgentService(repository AgentRepository) *AgentService {
	return &AgentService{
		repository: repository,
	}
}

func GenerateAgentToken() (string, error) {
	token := make([]byte, 32)

	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(token), nil
}

func HashAgentToken(token string) string {
	hash := sha256.Sum256([]byte(token))

	return hex.EncodeToString(hash[:])
}

func (s *AgentService) Create(ctx context.Context, serverID int, name string) (string, error) {
	err := model.ValidateServerID(serverID)
	if err != nil {
		return "", errors_custom.ErrInvalidServerID
	}

	err = model.ValidateAgentName(name)
	if err != nil {
		return "", err
	}

	token, err := GenerateAgentToken()
	if err != nil {
		return "", err
	}

	tokenHash := HashAgentToken(token)

	agent := model.Agent{
		ServerID:  serverID,
		Name:      name,
		TokenHash: tokenHash,
	}

	err = s.repository.Create(ctx, agent)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AgentService) GetByID(ctx context.Context, agentID int) (model.Agent, error) {
	return s.repository.GetByID(ctx, agentID)
}

func (s *AgentService) DeleteByID(ctx context.Context, agentID int) error {
	return s.repository.DeleteByID(ctx, agentID)
}

func (s *AgentService) GetByTokenHash(ctx context.Context, token string) (model.Agent, error) {
	tokenHash := HashAgentToken(token)

	agent, err := s.repository.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return model.Agent{}, errors_custom.ErrInvalidCredentials
	}

	return agent, nil
}

func (s *AgentService) Authenticate(ctx context.Context, token string) (model.Agent, error) {
	tokenHash := HashAgentToken(token)

	agent, err := s.repository.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, errors_custom.ErrAgentNotFound) {
			return model.Agent{}, err
		}

		return model.Agent{}, err
	}

	return agent, nil
}
