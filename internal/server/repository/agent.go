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

func (r *AgentRepository) Create(ctx context.Context, agent model.Agent) error {
	query := `
		INSERT INTO agents (server_id, name, token_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		agent.ServerID,
		agent.Name,
		agent.TokenHash,
	).Scan(
		&agent.ID,
		&agent.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *AgentRepository) GetByID(ctx context.Context, agentID int) (model.Agent, error) {}
