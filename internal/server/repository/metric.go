package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type MetricRepository struct {
	db *pgxpool.Pool
}

func NewMetricRepository(db *pgxpool.Pool) *MetricRepository {
	return &MetricRepository{
		db: db,
	}
}
