package repository

import "github.com/jackc/pgx/v5/pgxpool"

type AgentRepository struct {
	db *pgxpool.Pool
}
