package service

import (
	"context"
	"testing"
	"time"

	"github.com/alexnakagama/server-monitor/internal/model"
)

var _ MetricRepository = (*MetricRepositoryMock)(nil)

type MetricRepositoryMock struct {
	create          func(ctx context.Context, metric model.Metric) error
	getByID         func(ctx context.Context, metricID int) (model.Metric, error)
	getByServerID   func(ctx context.Context, serverID int, filters model.MetricFilters) ([]model.Metric, error)
	deleteOlderThan func(ctx context.Context, before time.Time) error
}

func (m *MetricRepositoryMock) Create(ctx context.Context, metric model.Metric) error {
	return m.create(ctx, metric)
}

func (m *MetricRepositoryMock) GetByID(ctx context.Context, metricID int) (model.Metric, error) {
	return m.getByID(ctx, metricID)
}

func (m *MetricRepositoryMock) GetByServerID(ctx context.Context, serverID int, filters model.MetricFilters) ([]model.Metric, error) {
	return m.getByServerID(ctx, serverID, filters)
}

func (m *MetricRepositoryMock) DeleteOlderThan(ctx context.Context, before time.Time) error {
	return m.deleteOlderThan(ctx, before)
}

func TestMetricService_DeleteOldMetrics(t *testing.T) {
	var receiveBefore time.Time

	repository := &MetricRepositoryMock{
		deleteOlderThan: func(ctx context.Context, before time.Time) error {
			receiveBefore = before
			return nil
		},
	}

	service := NewMetricService(repository, 30)

	start := time.Now()

	err := service.DeleteOldMetrics(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	end := time.Now()

	min := start.AddDate(0, 0, -30)
	max := end.AddDate(0, 0, -30)

	if receiveBefore.Before(min) || receiveBefore.After(max) {
		t.Errorf("unexpectec cutoff: %v", receiveBefore)
	}
}

func TestMetricService_Create(t *testing.T) {
	var createdMetric model.Metric

	repository := &MetricRepositoryMock{
		create: func(ctx context.Context, metric model.Metric) error {
			createdMetric = metric
			return nil
		},
	}

	service := NewMetricService(repository, 30)

	err := service.Create(context.Background(), MetricInput{
		ServerID:       1,
		CPUUsage:       50,
		MemoryUsage:    60,
		DiskUsage:      70,
		NetworkReceive: 1000,
		NetworkSent:    2000,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if createdMetric.ServerID != 1 {
		t.Errorf("expected server id 1, got: %d", createdMetric.ServerID)
	}

	if createdMetric.CPUUsage != 50 {
		t.Errorf("expected cpu usage 50, got: %f", createdMetric.CPUUsage)
	}

	if createdMetric.MemoryUsage != 60 {
		t.Errorf("expected memory usage 60, got: %f", createdMetric.MemoryUsage)
	}

	if createdMetric.DiskUsage != 70 {
		t.Errorf("expected disk usage 70, got: %f", createdMetric.DiskUsage)
	}

	if createdMetric.NetworkReceive != 1000 {
		t.Errorf("expected network received 1000, got: %d", createdMetric.NetworkReceive)
	}

	if createdMetric.NetworkSent != 2000 {
		t.Errorf("expected network received 2000, got: %d", createdMetric.NetworkSent)
	}
}

func TestMetricService_GetByID(t *testing.T) {
	metric := model.Metric{
		ID:             1,
		ServerID:       1,
		CPUUsage:       10,
		MemoryUsage:    20,
		DiskUsage:      30,
		NetworkReceive: 1000,
		NetworkSent:    2000,
	}

	repository := &MetricRepositoryMock{
		getByID: func(ctx context.Context, metricID int) (model.Metric, error) {
			return metric, nil
		},
	}

	service := NewMetricService(repository, 30)

	metricFound, err := service.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if metricFound.ServerID != metric.ServerID {
		t.Errorf("expected server id: %d, got: %d", metric.ServerID, metricFound.ServerID)
	}

	if metricFound.CPUUsage != metric.CPUUsage {
		t.Errorf("expected cpu usage: %f, got: %f", metric.CPUUsage, metricFound.CPUUsage)
	}

	if metricFound.MemoryUsage != metric.MemoryUsage {
		t.Errorf("expected memory usage usage: %f, got: %f", metric.MemoryUsage, metricFound.MemoryUsage)
	}

	if metricFound.DiskUsage != metric.DiskUsage {
		t.Errorf("expected disk usage: %f, got: %f", metric.DiskUsage, metricFound.DiskUsage)
	}

	if metricFound.NetworkReceive != metric.NetworkReceive {
		t.Errorf("expected network receive: %d, got: %d", metric.NetworkReceive, metricFound.NetworkReceive)
	}

	if metricFound.NetworkSent != metric.NetworkSent {
		t.Errorf("expected network sent: %d, got: %d", metric.NetworkSent, metricFound.NetworkSent)
	}
}
