package service

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type MetricRepository interface {
	Create(ctx context.Context, metric model.Metric) error
	GetByID(ctx context.Context, metricID int) (model.Metric, error)
	GetByServerID(ctx context.Context, serverID int, limit int) ([]model.Metric, error)
}

type MetricService struct {
	repository MetricRepository
}

func NewMetricService(repository MetricRepository) *MetricService {
	return &MetricService{
		repository: repository,
	}
}

type MetricInput struct {
	ServerID       int
	CPUUsage       float64
	MemoryUsage    float64
	DiskUsage      float64
	NetworkReceive uint64
	NetworkSent    uint64
}

func (s *MetricService) Create(ctx context.Context, input MetricInput) error {
	metric := model.Metric{
		ServerID:       input.ServerID,
		CPUUsage:       input.CPUUsage,
		MemoryUsage:    input.MemoryUsage,
		DiskUsage:      input.DiskUsage,
		NetworkReceive: input.NetworkReceive,
		NetworkSent:    input.NetworkSent,
	}

	err := metric.Validate()
	if err != nil {
		return err
	}

	return s.repository.Create(ctx, metric)
}

func (s *MetricService) GetByID(ctx context.Context, metricID int) (model.Metric, error) {
	return s.repository.GetByID(ctx, metricID)
}

func (s *MetricService) GetByServerID(ctx context.Context, serverID int, limit int) ([]model.Metric, error) {
	return s.repository.GetByServerID(ctx, serverID, limit)
}
