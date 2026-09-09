package collector

import (
	"github.com/alexnakagama/server-monitor/internal/agent/request"
)

func CollectMetrics() (request.CreateMetricRequest, error) {
	cpu, err := CPUUsage()
	if err != nil {
		return request.CreateMetricRequest{}, err
	}

	memory, err := MemoryUsage()
	if err != nil {
		return request.CreateMetricRequest{}, err
	}

	disk, err := DiskUsage()
	if err != nil {
		return request.CreateMetricRequest{}, err
	}

	network, err := NetworkUsage()
	if err != nil {
		return request.CreateMetricRequest{}, err
	}

	return request.CreateMetricRequest{
		CPUUsage:       cpu,
		MemoryUsage:    memory,
		DiskUsage:      disk,
		NetworkReceive: network.Received,
		NetworkSent:    network.Sent,
	}, nil
}
