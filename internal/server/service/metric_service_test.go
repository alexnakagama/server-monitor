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
}
