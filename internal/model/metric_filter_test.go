package model

import (
	"errors"
	"testing"
	"time"

	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
)

func TestMetricFiltersValidate(t *testing.T) {
	now := time.Date(2026, 9, 12, 12, 0, 0, 0, time.UTC)
	before := now.Add(-1 * time.Hour)
	after := now.Add(1 * time.Hour)

	tests := []struct {
		name    string
		filters MetricFilters
		wantErr error
	}{
		{
			name:    "limit below 1 is invalid",
			filters: MetricFilters{Limit: 0},
			wantErr: errors_custom.ErrInvalidLimit,
		},
		{
			name:    "negative limit is invalid",
			filters: MetricFilters{Limit: -10},
			wantErr: errors_custom.ErrInvalidLimit,
		},
		{
			name:    "limit above 1000 is invalid",
			filters: MetricFilters{Limit: 1001},
			wantErr: errors_custom.ErrInvalidLimit,
		},
		{
			name:    "limit at the minimum boundary is valid",
			filters: MetricFilters{Limit: 1},
			wantErr: nil,
		},
		{
			name:    "limit at the maximum boundary is valid",
			filters: MetricFilters{Limit: 1000},
			wantErr: nil,
		},
		{
			name:    "time range where From is before To is valid",
			filters: MetricFilters{From: &before, To: &after, Limit: 50},
			wantErr: nil,
		},
		{
			name:    "time range where From equals To is valid",
			filters: MetricFilters{From: &now, To: &now, Limit: 50},
			wantErr: nil,
		},
		{
			name:    "time range where From is after To is invalid",
			filters: MetricFilters{From: &after, To: &before, Limit: 50},
			wantErr: errors_custom.ErrInvalidTimeRange,
		},
		{
			name:    "time range skipped when From is nil",
			filters: MetricFilters{From: nil, To: &after, Limit: 50},
			wantErr: nil,
		},
		{
			name:    "time range skipped when To is nil",
			filters: MetricFilters{From: &before, To: nil, Limit: 50},
			wantErr: nil,
		},
		{
			name:    "time range skipped when both From and To are nil",
			filters: MetricFilters{Limit: 50},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.filters.Validate()

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("MetricFilters.Validate() error = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
