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

func (r *ServerRepository) GetByName(ctx context.Context, name string) (model.Server, error) {
	var server model.Server

	query := `
		SELECT id, name, hostname, os, created_at
		FROM servers
		WHERE name = $1
	`

	err := r.db.QueryRow(
		ctx,
		query,
		name,
	).Scan(
		&server.ID,
		&server.Name,
		&server.Hostname,
		&server.OS,
		&server.CreatedAt,
	)

	if err != nil {
		return model.Server{}, err
	}

	return server, nil
}

func (r *ServerRepository) GetAll(ctx context.Context) ([]model.Server, error) {
	query := `
		SELECT id, name, hostname, os, created_at
		FROM servers
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []model.Server

	for rows.Next() {
		var server model.Server

		err := rows.Scan(
			&server.ID,
			&server.Name,
			&server.Hostname,
			&server.OS,
			&server.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		servers = append(servers, server)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (r *ServerRepository) GetByOS(ctx context.Context, os string) ([]model.Server, error) {
	query := `
		SELECT id, name, hostname, os, created_at
		FROM servers
		WHERE os = $1
	`

	rows, err := r.db.Query(ctx, query, os)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []model.Server

	for rows.Next() {
		var server model.Server

		err := rows.Scan(
			&server.ID,
			&server.Name,
			&server.Hostname,
			&server.OS,
			&server.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		servers = append(servers, server)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return servers, nil
}
