package model

import "time"

type Metric struct {
	ID             int       `json:"id"`
	AgentID        int       `json:"agent_id"`
	CPUUsage       float64   `json:"cpu_usage"`
	MemoryUsage    float64   `json:"memory_usage"`
	DiskUsage      float64   `json:"disk_usage"`
	NetworkReceive uint64    `json:"network_receive"`
	NetworkSent    uint64    `json:"network_sent"`
	Timestamp      time.Time `json:"timestamp"`
}
