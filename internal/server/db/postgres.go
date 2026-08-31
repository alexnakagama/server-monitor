package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(ctx context.Context, url string) (*pgxpool.Pool, error) {}
