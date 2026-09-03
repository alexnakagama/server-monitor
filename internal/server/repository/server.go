package repository

import "github.com/jackc/pgx/v5/pgxpool"

type ServerRepository struct {
	db *pgxpool.Pool
}
