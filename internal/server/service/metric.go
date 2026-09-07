package service

import (
	"context"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type MetricRepository interface {
	Create(ctx context.Context, metric model.Metric) error
	GetByID(ctx context.Context, metricID int) (model.Metric, error)
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

	err := metric.ValidateMetric()
}

func (s *MetricService) GetByID(ctx context.Context, metricID int) (model.Metric, error) {}
