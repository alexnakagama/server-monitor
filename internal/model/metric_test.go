package model

import (
	"errors"
	"testing"

	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
)

func TestMetricValidate(t *testing.T) {
	tests := []struct {
		name    string
		metric  *Metric
		wantErr error
	}{
		{
			name: "valid metric passes validation",
			metric: &Metric{
				ServerID:    1,
				CPUUsage:    50,
				MemoryUsage: 50,
				DiskUsage:   50,
			},
			wantErr: nil,
		},
		{
			name: "usage values at the minimum boundary are valid",
			metric: &Metric{
				ServerID:    1,
				CPUUsage:    0,
				MemoryUsage: 0,
				DiskUsage:   0,
			},
			wantErr: nil,
		},
		{
			name: "usage values at the maximum boundary are valid",
			metric: &Metric{
				ServerID:    1,
				CPUUsage:    100,
				MemoryUsage: 100,
				DiskUsage:   100,
			},
			wantErr: nil,
		},
		{
			name: "zero server id is invalid",
			metric: &Metric{
				ServerID:    0,
				CPUUsage:    50,
				MemoryUsage: 50,
				DiskUsage:   50,
			},
			wantErr: errors_custom.ErrInvalidServerID,
		},
		{
			name: "negative cpu usage is invalid",
			metric: &Metric{
				ServerID:    1,
				CPUUsage:    -1,
				MemoryUsage: 50,
				DiskUsage:   50,
			},
			wantErr: errors_custom.ErrInvalidCPUUsage,
		},
		{
			name: "cpu usage above 100 is invalid",
			metric: &Metric{
				ServerID:    1,
				CPUUsage:    101,
				MemoryUsage: 50,
				DiskUsage:   50,
			},
			wantErr: errors_custom.ErrInvalidCPUUsage,
		},
		{
			name: "negative memory usage is invalid",
			metric: &Metric{
				ServerID:    1,
				CPUUsage:    50,
				MemoryUsage: -1,
				DiskUsage:   50,
			},
			wantErr: errors_custom.ErrInvalidMemoryUsage,
		},
		{
			name: "memory usage above 100 is invalid",
			metric: &Metric{
				ServerID:    1,
				CPUUsage:    50,
				MemoryUsage: 101,
				DiskUsage:   50,
			},
			wantErr: errors_custom.ErrInvalidMemoryUsage,
		},
		{
			name: "negative disk usage is invalid",
			metric: &Metric{
				ServerID:    1,
				CPUUsage:    50,
				MemoryUsage: 50,
				DiskUsage:   -1,
			},
			wantErr: errors_custom.ErrInvalidDiskUsage,
		},
		{
			name: "disk usage above 100 is invalid",
			metric: &Metric{
				ServerID:    1,
				CPUUsage:    50,
				MemoryUsage: 50,
				DiskUsage:   101,
			},
			wantErr: errors_custom.ErrInvalidDiskUsage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.metric.Validate()

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Metric.Validate() error = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
