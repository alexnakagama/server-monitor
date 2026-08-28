package model

import "time"

type Metric struct {
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsage float64   `json:"memory_usage"`
	DiskUsage   float64   `json:"disk_usage"`
	Timestamp   time.Time `json:"timestamp"`
}
