package model

import (
	"time"

	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
)

type Metric struct {
	ID             int       `json:"id"`
	ServerID       int       `json:"server_id"`
	CPUUsage       float64   `json:"cpu_usage"`
	MemoryUsage    float64   `json:"memory_usage"`
	DiskUsage      float64   `json:"disk_usage"`
	NetworkReceive uint64    `json:"network_receive"`
	NetworkSent    uint64    `json:"network_sent"`
	Timestamp      time.Time `json:"timestamp"`
}

func ValidateMetric(metric Metric) error {
	if metric.ServerID <= 0 {
		return errors_custom.ErrInvalidServerID
	}

	if metric.CPUUsage < 0 || metric.CPUUsage > 100 {
		return errors_custom.ErrInvalidCPUUsage
	}

	if metric.MemoryUsage < 0 || metric.MemoryUsage > 100 {
		return errors_custom.ErrInvalidMemoryUsage
	}

	if metric.DiskUsage < 0 || metric.DiskUsage > 100 {
		return errors_custom.ErrInvalidDiskUsage
	}

	return nil
}
