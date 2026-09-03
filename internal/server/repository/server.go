package repository

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServerRepository struct {
	db *pgxpool.Pool
}

func NewServerRepository(db *pgxpool.Pool) *ServerRepository {
	return &ServerRepository{
		db: db,
	}
}

func (r *ServerRepository) Create(ctx context.Context, server model.Server) error {
	query := `
		INSERT INTO servers (name, hostname, os)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		server.Name,
		server.Hostname,
		server.OS,
	)

	return err
}
