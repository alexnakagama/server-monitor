package repository

import (
	"context"
	"errors"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
	"github.com/jackc/pgx/v5"
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

func (r *AgentRepository) GetByID(ctx context.Context, agentID int) (model.Agent, error) {
	query := `
		SELECT id, server_id, name, token_hash, created_at
		FROM agents
		WHERE id = $1
	`

	var agent model.Agent

	err := r.db.QueryRow(
		ctx,
		query,
		agentID,
	).Scan(
		&agent.ID,
		&agent.ServerID,
		&agent.Name,
		&agent.TokenHash,
		&agent.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Agent{}, errors_custom.ErrAgentNotFound
		}

		return model.Agent{}, err
	}

	return agent, nil
}

func (r *AgentRepository) DeleteByID(ctx context.Context, agentID int) error {
	query := `
		DELETE FROM agents
		WHERE id = $1
	`

	result, err := r.db.Exec(
		ctx,
		query,
		agentID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors_custom.ErrAgentNotFound
	}

	return nil
}

func (r *AgentRepository) GetByTokenHash(ctx context.Context, tokenHash string) (model.Agent, error) {
}
