package repository

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AgentRepository struct {
	db *pgxpool.Pool
}

func NewAgentRepository(db *pgxpool.Pool) *AgentRepository {
	return &AgentRepository{
		db: db,
	}
}

func (r *AgentRepository) Create(ctx context.Context, agent model.Agent) error {}
